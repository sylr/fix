GO                ?= $(shell which go)
GIT_UPDATE_INDEX  := $(shell git update-index --refresh)
GIT_REVISION      ?= $(shell git rev-parse HEAD)
GIT_VERSION       ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)

GOENV_GOOS               := $(shell $(GO) env GOOS)
GOENV_GOARCH             := $(shell $(GO) env GOARCH)
GOENV_GOARM              := $(shell $(GO) env GOARM)
GOOS                     ?= $(GOENV_GOOS)
GOARCH                   ?= $(GOENV_GOARCH)
GOARM                    ?= $(GOENV_GOARM)
GO_BUILD_SRC             := $(shell find . -name \*.go -type f) go.mod go.sum
GO_BUILD_EXTLDFLAGS      :=
GO_BUILD_TAGS            ?= acceptor
GO_BUILD_TARGET_DEPS     :=
GO_BUILD_FLAGS           := -trimpath
GO_BUILD_LDFLAGS_OPTIMS  :=

ifeq ($(GOOS)/$(GOARCH),$(GOENV_GOOS)/$(GOENV_GOARCH))
GO_BUILD_TARGET          ?= dist/fix
GO_BUILD_VERSION_TARGET  ?= dist/fix-$(GIT_VERSION)
else
ifeq ($(GOARCH),arm)
GO_BUILD_TARGET          ?= dist/fix-$(GOOS)-$(GOARCH)v$(GOARM)
GO_BUILD_VERSION_TARGET  := dist/fix-$(GIT_VERSION)-$(GOOS)-$(GOARCH)v$(GOARM)
else
GO_BUILD_TARGET          ?= dist/fix-$(GOOS)-$(GOARCH)
GO_BUILD_VERSION_TARGET  := dist/fix-$(GIT_VERSION)-$(GOOS)-$(GOARCH)
endif # ($(GOARCH),arm)
endif # ($(GOOS)/$(GOARCH),$(GOENV_GOOS)/$(GOENV_GOARCH))

ifneq ($(DEBUG),)
GO_BUILD_FLAGS           += -race -gcflags="all=-N -l"
else
GO_BUILD_LDFLAGS_OPTIMS  += -s -w
endif # $(DEBUG)

GO_BUILD_FLAGS_TARGET           := .go-build-flags
GO_CROSSBUILD_WINDOWS_PLATFORMS := windows/amd64 windows/arm64
GO_CROSSBUILD_PLATFORMS         ?= linux/amd64 linux/arm64 \
                                   linux/mips/softfloat linux/mips64 linux/mips64le \
                                   linux/riscv64 linux/s390x \
                                   freebsd/amd64 freebsd/arm64 \
                                   openbsd/amd64 openbsd/arm64 openbsd/mips64 \
                                   netbsd/amd64 netbsd/arm64 \
                                   plan9/amd64 \
                                   darwin/amd64 darwin/arm64 \
                                   dragonfly/amd64 illumos/amd64 solaris/amd64

GO_CROSSBUILD_AMD64_PLATFORMS         := $(filter %/amd64,$(GO_CROSSBUILD_PLATFORMS))
GO_CROSSBUILD_ARM_PLATFORMS           := $(filter %/arm,$(GO_CROSSBUILD_PLATFORMS))
GO_CROSSBUILD_ARM64_PLATFORMS         := $(filter %/arm64,$(GO_CROSSBUILD_PLATFORMS))
GO_CROSSBUILD_ARMV6_PLATFORMS         := $(filter %/arm/v6,$(GO_CROSSBUILD_PLATFORMS))
GO_CROSSBUILD_ARMV7_PLATFORMS         := $(filter %/arm/v7,$(GO_CROSSBUILD_PLATFORMS))
GO_CROSSBUILD_MIPS_PLATFORMS          := $(filter %/mips,$(GO_CROSSBUILD_PLATFORMS))
GO_CROSSBUILD_MIPSLE_PLATFORMS        := $(filter %/mipsle,$(GO_CROSSBUILD_PLATFORMS))
GO_CROSSBUILD_MIPS64_PLATFORMS        := $(filter %/mips64,$(GO_CROSSBUILD_PLATFORMS))
GO_CROSSBUILD_MIPS64LE_PLATFORMS      := $(filter %/mips64le,$(GO_CROSSBUILD_PLATFORMS))
GO_CROSSBUILD_MIPSSOFTFLOAT_PLATFORMS := $(filter %/mips/softfloat,$(GO_CROSSBUILD_PLATFORMS))
GO_CROSSBUILD_RISCV64_PLATFORMS       := $(filter %/riscv64,$(GO_CROSSBUILD_PLATFORMS))
GO_CROSSBUILD_S390X_PLATFORMS         := $(filter %/s390x,$(GO_CROSSBUILD_PLATFORMS))
GO_CROSSBUILD_WINDOWS_PLATFORMS       := $(filter windows/%,$(GO_CROSSBUILD_WINDOWS_PLATFORMS))

