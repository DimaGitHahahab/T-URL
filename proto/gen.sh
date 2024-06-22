#!/bin/bash
mkdir -p shorteningpb
mkdir -p redirectionpb
mkdir -p storagepb
mkdir -p analyticspb

protoc --go_out=shorteningpb --go-grpc_out=shorteningpb shortening.proto
protoc --go_out=redirectionpb --go-grpc_out=redirectionpb redirection.proto
protoc --go_out=storagepb --go-grpc_out=storagepb storage.proto
protoc --go_out=analyticspb --go-grpc_out=analyticspb analytics.proto

mkdir -p ../shortening-microservice/proto/
mkdir -p ../redirection-microservice/proto/
mkdir -p ../storage-microservice/proto/
mkdir -p ../analytics-microservice/proto/
mkdir -p ../api-gateway/proto/

cp -r shorteningpb ../shortening-microservice/proto/
cp -r redirectionpb ../redirection-microservice/proto/
cp -r storagepb ../storage-microservice/proto/
cp -r analyticspb ../analytics-microservice/proto/

cp -r shorteningpb ../api-gateway/proto/
cp -r redirectionpb ../api-gateway/proto/
cp -r analyticspb ../api-gateway/proto/

cp -r storagepb ../shortening-microservice/proto/
cp -r storagepb ../redirection-microservice/proto/

rm -rf ./shorteningpb
rm -rf ./redirectionpb
rm -rf ./storagepb
rm -rf ./analyticspb