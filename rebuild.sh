#!/bin/sh

docker build -t backend . --no-cache
docker tag backend:latest harbor.3trolls.me/blog/backend:v0.0.1
docker push harbor.3trolls.me/blog/backend:v0.0.1

kubectl rollout restart deployment blog-deployment
