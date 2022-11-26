DOCKER?=podman
BIN_DIR?=dist/
REGISTRY_NAME?=kloyan
IMAGE_NAME=credstore-csi-provider
VERSION?=0.0.0-dev
IMAGE_TAG=$(REGISTRY_NAME)/$(IMAGE_NAME):$(VERSION)
KIND_CLUSTER_NAME?=credstore-cluster

default: build

.PHONY: clean
clean:
	rm -rfv dist/ *.o

.PHONY: test
test:
	go test ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: lint
lint: fmt
	golint ./...

.PHONY: vet
vet: fmt
	go vet ./...

.PHONY: build
build: vet test
	CGO_ENABLED=0 go build -o $(BIN_DIR) .

.PHONY: image
image:
	$(DOCKER) build \
		--no-cache \
		--tag $(IMAGE_TAG) \
		.

.PHONY: setup-kind
setup-kind:
	kind create cluster --name $(KIND_CLUSTER_NAME) --wait 30s

	$(DOCKER) save $(IMAGE_TAG) -o image.tar
	kind load image-archive --name $(KIND_CLUSTER_NAME) image.tar
	rm -f image.tar

	kubectl kustomize --enable-helm deploy/ | kubectl apply -f-

.PHONY: delete-kind
delete-kind:
	kind delete cluster --name $(KIND_CLUSTER_NAME)
