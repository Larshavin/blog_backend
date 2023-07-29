#!/bin/sh

# kill -9 `pgrep blog-webserver`

#delete process and remove pid file
echo "kill process"
kill -9 `cat save_pid.txt`

echo "remove pid file"
rm save_pid.txt

#rebuild and run and save pid file

echo "rebuild and run"
go build -o blog-webserver main.go
nohup ./blog-webserver > blog-webserver.log 2>&1 &
echo $! > save_pid.txt

cat save_pid.txt
ps -ef | grep blog-webserver
lsof -i :443