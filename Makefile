wallet-services:
	env GO111MODULE=on go build -v $(LDFLAGS) ./cmd/wallet-services

clean:
	rm wallet-services

test:
	go test -v ./...

lint:
	golangci-lint run ./...

.PHONY: \
	wallet-services \
	clean \
	test \
	lint
