#!/bin/bash
## build.sh compiles and builds focusedu for distribution


go build -o focusedu.exe utils/focusedu/*.go
zip -r focusedu.zip states static focusedu.exe