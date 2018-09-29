.PHONY: test build clean check
check:
	hack/verify.sh
build:
	scripts/build-go.sh
local-run:
	scripts/run-local-go.sh
image:
	hack/build-image.sh devops/apiserver:0.1
clean:
	rm -rf cmd/server
