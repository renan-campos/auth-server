all: bin/authentication-server bin/client bin/single-use-token bin/token-generator bin/pem-to-jwk

bin/authentication-server: cmd/authentication-server/*.go
	go build -o bin/authentication-server cmd/authentication-server/*.go

bin/client: cmd/client/*.go 
	go build -o bin/client cmd/client/*.go

bin/single-use-token: cmd/single-use-token/*.go
	go build -o bin/single-use-token cmd/single-use-token/*.go

bin/token-generator: cmd/token-generator/*.go
	go build -o bin/token-generator cmd/token-generator/*.go

bin/pem-to-jwk: cmd/pem-to-jwk/*.go
	go build -o bin/pem-to-jwk cmd/pem-to-jwk/*.go

.PHONY: clean
clean:
	rm bin/*
