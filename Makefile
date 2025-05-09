projectName ?= billing
hostOsType := $(shell go env GOHOSTOS)
ifeq ($(hostOsType),windows)
$(error build $(projectName) is not support on $(hostOsType))
endif
export CGO_ENABLED=0
appVersion ?= 1.3.5
appArchList := x32 x64
appBuildTime ?= $(shell TZ=Asia/Shanghai date "+%F %T GMT%:z")
appGitCommitHash ?= $(shell git rev-parse HEAD)
appModuleName := github.com/liuguangw/billing_go/services

# 获取机器信息
builderMachine ?= unknown

osReleasePath = $(wildcard /etc/os-release)
issuePath = $(wildcard /etc/issue.net)

ifneq ($(osReleasePath),)
	builderMachine=$(shell . $(osReleasePath); echo $$PRETTY_NAME)
else ifneq ($(issuePath),)
	builderMachine=$(file < $(issuePath))
else 
	builderMachine = $(hostOsType)
endif

# github
ifneq (${GITHUB_ACTIONS},)
	builderMachine += (GitHub Actions)
endif

buildLdFlags =-X $(appModuleName).appVersion=$(appVersion)
buildLdFlags += -X '$(appModuleName).appBuildTime=$(appBuildTime)'
buildLdFlags += -X $(appModuleName).gitCommitHash=$(appGitCommitHash)
buildLdFlags += -X '$(appModuleName).builderMachine=$(builderMachine)'
GO_BUILD:=go build -ldflags "-w -s $(buildLdFlags)"
EXTRA_FILES := config.yaml LICENSE README.md
releasePath := ./release

# 是否使用upx
useUpx ?= 0
upxBin :=
ifneq ($(useUpx),0)
	upxBin = $(shell which upx)
endif

# release
define release_app
	@echo build for $(2)
	@mkdir -p $(releasePath)
	@echo "build $(projectName) (linux/$(2))"
	@GOOS=linux GOARCH=$(1) $(GO_BUILD) -o $(releasePath)/$(projectName)
	@echo "build $(projectName) (windows/$(2))"
	@GOOS=windows GOARCH=$(1) $(GO_BUILD) -o $(releasePath)/$(projectName).exe
	@cp $(EXTRA_FILES) $(releasePath)/
endef

# 打包
define tar_app
	@mv $(releasePath) ./$(1)
	@tar -zcf $(1).tar.gz $(1)
	@rm -rf ./$(1)
	@echo output path: $(1).tar.gz
endef

build:
	@$(GO_BUILD) -o $(projectName)
	@echo build $(projectName) ok

# x32 x64
$(appArchList):
# call release_app,386,x32
# or
# call release_app,amd64,x64
	$(call release_app,$(subst x64,amd64,$(subst x32,386,$@)),$@)
ifneq ($(upxBin),)
	@$(upxBin) --best $(releasePath)/$(projectName)
	@$(upxBin) --best $(releasePath)/$(projectName).exe
endif
	$(call tar_app,$(projectName)-release-$@)

all:$(appArchList)

clean:
	@rm -rf ./$(projectName)*
	@rm -rf $(releasePath)

.PHONY: build x32 x64 all clean
