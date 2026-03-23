go env -w GOOS=linux
go env -w GOARCH=amd64
go mod tidy
go build .