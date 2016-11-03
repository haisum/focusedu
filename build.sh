#!/bin/bash
## build.sh compiles and builds focusedu for distribution

if [[ -f focusedu.zip ]]; then
	rm focusedu.zip
fi

go build -o focusedu.exe utils/focusedu/*.go
zip -r focusedu.zip templates static focusedu.exe