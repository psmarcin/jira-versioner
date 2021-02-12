lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix

update-deps:
	go get -u ./...

test:
	go test ./...
