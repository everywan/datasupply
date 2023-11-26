NAME = datahub
PACKAGE = git.in.zhihu.com/antispam/antispam-datahub-core
MAIN = $(PACKAGE)/entry

# BUILD_FLAGS= -ldflags "-s -w" -mod vendor -v -o  bin/$(NAME) entry/main.go

CL_RED  = "\033[0;31m"
CL_BLUE = "\033[0;34m"
CL_GREEN = "\033[0;32m"
CL_ORANGE = "\033[0;33m"
CL_NONE = "\033[0m"

define color_out
	@echo $(1)$(2)$(CL_NONE)
endef

lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1
	$(ENV) $(GOPATH)/bin/golangci-lint run ./... --timeout=5m

test:
	go test -coverprofile=coverage.out -v ./...

push_prepare:
	# @make linux
	# $(call color_out,$(CL_ORANGE),"Build Success.")
	@make lint
	$(call color_out,$(CL_ORANGE),"Lint Success.")
	@go test ./...
	$(call color_out,$(CL_ORANGE),"Test Success.")
	@go generate ./...
	$(call color_out,$(CL_ORANGE),"Generate Success.")
	$(call color_out,$(CL_ORANGE),"Done")

# build:
# 	@go mod vendor
# 	@go build $(BUILD_FLAGS)

# linux:
# 	@go mod vendor
# 	@GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS)

check:
	@make lint
	$(call color_out,$(CL_ORANGE),"Lint Success.")
	@go test ./...
	$(call color_out,$(CL_ORANGE),"Done")

.PHONY: all proto
all:
	build
