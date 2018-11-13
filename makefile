TAG=$(tag)
TEST?=$$(go list ./... |grep -v 'vendor')

.PHONY: test build clean check build-flyway
check:
	hack/verify.sh
fmt:
	hack/fmt.sh -l
fmt-all:
	hack/fmt.sh -a
build:
	hack/build-go.sh
local-run:
	hack/run-local-go.sh
image:
	hack/build-image.sh $(TAG)
clean:
	rm -rf cmd/server
test:
	go test $(TEST) -v -timeout 120m


build-image-%: ## build docker image
	@if [ "$*" = "latest" ];then \
	docker build -t kubesphere/devops:latest .; \
	docker build -t kubesphere/devops:flyway -f ./pkg/db/Dockerfile ./pkg/db/;\
	elif [ "`echo "$*" | grep -E "^v[0-9]+\.[0-9]+\.[0-9]+"`" != "" ];then \
	docker build -t kubesphere/devops:$* .; \
	docker build -t kubesphere/devops:flyway-$* -f ./pkg/db/Dockerfile ./pkg/db/; \
	fi

push-image-%: ## push docker image
	@if [ "$*" = "latest" ];then \
	docker push kubesphere/devops:latest; \
	docker push kubesphere/devops:flyway; \
	elif [ "`echo "$*" | grep -E "^v[0-9]+\.[0-9]+\.[0-9]+"`" != "" ];then \
	docker push kubesphere/devops:$*; \
	docker push kubesphere/devops:flyway-$*; \
	fi
