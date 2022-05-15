GO = /usr/bin/go

build:
	${GO} build -o bin/ ./...

install:
	cp bin/* /usr/local/bin/
	mkdir -p /etc/pnp-manager/
	cp example.config.toml /etc/pnp-manager/
