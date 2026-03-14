package service

import (
	"github.com/isa0-gh/resolv/internal/cache"
	"github.com/isa0-gh/resolv/internal/config"
	"github.com/isa0-gh/resolv/internal/local"
	"github.com/isa0-gh/resolv/internal/resolver"
)

type ServiceRepo struct {
	Config   *config.Config
	Cache    *cache.CacheDB
	Local    *local.Matcher
	Resolver *resolver.Resolver
}

func NewServiceRepo(conf *config.Config, cdb *cache.CacheDB, matcher *local.Matcher, r *resolver.Resolver) *ServiceRepo {
	return &ServiceRepo{
		Config:   conf,
		Cache:    cdb,
		Local:    matcher,
		Resolver: r,
	}
}
