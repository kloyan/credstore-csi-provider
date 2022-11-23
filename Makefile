REGISTRY_NAME?=kloyan
IMAGE_NAME=credstore-csi-provider
VERSION?=0.0.0-dev
IMAGE_TAG=$(REGISTRY_NAME)/$(IMAGE_NAME):$(VERSION)
KIND_CLUSTER_NAME?=credstore-cluster

default: test

clean:
	rm -rf dist/ *.o

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

kind:
	kind create cluster --name $(KIND_CLUSTER_NAME) --wait 30s

	podman save $(IMAGE_TAG) -o image.tar
	kind load image-archive --name $(KIND_CLUSTER_NAME) image.tar
	rm -f image.tar

	kubectl kustomize --enable-helm deploy/ | kubectl apply -f-

delete-kind:
	kind delete cluster --name $(KIND_CLUSTER_NAME)