GO_CROSSBUILD_AMD64_TARGET_PATTERN    := dist/fix-$(GIT_VERSION)-%-amd64
GO_CROSSBUILD_ARM_TARGET_PATTERN      := dist/fix-$(GIT_VERSION)-%-arm
GO_CROSSBUILD_ARM64_TARGET_PATTERN    := dist/fix-$(GIT_VERSION)-%-arm64
GO_CROSSBUILD_ARMV6_TARGET_PATTERN    := dist/fix-$(GIT_VERSION)-%-armv6
GO_CROSSBUILD_ARMV7_TARGET_PATTERN    := dist/fix-$(GIT_VERSION)-%-armv7
GO_CROSSBUILD_MIPS_TARGET_PATTERN     := dist/fix-$(GIT_VERSION)-%-mips
GO_CROSSBUILD_MIPSLE_TARGET_PATTERN   := dist/fix-$(GIT_VERSION)-%-mipsle
GO_CROSSBUILD_MIPS64_TARGET_PATTERN   := dist/fix-$(GIT_VERSION)-%-mips64
GO_CROSSBUILD_MIPS64LE_TARGET_PATTERN := dist/fix-$(GIT_VERSION)-%-mips64le
GO_CROSSBUILD_RISCV64_TARGET_PATTERN  := dist/fix-$(GIT_VERSION)-%-riscv64
GO_CROSSBUILD_S390X_TARGET_PATTERN    := dist/fix-$(GIT_VERSION)-%-s390x
GO_CROSSBUILD_WINDOWS_TARGET_PATTERN  := dist/fix-$(GIT_VERSION)-windows-%.exe

GO_CROSSBUILD_TARGETS += $(patsubst %/amd64,$(GO_CROSSBUILD_AMD64_TARGET_PATTERN),$(GO_CROSSBUILD_AMD64_PLATFORMS))
GO_CROSSBUILD_TARGETS += $(patsubst %/arm,$(GO_CROSSBUILD_ARM_TARGET_PATTERN),$(GO_CROSSBUILD_ARM_PLATFORMS))
GO_CROSSBUILD_TARGETS += $(patsubst %/arm64,$(GO_CROSSBUILD_ARM64_TARGET_PATTERN),$(GO_CROSSBUILD_ARM64_PLATFORMS))
GO_CROSSBUILD_TARGETS += $(patsubst %/arm/v6,$(GO_CROSSBUILD_ARMV6_TARGET_PATTERN),$(GO_CROSSBUILD_ARMV6_PLATFORMS))
GO_CROSSBUILD_TARGETS += $(patsubst %/arm/v7,$(GO_CROSSBUILD_ARMV7_TARGET_PATTERN),$(GO_CROSSBUILD_ARMV7_PLATFORMS))
GO_CROSSBUILD_TARGETS += $(patsubst %/mips,$(GO_CROSSBUILD_MIPS_TARGET_PATTERN),$(GO_CROSSBUILD_MIPS_PLATFORMS))
GO_CROSSBUILD_TARGETS += $(patsubst %/mipsle,$(GO_CROSSBUILD_MIPSLE_TARGET_PATTERN),$(GO_CROSSBUILD_MIPSLE_PLATFORMS))
GO_CROSSBUILD_TARGETS += $(patsubst %/mips64,$(GO_CROSSBUILD_MIPS64_TARGET_PATTERN),$(GO_CROSSBUILD_MIPS64_PLATFORMS))
GO_CROSSBUILD_TARGETS += $(patsubst %/mips64le,$(GO_CROSSBUILD_MIPS64LE_TARGET_PATTERN),$(GO_CROSSBUILD_MIPS64LE_PLATFORMS))
GO_CROSSBUILD_TARGETS += $(patsubst %/mips/softfloat,$(GO_CROSSBUILD_LINUX_MIPS_TARGET_PATTERN),$(GO_CROSSBUILD_LINUX_MIPS_PLATFORMS))
GO_CROSSBUILD_TARGETS += $(patsubst %/riscv64,$(GO_CROSSBUILD_RISCV64_TARGET_PATTERN),$(GO_CROSSBUILD_RISCV64_PLATFORMS))
GO_CROSSBUILD_TARGETS += $(patsubst %/s390x,$(GO_CROSSBUILD_S390X_TARGET_PATTERN),$(GO_CROSSBUILD_S390X_PLATFORMS))
GO_CROSSBUILD_TARGETS += $(patsubst windows/%,$(GO_CROSSBUILD_WINDOWS_TARGET_PATTERN),$(GO_CROSSBUILD_WINDOWS_PLATFORMS))

