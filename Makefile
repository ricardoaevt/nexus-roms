.PHONY: test coverage race coverage-html

test:
	go test ./... -v

race:
	go test ./... -race -v

coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out

coverage-html:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
