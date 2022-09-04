ENTRY_FILE = main.go
GO_BUILD = go build -o
OUTPUT_PATH = bin
.PHONY: build

build: buildwindows64 buildwindows32 buildmacarm buildmacamd buildlinux64 buildlinux32

buildwindows64:
	GOOS=windows GOARCH=amd64 ${GO_BUILD} ${OUTPUT_PATH}/windows/jo.exe ${ENTRY_FILE}

buildwindows32:
	GOOS=windows GOARCH=386 ${GO_BUILD} ${OUTPUT_PATH}/windows/jo32.exe ${ENTRY_FILE}

buildmacarm:
	GOOS=darwin GOARCH=arm64 ${GO_BUILD} ${OUTPUT_PATH}/mac/jo ${ENTRY_FILE}

buildmacamd:
	GOOS=darwin GOARCH=amd64 ${GO_BUILD} ${OUTPUT_PATH}/mac/joAMD ${ENTRY_FILE}

buildlinux64:
	GOOS=linux GOARCH=amd64 ${GO_BUILD} ${OUTPUT_PATH}/linux/jo ${ENTRY_FILE}

buildlinux32:
	GOOS=linux GOARCH=386 ${GO_BUILD} ${OUTPUT_PATH}/linux/jo32 ${ENTRY_FILE}