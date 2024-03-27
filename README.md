# MONIKA GATEWAY (API GATEWAY FOR MONIKA)

## Compile

This application uses [mattn/go-sqlite3](github.com/mattn/go-sqlite3) and thus needs gcc and has to be compiled

### Cross Compiling from macOS

The simplest way to cross compile from macOS is to use xgo.

Steps:

- Install musl-cross (brew install FiloSottile/musl-cross/musl-cross).
- Run CC=x86_64-linux-musl-gcc CXX=x86_64-linux-musl-g++ GOARCH=amd64 GOOS=linux CGO_ENABLED=1 go build -ldflags "-linkmode external -extldflags -static".
