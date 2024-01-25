run:
	go run main.go

unit-tests:
	go test ./...

functional-tests:
	go test ./functional_tests/transformer_test.go

build:
	docker build . -t shadowshotx/product-go-micro
