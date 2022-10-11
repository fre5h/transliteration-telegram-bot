aws_lambda:
    GOARCH=amd64 GOOS=linux go build -o main cmd/main.go && zip archive.zip main