GIT_COMMIT=$(shell git rev-list -1 HEAD)
BUILD_DATE=$(shell date +"%Y-%m-%dT%H:%M:%SZ")

DOCKER?=podman
BIN_DIR?=dist/
REGISTRY_NAME?=ghcr.io/kloyan
IMAGE_NAME=credstore-csi-provider
VERSION?=dev
IMAGE_TAG=$(REGISTRY_NAME)/$(IMAGE_NAME):$(VERSION)
KIND_CLUSTER_NAME?=credstore-cluster
PKG=github.com/kloyan/credstore-csi-provider/internal/version
LDFLAGS?="-X '$(PKG).BuildVersion=$(VERSION)' \
		  -X '$(PKG).BuildDate=$(BUILD_DATE)' \
		  -X '$(PKG).GitCommit=$(GIT_COMMIT)'"

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
	CGO_ENABLED=0 go build -ldflags $(LDFLAGS) -o $(BIN_DIR) .

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
