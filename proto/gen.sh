#!/bin/bash
mkdir -p shortening
mkdir -p redirection
mkdir -p storage
mkdir -p analytics

protoc --go_out=shortening --go-grpc_out=shortening shortening.proto
protoc --go_out=redirection --go-grpc_out=redirection redirection.proto
protoc --go_out=storage --go-grpc_out=storage storage.proto
protoc --go_out=analytics --go-grpc_out=analytics analytics.proto

mkdir -p ../shortening-microservice/proto/
mkdir -p ../redirection-microservice/proto/
mkdir -p ../storage-microservice/proto/
mkdir -p ../analytics-microservice/proto/
mkdir -p ../api-gateway/proto/

cp -r shortening ../shortening-microservice/proto/
cp -r redirection ../redirection-microservice/proto/
cp -r storage ../storage-microservice/proto/
cp -r analytics ../analytics-microservice/proto/

cp -r shortening ../api-gateway/proto/
cp -r redirection ../api-gateway/proto/
cp -r analytics ../api-gateway/proto/