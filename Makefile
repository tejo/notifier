test:
	@(go test ./... -v)

build:
	@(go build -o notify cmd/notify/main.go)

run:
	@(go run cmd/notify/main.go)

server:
	@(go build -o notify cmd/notify/main.go && ./notify -dummyserver) 