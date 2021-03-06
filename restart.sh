#!/bin/bash

cd "/root/go/src/github.com/Zamiell/isaac-racing-server/src"
GOPATH=/root/go /usr/local/go/bin/go install
if [ $? -eq 0 ]; then
        mv "/root/go/bin/src" "/root/go/bin/isaac-racing-server"
	supervisorctl restart isaac-racing-server
else
	echo "isaac-racing-server - Go compilation failed!"
fi
