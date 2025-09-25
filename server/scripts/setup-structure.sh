#!/bin/bash

# Script to organize project structure
# This script creates the recommended directory structure

# Create deployment configurations
mkdir -p deployments/docker/{auth-service,user-service,ai-service,api-gateway}
mkdir -p deployments/k8s/{auth-service,user-service,ai-service,api-gateway}

# Create shared utilities
mkdir -p shared/utils/{validation,crypto,logging,monitoring}
mkdir -p shared/models
mkdir -p shared/constants

# Create API documentation
mkdir -p docs/api/{auth,user,ai,gateway}

# Create configuration management
mkdir -p configs/{development,staging,production}

# Create scripts for development
mkdir -p scripts/{build,deploy,test,migration}

# Create test utilities
mkdir -p test/{fixtures,mocks,integration}

echo "Directory structure created successfully!"
