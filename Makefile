GO = /usr/bin/go

all: clean build

clean:
	rm -rf bin/

build:
	${GO} build -o bin/ ./...

install:
	cp bin/pnp-manager /usr/local/bin/
	mkdir -p /etc/pnp-manager/
	cp example.config.toml /etc/pnp-manager/config.toml

uninstall:
	rm -rf /etc/pnp-manager/
	rm /usr/local/bin/pnp-manager
