TAG=$(tag)

.PHONY: test build clean check
check:
	hack/verify.sh
fmt:
	hack/fmt.sh -l
fmt-all:
	hack/fmt.sh -a
build:
	scripts/build-go.sh
local-run:
	scripts/run-local-go.sh
image:
	hack/build-image.sh $(TAG)
clean:
	rm -rf cmd/server

