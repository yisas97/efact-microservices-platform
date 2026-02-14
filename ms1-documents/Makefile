.PHONY: swagger build run test clean

swagger:
	@swag init -g cmd/api/main.go --output ./docs

build: swagger
	@go build -o bin/ms1-documents.exe cmd/api/main.go

run: swagger
	@go run cmd/api/main.go

test:
	@go test -v -cover ./...

clean:
	@rm -rf bin/ docs/

swagger-install:
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go get -u github.com/swaggo/gin-swagger
	@go get -u github.com/swaggo/files
