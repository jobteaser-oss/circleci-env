BIN_DIR = $(CURDIR)/bin

all: build

build: FORCE
	GOBIN=$(BIN_DIR) go install ./...

vet:
	go list ./... | xargs go vet

clean:
	$(RM) -r ./bin

FORCE:

.PHONY: all build vet clean
