#!/bin/sh

if [ ! -f .BUILD_ROOT ]; then
    echo "Please run this script from the root of the project."
    exit 1
fi

# Generate the client
openapi-generator-cli generate -i spec/openapi.yaml -g go -o client --additional-properties=packageName=client

# Generate the server
openapi-generator-cli generate -i spec/openapi.yaml -g go-server -o server --additional-properties=addResponseHeaders=true,outputAsLibrary=true,packageName=server

# Copy the spec to the swagger docs
cp spec/openapi.yaml swagger/problematic_api.yaml
