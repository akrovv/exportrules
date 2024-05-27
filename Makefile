PLUGIN_PATH := .

build_plugin:
	go build -buildmode=plugin -o ${PLUGIN_PATH}/exportrules.so cmd/exportrules/main.go

.PHONY: build_plugin