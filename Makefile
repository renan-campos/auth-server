all: bin/authentication-server bin/example-client

bin/authentication-server: cmd/authentication-server/*.go
	go build -o bin/authentication-server cmd/authentication-server/*.go

bin/example-client: cmd/client/*.go
	go build -o bin/example-client cmd/client/*.go

.PHONY: clean
clean:
	rm bin/*
