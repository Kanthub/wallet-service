RRM_ABI_ARTIFACT := ./abis/ReferralRewardManager.sol/ReferralRewardManager.json

event-services:
	env GO111MODULE=on go build -v $(LDFLAGS) ./cmd/event-services

clean:
	rm event-services

test:
	go test -v ./...

lint:
	golangci-lint run ./...

bindings: binding-rrm

binding-rrm:
	$(eval temp := $(shell mktemp))

	cat $(RRM_ABI_ARTIFACT) \
		| jq -r .bytecode.object > $(temp)

	cat $(RRM_ABI_ARTIFACT) \
		| jq .abi \
		| abigen --pkg bindings \
		--abi - \
		--out bindings/rrm_manager.go \
		--type ReferralRewardManager \
		--bin $(temp)

		rm $(temp)

.PHONY: \
	event-services \
	bindings \
	binding-rrm \
	clean \
	test \
	lint
