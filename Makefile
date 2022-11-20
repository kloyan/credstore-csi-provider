REGISTRY_NAME?=docker.io/kloyan
IMAGE_NAME=credstore-csi-provider
VERSION?=0.0.0-dev
IMAGE_TAG=$(REGISTRY_NAME)/$(IMAGE_NAME):$(VERSION)

default: test

clean: rm -rf *.o

test:
	go test ./...

mod:
	@go mod tidy

build:
	CGO_ENABLED=0 go build -o dist/ .

image:
	podman build \
		--no-cache \
		--tag $(IMAGE_TAG) \
		.
