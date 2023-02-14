BUILD_DIR=./build
NFPM_CONFIG=./packaging/config
build:  clean init
	go build -o $(BUILD_DIR)/v2ray-api cmd/main.go
clean:
		@rm -rf $(BUILD_DIR)

help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

init: ## create base files and directories
	@mkdir -p $(BUILD_DIR)

