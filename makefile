.PHONY: test build clean check
check:
	hack/verify.sh
build:
	scripts/build-go.sh
local-run:
	scripts/run-local-go.sh
image:
	scripts/build-image.sh
clean:
	rm -rf cmd/server