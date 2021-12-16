appVersion ?= 0.0.0
appBuildTime ?= $(shell TZ=Asia/Shanghai date "+%F %T GMT%:z")
appGitCommitHash ?= $(shell git rev-parse HEAD)
projectName ?= go_app
appModuleName = github.com/liuguangw/billing_go/services
buildLdFlags =-X $(appModuleName).appVersion=$(appVersion)
buildLdFlags += -X '$(appModuleName).appBuildTime=$(appBuildTime)'
buildLdFlags += -X $(appModuleName).gitCommitHash=$(appGitCommitHash)
CGO_ENABLED ?= 0
GO_BUILD=go build -v -ldflags "-w -s $(buildLdFlags)"
EXTRA_FILES = config.yaml LICENSE README.md
releasePath ?= ./release

define build_app
	mkdir -p $(releasePath)
	echo "build $(projectName)\(linux/$(2)\)"
	@GOOS=linux GOARCH=$(1) $(GO_BUILD) -o $(releasePath)/$(projectName)
	echo "build $(projectName)\(windows/$(2)\)"
	@GOOS=windows GOARCH=$(1) $(GO_BUILD) -o $(releasePath)/$(projectName).exe
	cp $(EXTRA_FILES) $(releasePath)/
	mv $(releasePath) ./$(projectName)-$(2)
	tar -zcvf $(projectName)-$(2).tar.gz ./$(projectName)-$(2)
	rm -rf ./$(projectName)-$(2)
endef

build:
	@$(GO_BUILD) -o $(projectName)

all:
	#build for x32
	$(call build_app,386,x32)
	#build for x64
	$(call build_app,amd64,x64)

clean:
	rm -rf ./*.tar.gz

.PHONY: build all clean
