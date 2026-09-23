PLUGIN_ID := com.tabihard.mattermost-plugin-template
DIST_DIR := dist
BUNDLE_DIR := $(DIST_DIR)/$(PLUGIN_ID)
BUNDLE_NAME := $(PLUGIN_ID).tar.gz

GO := go

PLATFORMS := linux-amd64 linux-arm64 darwin-amd64 darwin-arm64 windows-amd64

.PHONY: all build dist clean $(PLATFORMS)

all: dist

build: $(PLATFORMS)

linux-amd64:
	GOOS=linux GOARCH=amd64 $(GO) build -o $(BUNDLE_DIR)/server/dist/plugin-linux-amd64 ./server

linux-arm64:
	GOOS=linux GOARCH=arm64 $(GO) build -o $(BUNDLE_DIR)/server/dist/plugin-linux-arm64 ./server

darwin-amd64:
	GOOS=darwin GOARCH=amd64 $(GO) build -o $(BUNDLE_DIR)/server/dist/plugin-darwin-amd64 ./server

darwin-arm64:
	GOOS=darwin GOARCH=arm64 $(GO) build -o $(BUNDLE_DIR)/server/dist/plugin-darwin-arm64 ./server

windows-amd64:
	GOOS=windows GOARCH=amd64 $(GO) build -o $(BUNDLE_DIR)/server/dist/plugin-windows-amd64.exe ./server

dist: build
	cp plugin.json $(BUNDLE_DIR)/plugin.json
	tar -C $(DIST_DIR) -czvf $(DIST_DIR)/$(BUNDLE_NAME) $(PLUGIN_ID)

clean:
	rm -rf $(DIST_DIR)
