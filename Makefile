all: bin/authentication-server bin/client bin/single-use-token

bin/authentication-server: cmd/authentication-server/*.go
	go build -o bin/authentication-server cmd/authentication-server/*.go

bin/client: cmd/client/*.go 
	go build -o bin/client cmd/client/*.go

bin/single-use-token: cmd/single-use-token/*.go
	go build -o bin/single-use-token cmd/single-use-token/*.go

.PHONY: clean
clean:
	rm bin/*
