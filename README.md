# Hashcash tls-client

A Go implementation of a TLS client for secure communication with a server over a TCP connection.
It generates a string with a hash that has a specified number of leading zeros (the difficulty) by concatenating a nonce with a provided prefix data. And sends required data to server like contact information.

## Features
 - Loads a certificate and private key pair to configure a secure TLS connection.
 - Connects to a server specified by a host:port address.
 - Authenticates the server using the provided certificate.
 - Reads and writes messages to the server.
 - Handles errors and logging.
 - Requirements
 - Go (1.14+)

## Usage

Use the following flags to configure the tls-client:

- `addr`: host:port address of the server to connect to (default: "localhost:8080").
- `cert`: location of the certificate to use for the TLS connection (default: ".\cert\cert1.pem").
- `key`: location of the private key to use for the TLS connection (default: ".\cert\key.pem").
- `l`: string length (default: 8).
- `g`: max number of goroutines (default: 8).
- `p`: write profiling (default: false).


```sh
go run cmd/client/client.go -addr localhost:8080 -cert cert1.pem -key key.pem
```


# Hashcash

Hashcash is a Go program that implements the hashcash algorithm. It generates a string with a hash that has a specified number of leading zeros (the difficulty) by concatenating a nonce with a provided prefix data.

## Usage

Run the program using the following command:

```bash
go run cmd/main.go -d [difficulty] -l [string length] -g [max goroutines] -a [auth] -p [profiling status]
```

### where:

- `difficulty` is the number of leading zeros required in the hash (1-9). Default is 7.
- `string length` is the length of the string generated. Default is 8.
- `max goroutines` is the maximum number of goroutines to use. Default is 8.
- `auth` is the prefix data to use in the hash. Default is "auth".
- `profiling` status is a boolean indicating whether to write profiling information to disk (true) or not (false). Default is false.


## Profiling
The program has the option to write profiling information to disk by setting the `-p` flag to `true`. 
This will generate two files `cpu.pro`f and `mem.prof` in the current directory, which can be analyzed using tools such as `go tool pprof`.

## Example
```sh
go run cmd/main.go -d 7 -l 8 -g 8 -a "prefix" -p true
```

This command will run the program with a difficulty of 7, a string length of 8, a maximum of 8 goroutines, and a prefix of "prefix". It will also write profiling information to disk.
