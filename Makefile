name := $(shell basename "$(CURDIR)")
registry := $GCR_REGISTRY

ifdef CI
git_hash := $(shell echo ${CI_COMMIT_SHA} | cut -c1-10 )
git_branch := $(or ${CI_COMMIT_REF_NAME}, unknown)
git_tag := $(or ${CI_COMMIT_TAG}, ${CI_COMMIT_REF_NAME}, unknown)
else
git_hash := $(shell git rev-parse HEAD | cut -c1-10)
git_branch := $(shell git rev-parse --abbrev-ref HEAD || echo "unknown")
git_tag := $(or ${git_branch}, unknown)
endif

build:
	docker build \
		-t \
		$(registry)/$(name):${git_branch}-${git_hash} .

push:
	docker push \
		$(registry)/$(name):${git_branch}-${git_hash}

test:
	go test -count=1 -v ./...