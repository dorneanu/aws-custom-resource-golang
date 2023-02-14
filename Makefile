build: lambda

lambda:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./build/lambda-function.bin ./cmd/main.go
