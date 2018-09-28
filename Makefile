TRAG.Org:=kubesphere
TRAG.Name:=ks-devops-apiserver
TRAG.Gopkg:=kubesphere.io/devops

DOCKER_TAGS=latest
RUN_IN_DOCKER:=docker run -it --rm -v `pwd`:/go/src/$(TRAG.Gopkg) -v `pwd`/tmp/cache:/root/.cache/go-build  -w /go/src/$(TRAG.Gopkg) -e GOBIN=/go/src/$(TRAG.Gopkg)/tmp/bin -e USER_ID=`id -u` -e GROUP_ID=`id -g` kubesphere/kubesphere-builder
GO_FMT:=goimports -l -w -e -local=kubesphere -srcdir=/go/src/$(TRAG.Gopkg)
GO_FILES:=./cmd ./pkg

define get_diff_files
    $(eval DIFF_FILES=$(shell git diff --name-only --diff-filter=ad | grep -E "^(test|cmd|pkg)/.+\.go"))
endef
define get_build_flags
    $(eval SHORT_VERSION=$(shell git describe --tags --always --dirty="-dev"))
    $(eval SHA1_VERSION=$(shell git show --quiet --pretty=format:%H))
	$(eval DATE=$(shell date +'%Y-%m-%dT%H:%M:%S'))
	$(eval BUILD_FLAG= -X $(TRAG.Version).ShortVersion="$(SHORT_VERSION)" \
		-X $(TRAG.Version).GitSha1Version="$(SHA1_VERSION)" \
		-X $(TRAG.Version).BuildDate="$(DATE)")
endef


.PHONY: fmt-all
fmt-all:
	$(RUN_IN_DOCKER) $(GO_FMT) $(GO_FILES)
	@echo "fmt done"

.PHONY: fmt
fmt:
	$(call get_diff_files)
	$(if $(DIFF_FILES), \
		$(RUN_IN_DOCKER) $(GO_FMT) ${DIFF_FILES}, \
		$(info cannot find modified files from git) \
	)
	@echo "fmt done"

.PHONY: fmt-check
fmt-check: fmt-all
	$(call get_diff_files)
	$(if $(DIFF_FILES), \
		exit 2 \
	)
