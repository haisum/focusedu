#!/bin/bash
## release.sh zips and uploads latest build zip to github releases

TAG=$1
MSG=$2

if [[ -z "$TAG" || -z "$MSG" ]]; then
	echo "TAG and MSG can't be empty"
fi

sh build.sh
github-release.exe -v release --repo focusedu --user haisum -t "$TAG" -d "$MSG" 
github-release.exe -v upload --repo focusedu --user haisum -t "$TAG" -n focusedu.zip -f focusedu.zip