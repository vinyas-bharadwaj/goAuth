#!/bin/bash

# Generate Python gRPC client code from proto files
# Run this from the root directory

echo "Generating Python gRPC client code..."

# Create output directory for Python proto files
mkdir -p services/api-gateway/proto

# Generate Python gRPC code
python -m grpc_tools.protoc \
    -I shared/proto \
    --python_out=services/api-gateway/proto \
    --grpc_python_out=services/api-gateway/proto \
    shared/proto/auth.proto

echo "Python gRPC client code generated in services/api-gateway/proto/"
