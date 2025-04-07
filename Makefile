#!/usr/bin/make -f

PACKAGES_NOSIMULATION=$(shell go list ./... | grep -v '/simulation')
PACKAGES_NOSIMULATION_TEST=$(shell go list ./... | grep -v '/simulation' | grep -v '/tests' | grep -v '/feemarket' )
PACKAGES_SIMTEST=$(shell go list ./... | grep '/simulation')
VERSION ?= $(shell echo $(shell git describe --tags `git rev-list --tags="v*" --max-count=1`) | sed 's/^v//')
TMVERSION := $(shell go list -m github.com/tendermint/tendermint | sed 's:.* ::')
COMMIT := $(shell git log -1 --format='%H')
INSTALL_DIR ?= $(shell go env GOPATH)/bin
LEDGER_ENABLED ?= true
BINDIR ?= $(GOPATH)/bin
TREASURENET_BINARY = treasurenetd
TREASURENET_DIR = treasurenet
BUILDDIR ?= $(CURDIR)/build
SUPPORTED_ARCHS = amd64 arm arm64
BUILD_ARCH = $(shell uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')
SIMAPP = ./app
HTTPS_GIT := https://github.com/treasurenetprotocol/treasurenet.git
PROJECT_NAME = $(shell git remote get-url origin | xargs basename -s .git)
DOCKER := $(shell which docker)
DOCKER_BUF := $(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace bufbuild/buf:1.0.0-rc8
NAMESPACE := treasurenetd
DOCKER_IMAGE := $(NAMESPACE)
DOCKER_TAG := latest
export GO111MODULE = on

# Default target executed when no arguments are given to make.
default_target: all

.PHONY: default_target

# process build tags

build_tags = netgo
ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags += ledger
    endif
  else
    UNAME_S = $(shell uname -s)
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        build_tags += ledger
      endif
    endif
  endif
endif

ifeq (cleveldb,$(findstring cleveldb,$(COSMOS_BUILD_OPTIONS)))
  build_tags += gcc
endif
build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

# process linker flags

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=treasurenet \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=$(TREASURENET_BINARY) \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
			-X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)" \
			-X github.com/tendermint/tendermint/version.TMCoreSemVer=$(TMVERSION)

# DB backend selection
ifeq (cleveldb,$(findstring cleveldb,$(COSMOS_BUILD_OPTIONS)))
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=cleveldb
endif
ifeq (badgerdb,$(findstring badgerdb,$(COSMOS_BUILD_OPTIONS)))
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=badgerdb
endif
# handle rocksdb
ifeq (rocksdb,$(findstring rocksdb,$(COSMOS_BUILD_OPTIONS)))
  CGO_ENABLED=1
  BUILD_TAGS += rocksdb
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=rocksdb
endif
# handle boltdb
ifeq (boltdb,$(findstring boltdb,$(COSMOS_BUILD_OPTIONS)))
  BUILD_TAGS += boltdb
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=boltdb
endif

ifeq (,$(findstring nostrip,$(COSMOS_BUILD_OPTIONS)))
  ldflags += -w -s
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'
# check for nostrip option
ifeq (,$(findstring nostrip,$(COSMOS_BUILD_OPTIONS)))
  BUILD_FLAGS += -trimpath
endif

# # The below include contains the tools and runsim targets.
# include contrib/devtools/Makefile

###############################################################################
###                                  Build                                  ###
###############################################################################

# BUILD_TARGETS := build install

# build: BUILD_ARGS=-o $(BUILDDIR)/
# build-linux:
# 	GOOS=linux GOARCH=amd64 LEDGER_ENABLED=false $(MAKE) build

BUILD_TARGETS := build install

define ARCH_BUILD_TEMPLATE
build-$(1):
	@mkdir -p $(BUILDDIR)/$(1)  # 先创建目录
	@echo "Building for $(1)..."
	GOOS=linux GOARCH=$(1) GOARM=7 LEDGER_ENABLED=$(LEDGER_ENABLED) \
	go build $(BUILD_FLAGS) -o $(BUILDDIR)/$(1)/$(TREASURENET_BINARY) ./cmd/treasurenetd
endef

$(foreach arch,$(SUPPORTED_ARCHS),$(eval $(call ARCH_BUILD_TEMPLATE,$(arch))))

build: build-$(BUILD_ARCH)

build-multiarch: $(addprefix build-,$(SUPPORTED_ARCHS))
	@echo "Multi-architecture build completed in $(BUILDDIR)/"
	

install: | $(BUILDDIR)
	@echo "Installing to GoPATH: $(INSTALL_DIR)"
	@mkdir -p $(INSTALL_DIR)
	@for arch in $(SUPPORTED_ARCHS); do \
		bin_path="$(BUILDDIR)/$${arch}/$(TREASURENET_BINARY)"; \
		if [ -f "$${bin_path}" ]; then \
			install -m 755 "$${bin_path}" "$(INSTALL_DIR)/$(TREASURENET_BINARY)-$${arch}"; \
			echo "✔ Installed $${arch} => $(INSTALL_DIR)/$(TREASURENET_BINARY)-$${arch}"; \
		else \
			echo "⚠ $${arch} binary missing at $${bin_path}"; \
		fi \
	done

$(BUILDDIR)/:
	mkdir -p $(BUILDDIR)/

docker-build:
	# TODO replace with kaniko
	#docker build -t ${DOCKER_IMAGE}:${DOCKER_TAG} .
	#docker tag ${DOCKER_IMAGE}:${DOCKER_TAG} ${DOCKER_IMAGE}:latest
		@$(foreach arch,$(SUPPORTED_ARCHS), \
		docker buildx build --platform linux/$(arch) -t ${DOCKER_IMAGE}:${DOCKER_TAG}-$(arch) . ; \
	)
	docker tag ${DOCKER_IMAGE}:${DOCKER_TAG}-$(BUILD_ARCH) ${DOCKER_IMAGE}:latest
	# docker tag ${DOCKER_IMAGE}:${DOCKER_TAG} ${DOCKER_IMAGE}:${COMMIT_HASH}
	# update old container
	docker rm treasurenetd || true
	# create a new container from the latest image
	docker create --name treasurenetd -t -i treasurenetd:latest treasurenetd
	# move the binaries to the ./build directory
	mkdir -p ./build/
	docker cp treasurenetd:/usr/bin/treasurenetd ./build/

$(MOCKS_DIR):
	mkdir -p $(MOCKS_DIR)

distclean: clean tools-clean

# clean:
# 	rm -rf \
#     $(BUILDDIR)/ \
#     artifacts/ \
#     tmp-swagger-gen/
clean:
	rm -rf \
		$(BUILDDIR)/ \
		artifacts/ \
		tmp-swagger-gen/ \
		$(foreach arch,$(SUPPORTED_ARCHS),$(INSTALL_DIR)/$(TREASURENET_BINARY)-$(arch))

all: build

build-all: tools build lint test

.PHONY: distclean clean build-all

###############################################################################
###                                Releasing                                ###
###############################################################################

PACKAGE_NAME:=github.com/treasurenetprotocol/treasurenet
GOLANG_CROSS_VERSION  = v1.18
GOPATH ?= '$(HOME)/go'
release-dry-run:
	docker run \
		--rm \
		--privileged \
		-e CGO_ENABLED=1 \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-v ${GOPATH}/pkg:/go/pkg \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/troian/golang-cross:${GOLANG_CROSS_VERSION} \
		--rm-dist --skip-validate --skip-publish

release:
	@if [ ! -f ".release-env" ]; then \
		echo "\033[91m.release-env is required for release\033[0m";\
		exit 1;\
	fi
	docker run \
		--rm \
		--privileged \
		-e CGO_ENABLED=1 \
		--env-file .release-env \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/troian/golang-cross:${GOLANG_CROSS_VERSION} \
		release --rm-dist --skip-validate

.PHONY: release-dry-run release

###############################################################################
###                          Tools & Dependencies                           ###
###############################################################################

TOOLS_DESTDIR  ?= $(GOPATH)/bin
STATIK         = $(TOOLS_DESTDIR)/statik
RUNSIM         = $(TOOLS_DESTDIR)/runsim

# Install the runsim binary with a temporary workaround of entering an outside
# directory as the "go get" command ignores the -mod option and will polute the
# go.{mod, sum} files.
#
# ref: https://github.com/golang/go/issues/30515
runsim: $(RUNSIM)
$(RUNSIM):
	@echo "Installing runsim..."
	@(cd /tmp && ${GO_MOD} go install github.com/cosmos/tools/cmd/runsim@master)

statik: $(STATIK)
$(STATIK):
	@echo "Installing statik..."
	@(cd /tmp && go get github.com/rakyll/statik@v0.1.6)

contract-tools:
ifeq (, $(shell which stringer))
	@echo "Installing stringer..."
	@go get golang.org/x/tools/cmd/stringer
else
	@echo "stringer already installed; skipping..."
endif

ifeq (, $(shell which go-bindata))
	@echo "Installing go-bindata..."
	@go get github.com/kevinburke/go-bindata/go-bindata
else
	@echo "go-bindata already installed; skipping..."
endif

ifeq (, $(shell which gencodec))
	@echo "Installing gencodec..."
	@go get github.com/fjl/gencodec
else
	@echo "gencodec already installed; skipping..."
endif

ifeq (, $(shell which protoc-gen-go))
	@echo "Installing protoc-gen-go..."
	@go get github.com/fjl/gencodec github.com/golang/protobuf/protoc-gen-go
else
	@echo "protoc-gen-go already installed; skipping..."
endif

ifeq (, $(shell which solcjs))
	@echo "Installing solcjs..."
	@npm install -g solc@0.5.11
else
	@echo "solcjs already installed; skipping..."
endif

tools: tools-stamp
tools-stamp: contract-tools proto-tools statik runsim
	# Create dummy file to satisfy dependency and avoid
	# rebuilding when this Makefile target is hit twice
	# in a row.
	touch $@

tools-clean:
	rm -f $(RUNSIM)
	rm -f tools-stamp

.PHONY: runsim statik tools contract-tools proto-tools  tools-stamp tools-clean

go.sum: go.mod
	echo "Ensure dependencies have not been modified ..." >&2
	go mod verify
	go mod tidy

###############################################################################
###                              Documentation                              ###
###############################################################################

update-swagger-docs: statik
	$(BINDIR)/statik -src=client/docs/swagger-ui -dest=client/docs -f -m
	@if [ -n "$(git status --porcelain)" ]; then \
        echo "\033[91mSwagger docs are out of sync!!!\033[0m";\
        exit 1;\
    else \
        echo "\033[92mSwagger docs are in sync\033[0m";\
    fi
.PHONY: update-swagger-docs

godocs:
	@echo "--> Wait a few seconds and visit http://localhost:6060/pkg/github.com/evmos/ethermint/types"
	godoc -http=:6060

###############################################################################
###                           Tests & Simulation                            ###
###############################################################################

test:
	@go test -mod=readonly $(PACKAGES_NOSIMULATION_TEST)

###############################################################################
###                                Linting                                  ###
###############################################################################

lint:
	@@test -n "$$golangci-lint version | awk '$4 >= 1.42')"
	golangci-lint run --out-format=tab -n
run-integration-tests:
	@nix-shell ./tests/integration_tests/shell.nix --run ./scripts/run-integration-tests.sh

lint-py:
	flake8 --show-source --count --statistics \
          --format="::error file=%(path)s,line=%(row)d,col=%(col)d::%(path)s:%(row)d:%(col)d: %(code)s %(text)s" \

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/docs/statik/statik.go" -not -name '*.pb.go' -not -name '*.pb.gw.go' | xargs gofumpt -d -e -extra

lint-fix:
	golangci-lint run --fix --out-format=tab --issues-exit-code=0
.PHONY: lint lint-fix lint-py

format-fix:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/docs/statik/statik.go" -not -name '*.pb.go' -not -name '*.pb.gw.go' | xargs gofumpt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/docs/statik/statik.go" -not -name '*.pb.go' -not -name '*.pb.gw.go' | xargs misspell -w
.PHONY: format

###############################################################################
###                                Protobuf                                 ###
###############################################################################

protoVer=v0.2
protoImageName=tendermintdev/sdk-proto-gen:$(protoVer)
containerProtoGen=$(PROJECT_NAME)-proto-gen-$(protoVer)
containerProtoGenAny=$(PROJECT_NAME)-proto-gen-any-$(protoVer)
containerProtoGenSwagger=$(PROJECT_NAME)-proto-gen-swagger-$(protoVer)
containerProtoFmt=$(PROJECT_NAME)-proto-fmt-$(protoVer)

proto-all: proto-format proto-lint proto-gen

proto-gen:
	@echo "Generating Protobuf files"
	@if docker ps -a --format '{{.Names}}' | grep -Eq "^${containerProtoGen}$$"; then docker start -a $(containerProtoGen); else docker run --name $(containerProtoGen) -v $(CURDIR):/workspace --workdir /workspace $(protoImageName) \
		sh ./scripts/protocgen.sh; fi

proto-swagger-gen:
	@echo "Generating Protobuf Swagger"
	@./scripts/proto-tools-installer.sh
	@./scripts/protoc-swagger-gen.sh

proto-format:
	@echo "Formatting Protobuf files"
	@if docker ps -a --format '{{.Names}}' | grep -Eq "^${containerProtoFmt}$$"; then docker start -a $(containerProtoFmt); else docker run --name $(containerProtoFmt) -v $(CURDIR):/workspace --workdir /workspace tendermintdev/docker-build-proto \
		find ./ -not -path "./third_party/*" -name "*.proto" -exec clang-format -i {} \; ; fi

proto-lint:
	@$(DOCKER_BUF) lint --error-format=json

proto-check-breaking:
	@$(DOCKER_BUF) breaking --against $(HTTPS_GIT)#branch=main


TM_URL              = https://raw.githubusercontent.com/tendermint/tendermint/v0.34.12/proto/tendermint
GOGO_PROTO_URL      = https://raw.githubusercontent.com/regen-network/protobuf/cosmos
COSMOS_SDK_URL      = https://raw.githubusercontent.com/cosmos/cosmos-sdk/v0.43.0
COSMOS_PROTO_URL    = https://raw.githubusercontent.com/regen-network/cosmos-proto/master

TM_CRYPTO_TYPES     = third_party/proto/tendermint/crypto
TM_ABCI_TYPES       = third_party/proto/tendermint/abci
TM_TYPES            = third_party/proto/tendermint/types

GOGO_PROTO_TYPES    = third_party/proto/gogoproto

COSMOS_PROTO_TYPES  = third_party/proto/cosmos_proto

proto-update-deps:
	@mkdir -p $(GOGO_PROTO_TYPES)
	@curl -sSL $(GOGO_PROTO_URL)/gogoproto/gogo.proto > $(GOGO_PROTO_TYPES)/gogo.proto

	@mkdir -p $(COSMOS_PROTO_TYPES)
	@curl -sSL $(COSMOS_PROTO_URL)/cosmos.proto > $(COSMOS_PROTO_TYPES)/cosmos.proto

## Importing of tendermint protobuf definitions currently requires the
## use of `sed` in order to build properly with cosmos-sdk's proto file layout
## (which is the standard Buf.build FILE_LAYOUT)
## Issue link: https://github.com/tendermint/tendermint/issues/5021
	@mkdir -p $(TM_ABCI_TYPES)
	@curl -sSL $(TM_URL)/abci/types.proto > $(TM_ABCI_TYPES)/types.proto

	@mkdir -p $(TM_TYPES)
	@curl -sSL $(TM_URL)/types/types.proto > $(TM_TYPES)/types.proto

	@mkdir -p $(TM_CRYPTO_TYPES)
	@curl -sSL $(TM_URL)/crypto/proof.proto > $(TM_CRYPTO_TYPES)/proof.proto
	@curl -sSL $(TM_URL)/crypto/keys.proto > $(TM_CRYPTO_TYPES)/keys.proto



.PHONY: proto-all proto-gen proto-gen-any proto-swagger-gen proto-format proto-lint proto-check-breaking proto-update-deps

###############################################################################
###                                Localnet                                 ###
###############################################################################

# Build image for a local testnet
localnet-build:
	@$(MAKE) -C networks/local

# Start a 4-node testnet locally
localnet-start: localnet-stop
ifeq ($(OS),Windows_NT)
	mkdir localnet-setup &
	@$(MAKE) localnet-build

	IF not exist "build/node0/$(TREASURENET_BINARY)/config/genesis.json" docker run --rm -v $(CURDIR)/build\treasurenet\Z treasurenetd/node "./treasurenetd testnet init-files --v 4 -o /treasurenet --keyring-backend=test  --starting-ip-address 192.167.10.2"
	docker-compose up -d
else
	mkdir -p localnet-setup
	@$(MAKE) localnet-build

	if ! [ -f localnet-setup/node0/$(TREASURENET_BINARY)/config/genesis.json ]; then docker run --rm -v $(CURDIR)/localnet-setup:/treasurenet:Z treasurenetd/node "./treasurenetd testnet init-files --v 4 -o /treasurenet --keyring-backend=test  --starting-ip-address 192.167.10.2"; fi
	docker-compose up -d
endif

# Stop testnet
localnet-stop:
	docker-compose down

# Clean testnet
localnet-clean:
	docker-compose down
	sudo rm -rf localnet-setup

 # Reset testnet
localnet-unsafe-reset:
	docker-compose down
ifeq ($(OS),Windows_NT)
	@docker run --rm -v $(CURDIR)\localnet-setup\node0\treasurenetd:treasurenet\Z treasurenetd/node "./treasurenetd unsafe-reset-all --home=/treasurenet"
	@docker run --rm -v $(CURDIR)\localnet-setup\node1\treasurenetd:treasurenet\Z treasurenetd/node "./treasurenetd unsafe-reset-all --home=/treasurenet"
	@docker run --rm -v $(CURDIR)\localnet-setup\node2\treasurenetd:treasurenet\Z treasurenetd/node "./treasurenetd unsafe-reset-all --home=/treasurenet"
	@docker run --rm -v $(CURDIR)\localnet-setup\node3\treasurenetd:treasurenet\Z treasurenetd/node "./treasurenetd unsafe-reset-all --home=/treasurenet"
else
	@docker run --rm -v $(CURDIR)/localnet-setup/node0/treasurenetd:/treasurenet:Z treasurenetd/node "./treasurenetd unsafe-reset-all --home=/treasurenet"
	@docker run --rm -v $(CURDIR)/localnet-setup/node1/treasurenetd:/treasurenet:Z treasurenetd/node "./treasurenetd unsafe-reset-all --home=/treasurenet"
	@docker run --rm -v $(CURDIR)/localnet-setup/node2/treasurenetd:/treasurenet:Z treasurenetd/node "./treasurenetd unsafe-reset-all --home=/treasurenet"
	@docker run --rm -v $(CURDIR)/localnet-setup/node3/treasurenetd:/treasurenet:Z treasurenetd/node "./treasurenetd unsafe-reset-all --home=/treasurenet"
endif

# Clean testnet
localnet-show-logstream:
	docker-compose logs --tail=1000 -f

.PHONY: build-docker-local-treasurenet localnet-start localnet-stop
verify-arch:
	@for arch in $(SUPPORTED_ARCHS); do \
		if [ -f "$(BUILDDIR)/$${arch}/$(TREASURENET_BINARY)" ]; then \
			echo "Verifying $${arch} binary:"; \
			file "$(BUILDDIR)/$${arch}/$(TREASURENET_BINARY)" | grep -E 'ARM|x86'; \
		fi \
	done

test-arm:
	@if [ ! -f "$(BUILDDIR)/arm/$(TREASURENET_BINARY)" ]; then \
		echo "ARM binary not found, building first..."; \
		$(MAKE) build-arm; \
	fi
	qemu-arm -L /usr/arm-linux-gnueabihf $(BUILDDIR)/arm/$(TREASURENET_BINARY) version