GO_BUILD_EXTLDFLAGS     := $(strip $(GO_BUILD_EXTLDFLAGS))
GO_BUILD_TAGS           := $(strip $(GO_BUILD_TAGS))
GO_BUILD_TARGET_DEPS    := $(strip $(GO_BUILD_TARGET_DEPS))
GO_BUILD_FLAGS          := $(strip $(GO_BUILD_FLAGS))
GO_BUILD_LDFLAGS_OPTIMS := $(strip $(GO_BUILD_LDFLAGS_OPTIMS))
GO_BUILD_LDFLAGS        := -ldflags '$(GO_BUILD_LDFLAGS_OPTIMS) -X sylr.dev/fix/cmd.Version=$(GIT_VERSION) -extldflags "$(GO_BUILD_EXTLDFLAGS)"'

GO_TOOLS_GOLANGCI_LINT ?= $(shell $(GO) env GOPATH)/bin/golangci-lint

DOCKER_BUILD_IMAGE      ?= ghcr.io/sylr/fix
DOCKER_BUILD_VERSION    ?= $(GIT_VERSION)
DOCKER_BUILD_GO_VERSION ?= 1.17
DOCKER_BUILD_LABELS      = --label org.opencontainers.image.title=fix
DOCKER_BUILD_LABELS     += --label org.opencontainers.image.description="fix client"
DOCKER_BUILD_LABELS     += --label org.opencontainers.image.url="https://github.com/sylr/fix"
DOCKER_BUILD_LABELS     += --label org.opencontainers.image.source="https://github.com/sylr/fix"
DOCKER_BUILD_LABELS     += --label org.opencontainers.image.revision=$(GIT_REVISION)
DOCKER_BUILD_LABELS     += --label org.opencontainers.image.version=$(GIT_VERSION)
DOCKER_BUILD_LABELS     += --label org.opencontainers.image.created=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
DOCKER_BUILD_BUILD_ARGS ?= --build-arg=GO_VERSION=$(DOCKER_BUILD_GO_VERSION)
DOCKER_BUILDX_PLATFORMS ?= linux/amd64,linux/arm64
DOCKER_BUILDX_CACHE     ?= /tmp/.buildx-cache

export CGO_ENABLED=1

# ------------------------------------------------------------------------------

.PHONY: all build crossbuild crossbuild-checksums .FORCE

build: $(GO_BUILD_VERSION_TARGET) $(GO_BUILD_TARGET)

all: crossbuild crossbuild-checksums

install:
	$(GO) install -tags $(GO_BUILD_TAGS) $(GO_BUILD_FLAGS) $(GO_BUILD_LDFLAGS)

$(GO_BUILD_FLAGS_TARGET) : .FORCE
	@(echo "GO_VERSION=$(shell $(GO) version)"; \
	  echo "GO_GOOS=$(GOOS)"; \
	  echo "GO_GOARCH=$(GOARCH)"; \
	  echo "GO_GOARM=$(GOARM)"; \
	  echo "GO_BUILD_TAGS=$(GO_BUILD_TAGS)"; \
	  echo "GO_BUILD_FLAGS=$(GO_BUILD_FLAGS)"; \
	  echo 'GO_BUILD_LDFLAGS=$(subst ','\'',$(GO_BUILD_LDFLAGS))') \
	    | cmp -s - $@ \
	        || (echo "GO_VERSION=$(shell $(GO) version)"; \
	            echo "GO_GOOS=$(GOOS)"; \
	            echo "GO_GOARCH=$(GOARCH)"; \
	            echo "GO_GOARM=$(GOARM)"; \
	            echo "GO_BUILD_TAGS=$(GO_BUILD_TAGS)"; \
	            echo "GO_BUILD_FLAGS=$(GO_BUILD_FLAGS)"; \
	            echo 'GO_BUILD_LDFLAGS=$(subst ','\'',$(GO_BUILD_LDFLAGS))') > $@

