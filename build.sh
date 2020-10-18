#!/bin/bash

if [ $# -eq 0 ]; then
    echo "Version not provided"
    exit 1
fi

CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build

rm version.json
echo "{\"version\":\"$1\"}" > version.json

sshpass -p "fgldbpass" sftp root@134.209.25.57 <<EOF
rm fgl-client.exe version.json
put /home/tre/Projects/Go/fgl-database/fgl-client/fgl-client.exe
put /home/tre/Projects/Go/fgl-database/fgl-client/version.json
exit
EOF