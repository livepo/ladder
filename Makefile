default: build

.PHONEY: build
build:
	go build -o ladder-server cmd/ladder-server.go
	go build -o ladder-local cmd/ladder-local.go

.PHONEY: clean
clean:
	rm -rf .gosocks
	rm -rf ladder-local
	rm -rf ladder-server
