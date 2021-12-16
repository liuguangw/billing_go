export CGO_ENABLED=0
appVersion ?= 1.3.3
appBuildTime ?= $(shell TZ=Asia/Shanghai date "+%F %T GMT%:z")
appGitCommitHash ?= $(shell git rev-parse HEAD)
projectName ?= billing
appModuleName = github.com/liuguangw/billing_go/services

# 获取机器信息
builderMachine ?= unknown

osReleasePath = $(wildcard /etc/os-release)
issuePath = $(wildcard /etc/issue.net)

ifneq ($(osReleasePath),)
builderMachine=$(shell . $(osReleasePath); echo $$PRETTY_NAME)
else ifneq ($(issuePath),)
builderMachine=$(file < $(issuePath))
else 
builderMachine = $(shell go env GOHOSTOS)
endif

# github
ifneq (${GITHUB_ACTIONS},)
builderMachine += (GitHub Actions)
endif

buildLdFlags =-X $(appModuleName).appVersion=$(appVersion)
buildLdFlags += -X '$(appModuleName).appBuildTime=$(appBuildTime)'
buildLdFlags += -X $(appModuleName).gitCommitHash=$(appGitCommitHash)
buildLdFlags += -X '$(appModuleName).builderMachine=$(builderMachine)'
GO_BUILD=go build -ldflags "-w -s $(buildLdFlags)"
EXTRA_FILES = config.yaml LICENSE README.md
releasePath = ./release

# 如果upx存在，是否使用upx
useUpx ?= 1
upxBin = $(shell which upx)

define release_app
	@echo build for $(2)
	@mkdir -p $(releasePath)
	@echo "build $(projectName) (linux/$(2))"
	@GOOS=linux GOARCH=$(1) $(GO_BUILD) -o $(releasePath)/$(projectName)
	@echo "build $(projectName) (windows/$(2))"
	@GOOS=windows GOARCH=$(1) $(GO_BUILD) -o $(releasePath)/$(projectName).exe
	@if [ $(useUpx) -eq 1 ] && [ -n "$(upxBin)" ]; then \
		$(upxBin) --best $(releasePath)/$(projectName); \
		$(upxBin) --best $(releasePath)/$(projectName).exe; \
	fi
	@cp $(EXTRA_FILES) $(releasePath)/
	@mv $(releasePath) ./$(projectName)-release-$(2)
	@tar -zcf $(projectName)-release-$(2).tar.gz $(projectName)-release-$(2)
	@rm -rf ./$(projectName)-release-$(2)
	@echo output path: $(projectName)-release-$(2).tar.gz
endef

build:
	@$(GO_BUILD) -o $(projectName)
	@echo build $(projectName) ok

x32:
	$(call release_app,386,x32)

x64:
	$(call release_app,amd64,x64)

all:x32 x64

clean:
	@rm -rf ./billing*

.PHONY: build x32 x64 all clean
