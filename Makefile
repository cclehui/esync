build:goget
	go build -o run main.go
autocompletor:goget
	  go install ./...
goget:
	  go mod tidy
gotest:
	  go test -v ./...
