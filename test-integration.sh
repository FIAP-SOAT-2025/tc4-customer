#!/bin/bash

# Script to run integration tests with MongoDB in Docker

echo "Starting MongoDB container for testing..."
docker run -d --name mongodb-test -p 27018:27017 mongo:7.0

# Wait for MongoDB to be ready
echo "Waiting for MongoDB to be ready..."
sleep 5

# Run integration tests
echo "Running integration tests..."
TEST_MONGODB_URI="mongodb://localhost:27018" go test -v -tags=integration ./internal/repository/

# Capture test exit code
TEST_EXIT_CODE=$?

# Cleanup
echo "Cleaning up..."
docker stop mongodb-test
docker rm mongodb-test

exit $TEST_EXIT_CODE
