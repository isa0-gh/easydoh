
1. **Build and install using Makefile**

```bash
git clone https://github.com/isa0-gh/resolv.git
cd resolv
make
sudo make install
````

> The Makefile will automatically detect your init system and install the appropriate service script.

2. **Configuration file**

Create or edit `/etc/resolv/config.toml`:

```toml
resolver = "https://dns.adguard-dns.com/dns-query"
ttl = 300
bind_address = "0.0.0.0:53"

[hosts]
"*.home" = "127.0.0.1"
"test.local" = "192.168.1.100"
```

* `resolver` — choose from [this list](docs/servers.md)
* `bind_address` — IP and port for the server to listen on

---

## Service Management (Docker)

To run resolv continuously, we recommend using Docker. You should mount a volume to persist and provide your `config.toml`.

### Using Docker Compose (Recommended)
```bash
# Start the service in the background
docker compose up -d

# View logs
docker compose logs -f
```

### Using Docker Run
```bash
# Pull the image
docker pull ghcr.io/isa0-gh/resolv:latest

# Run the container (mapping host UDP 53 to container UDP 53)
docker run -d --name resolv \
  -p 53:53/udp \
  -v /etc/resolv:/etc/resolv \
  --restart unless-stopped \
  ghcr.io/isa0-gh/resolv:latest
```

## Usage

Once installed and started, EasyDoH will listen on the configured `bind_address` and forward queries to the chosen resolver. You can use it as your local DoH server by pointing your applications or system DNS settings to it.