$(GO_BUILD_TARGET): $(GO_BUILD_VERSION_TARGET)
	@(test -e $@ && unlink $@) || true
	@mkdir -p $$(dirname $@)
	@ln $< $@

$(GO_BUILD_VERSION_TARGET): $(GO_BUILD_SRC) $(GO_GENERATE_TARGET) $(GO_BUILD_FLAGS_TARGET) | $(GO_BUILD_TARGET_DEPS)
	GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) $(GO) build -tags $(GO_BUILD_TAGS) $(GO_BUILD_FLAGS) $(GO_BUILD_LDFLAGS) -o $@

crossbuild: $(GO_BUILD_VERSION_TARGET) $(GO_CROSSBUILD_TARGETS)

$(GO_CROSSBUILD_WINDOWS_TARGET_PATTERN): $(GO_BUILD_SRC) $(GO_BUILD_FLAGS_TARGET)
	GOOS=windows GOARCH=$* $(GO) build -tags $(GO_BUILD_TAGS),crossbuild $(GO_BUILD_FLAGS) $(GO_BUILD_LDFLAGS) -o $@

$(GO_CROSSBUILD_386_TARGET_PATTERN): $(GO_BUILD_SRC) $(GO_BUILD_FLAGS_TARGET)
	GOOS=$* GOARCH=386 $(GO) build -tags $(GO_BUILD_TAGS),crossbuild $(GO_BUILD_FLAGS) $(GO_BUILD_LDFLAGS) -o $@

$(GO_CROSSBUILD_AMD64_TARGET_PATTERN): $(GO_BUILD_SRC) $(GO_BUILD_FLAGS_TARGET)
	GOOS=$* GOARCH=amd64 $(GO) build -tags $(GO_BUILD_TAGS),crossbuild $(GO_BUILD_FLAGS) $(GO_BUILD_LDFLAGS) -o $@

$(GO_CROSSBUILD_ARM_TARGET_PATTERN): $(GO_BUILD_SRC) $(GO_BUILD_FLAGS_TARGET)
	GOOS=$* GOARCH=arm $(GO) build -tags $(GO_BUILD_TAGS),crossbuild $(GO_BUILD_FLAGS) $(GO_BUILD_LDFLAGS) -o $@

$(GO_CROSSBUILD_ARM64_TARGET_PATTERN): $(GO_BUILD_SRC) $(GO_BUILD_FLAGS_TARGET)
	GOOS=$* GOARCH=arm64 $(GO) build -tags $(GO_BUILD_TAGS),crossbuild $(GO_BUILD_FLAGS) $(GO_BUILD_LDFLAGS) -o $@

$(GO_CROSSBUILD_ARMV6_TARGET_PATTERN): $(GO_BUILD_SRC) $(GO_BUILD_FLAGS_TARGET)
	GOOS=$* GOARCH=arm GOARM=6 $(GO) build -tags $(GO_BUILD_TAGS),crossbuild $(GO_BUILD_FLAGS) $(GO_BUILD_LDFLAGS) -o $@

$(GO_CROSSBUILD_ARMV7_TARGET_PATTERN): $(GO_BUILD_SRC) $(GO_BUILD_FLAGS_TARGET)
	GOOS=$* GOARCH=arm GOARM=7 $(GO) build -tags $(GO_BUILD_TAGS),crossbuild $(GO_BUILD_FLAGS) $(GO_BUILD_LDFLAGS) -o $@

$(GO_CROSSBUILD_MIPS_TARGET_PATTERN): $(GO_BUILD_SRC) $(GO_BUILD_FLAGS_TARGET)
	GOOS=$* GOARCH=mips $(GO) build -tags $(GO_BUILD_TAGS),crossbuild $(GO_BUILD_FLAGS) $(GO_BUILD_LDFLAGS) -o $@

