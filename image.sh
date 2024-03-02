#!/bin/sh

docker build . -t back/back-app:latest
docker tag back/back-app:latest harbor.3trolls.me/blog/back/back-app:latest
docker push harbor.3trolls.me/blog/back/back-app:latest
