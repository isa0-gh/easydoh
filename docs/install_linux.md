
1. **Build and install using Makefile**

```bash
git clone https://github.com/isa0-gh/easydoh.git
cd easydoh
make
sudo make install
````

> The Makefile will automatically detect your init system and install the appropriate service script.

2. **Configuration file**

Create or edit `/etc/easydoh/config.json`:

```json
{
  "resolver": "https://dns.adguard-dns.com/dns-query",
  "ttl": 300,
  "bind_address": "127.0.0.1:53"
}
```

* `resolver` — choose from [this list](docs/servers.md)
* `bind_address` — IP and port for the server to listen on

---

## Service Management

### Systemd

```bash
# Reload systemd and enable the service
sudo systemctl daemon-reload
sudo systemctl enable easydoh.service

# Start/Stop the service
sudo systemctl start easydoh.service
sudo systemctl stop easydoh.service

# Check status
sudo systemctl status easydoh.service
```

## Usage

Once installed and started, EasyDoH will listen on the configured `bind_address` and forward queries to the chosen resolver. You can use it as your local DoH server by pointing your applications or system DNS settings to it.
