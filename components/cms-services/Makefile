DEPLOYS := $(wildcard deploy/*/)
DOCKER_FILES := $(patsubst deploy/%/,deploy/%/Dockerfile,$(DEPLOYS))

build-images:
	echo "$(DOCKER_FILES)"