$(GO_CROSSBUILD_MIPSLE_TARGET_PATTERN): $(GO_BUILD_SRC) $(GO_BUILD_FLAGS_TARGET)
	GOOS=$* GOARCH=mipsle $(GO) build -tags $(GO_BUILD_TAGS),crossbuild $(GO_BUILD_FLAGS) $(GO_BUILD_LDFLAGS) -o $@

$(GO_CROSSBUILD_MIPS64_TARGET_PATTERN): $(GO_BUILD_SRC) $(GO_BUILD_FLAGS_TARGET)
	GOOS=$* GOARCH=mips64 $(GO) build -tags $(GO_BUILD_TAGS),crossbuild $(GO_BUILD_FLAGS) $(GO_BUILD_LDFLAGS) -o $@

$(GO_CROSSBUILD_MIPS64LE_TARGET_PATTERN): $(GO_BUILD_SRC) $(GO_BUILD_FLAGS_TARGET)
	GOOS=$* GOARCH=mips64le $(GO) build -tags $(GO_BUILD_TAGS),crossbuild $(GO_BUILD_FLAGS) $(GO_BUILD_LDFLAGS) -o $@

$(GO_CROSSBUILD_MIPSSOFFLOAT_TARGET_PATTERN): $(GO_BUILD_SRC) $(GO_BUILD_FLAGS_TARGET)
	GOOS=$* GOARCH=mips GOMIPS=softfloat $(GO) build -tags $(GO_BUILD_TAGS),crossbuild $(GO_BUILD_FLAGS) $(GO_BUILD_LDFLAGS) -o $@

$(GO_CROSSBUILD_RISCV64_TARGET_PATTERN): $(GO_BUILD_SRC) $(GO_BUILD_FLAGS_TARGET)
	GOOS=$* GOARCH=riscv64 $(GO) build -tags $(GO_BUILD_TAGS),crossbuild $(GO_BUILD_FLAGS) $(GO_BUILD_LDFLAGS) -o $@

$(GO_CROSSBUILD_S390X_TARGET_PATTERN): $(GO_BUILD_SRC) $(GO_BUILD_FLAGS_TARGET)
	GOOS=$* GOARCH=s390x $(GO) build -tags $(GO_BUILD_TAGS),crossbuild $(GO_BUILD_FLAGS) $(GO_BUILD_LDFLAGS) -o $@

crossbuild-checksums: dist/checksums

dist/checksums : $(GO_CROSSBUILD_TARGETS)
	cd dist && shasum -a 256 fix-*-* > checksums

# -- go mod --------------------------------------------------------------------

.PHONY: go-mod-verify go-mod-tidy

go-mod-verify:
	$(GO) mod download
	git diff --quiet go.* || git diff --exit-code go.* || exit 1

go-mod-tidy:
	$(GO) mod download
	$(GO) mod tidy

# ------------------------------------------------------------------------------

test:
	$(GO) test ./...

lint: $(GO_TOOLS_GOLANGCI_LINT)
	$(GO_TOOLS_GOLANGCI_LINT) run

# -- tools ---------------------------------------------------------------------

.PHONY: tools

tools: $(GO_TOOLS_GOLANGCI_LINT)

$(GO_TOOLS_GOLANGCI_LINT):
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.36.0

# -- docker --------------------------------------------------------------------

.PHONY: docker-buildx-build docker-buildx-push docker-buildx-inspect

docker-buildx-build:
	@docker buildx build . -f Dockerfile \
		-t $(DOCKER_BUILD_IMAGE):$(DOCKER_BUILD_VERSION) \
		--cache-to=type=local,dest=$(DOCKER_BUILDX_CACHE) \
		--platform=$(DOCKER_BUILDX_PLATFORMS) \
		$(DOCKER_BUILD_BUILD_ARGS) \
		$(DOCKER_BUILD_LABELS)

docker-buildx-push:
	@docker buildx build . -f Dockerfile \
		-t $(DOCKER_BUILD_IMAGE):$(DOCKER_BUILD_VERSION) \
		--cache-from=type=local,src=$(DOCKER_BUILDX_CACHE) \
		--platform=$(DOCKER_BUILDX_PLATFORMS) \
		$(DOCKER_BUILD_BUILD_ARGS) \
		$(DOCKER_BUILD_LABELS) \
		--push

docker-buildx-inspect:
	@docker buildx imagetools inspect $(DOCKER_BUILD_IMAGE):$(DOCKER_BUILD_VERSION)
