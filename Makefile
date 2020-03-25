GOVARS  := CGO_ENABLED=1
VERSION := $(shell git describe --tags --always --dirty="-dev")

# Currently arm and mac builds broken based on the cross compile dependencies.

release: clean github-release dist
	github-release release \
		--user ahmet2mir \
		--repo terraform-provider-freeipa \
		--tag $(VERSION) \
		--name $(VERSION)
		--security-token $$GITHUB_TOKEN

	# GNU/Linux - X86
	github-release upload \
		--user ahmet2mir \
		--repo terraform-provider-freeipa \
		--tag $(VERSION) \
		--name terraform-provider-freeipa_$(VERSION)-linux-amd64 \
		--file terraform-provider-freeipa_$(VERSION)-linux-amd64 \
		--security-token $$GITHUB_TOKEN

dist: goget
	# GNU/Linux - X86
	$(GOVARS) GOOS=linux GOARCH=amd64 go build -o terraform-provider-freeipa_$(VERSION)-linux-amd64

	# arm
	# $(GOVARS) GOOS=linux CC=arm-linux-gnueabi-gcc GOARCH=arm go build -o terraform-provider-k8s_$(VERSION)-linux-arm
	# $(GOVARS) GOOS=linux GOARCH=arm64 go build -o terraform-provider-k8s_$(VERSION)-linux-arm64

	# macOS
	# $(GOVARS) GOOS=darwin GOARCH=amd64 go build -o terraform-provider-k8s_$(VERSION)-darwin-amd64

goget:
	go get

clean:
	rm -rf terraform-provider-freeipa*

github-release:
	go get -u github.com/aktau/github-release

.PHONY: clean github-release
