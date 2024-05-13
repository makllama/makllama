################################################################################
# ========================== Capture Environment ===============================
# get the repo root and output path
REPO_ROOT:=${CURDIR}
OUT_DIR=$(REPO_ROOT)/bin
################################################################################
# ============================== OPTIONS =======================================
# the output binary name, overridden when cross compiling
DEMO_BINARY_NAME?=demo
################################################################################
# ================================= Building ===================================
# standard "make" target -> builds
.PHONY: all build
all: build
# builds demo, outputs to $(OUT_DIR)
build:
	go build -v -o "$(OUT_DIR)/$(DEMO_BINARY_NAME)" .
################################################################################

.PHONY: image openai localai
image:
	@docker buildx build --platform=linux/amd64,linux/arm64 --push -t r0ckstar/mods:demo -f Dockerfile.mods .
	@docker push r0ckstar/mods:demo

openai:
	@docker run -it --env OPENAI_API_KEY=${OPENAI_API_KEY} r0ckstar/mods:demo bash

localai:
	@docker run -it --net host --mount type=bind,source="$(CURDIR)"/mods.yml,target=/root/.config/mods/mods.yml,readonly r0ckstar/mods:demo bash
