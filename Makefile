PLUGIN_PATH := .

build_plugin:
	go build -buildmode=plugin -o ${PLUGIN_PATH}/exportes.so cmd/exportes/main.go

.PHONY: build_plugin