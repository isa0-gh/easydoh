package cache

import (
	"fmt"
	"sync"
	"time"

	"github.com/miekg/dns"
)

type CacheDB struct {
	Mu sync.RWMutex
	DB map[string]dns.Msg
}

func New() *CacheDB {
	return &CacheDB{
		DB: make(map[string]dns.Msg),
	}
}

func (c *CacheDB) Add(message []byte, response []byte) error {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	req := new(dns.Msg)
	resp := new(dns.Msg)
	if err := req.Unpack(message); err != nil {
		return err
	}
	if len(req.Question) == 0 {
		return fmt.Errorf("request contains no question section")
	}

	if err := resp.Unpack(response); err != nil {
		return err
	}

	question := req.Question[0]
	key := fmt.Sprintf("%s|%d", question.Name, question.Qtype)
	c.DB[key] = *resp
	return nil
}

func (c *CacheDB) Get(message []byte) ([]byte, bool, error) {
	msg := new(dns.Msg)
	if err := msg.Unpack(message); err != nil {
		return nil, false, err
	}
	if len(msg.Question) == 0 {
		return nil, false, fmt.Errorf("request contains no question section")
	}
	question := msg.Question[0]
	key := fmt.Sprintf("%s|%d", question.Name, question.Qtype)
	c.Mu.RLock()
	defer c.Mu.RUnlock()
	for k, v := range c.DB {
		if k == key {
			v.Id = msg.Id
			resp, err := v.Pack()
			return resp, true, err
		}
	}
	return nil, false, nil
}

func (c *CacheDB) Expire(data []byte, ttl int) {
	msg := new(dns.Msg)
	if err := msg.Unpack(data); err != nil {
		return
	}
	if len(msg.Question) == 0 {
		return
	}
	question := msg.Question[0]
	key := fmt.Sprintf("%s|%d", question.Name, question.Qtype)
	c.Mu.Lock()
	go func(k string, d int) {
		time.Sleep(time.Duration(d) * time.Second)
		delete(c.DB, key)
	}(key, ttl)

	c.Mu.Unlock()
}
