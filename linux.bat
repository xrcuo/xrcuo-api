go env -w GOOS=linux
go env -w GOARCH=amd64
go get -u
go mod tidy
go build .