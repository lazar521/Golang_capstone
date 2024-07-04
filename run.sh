#!/bin/bash

PROJECT1="./location_history"
PROJECT2="./users"

export DATA_FOLDER="$(pwd)/data"

export LOCATION_HISTORY_REST_HOST="localhost"
export LOCATION_HISTORY_REST_PORT="8000"
export LOCATION_HISTORY_GRPC_HOST="localhost"
export LOCATION_HISTORY_GRPC_PORT="50051"

export USERS_REST_HOST="localhost"
export USERS_REST_PORT="8001"
export USERS_GRPC_HOST="localhost"
export USERS_GRPC_PORT="50051"

export USERS_DATABASE_URL="$(pwd)/data/users.db"
export LOCATION_HISTORY_DATABASE_URL="$(pwd)/data/location_history.db"

echo "Starting Project 1..."
(cd "$PROJECT1" && go run . &)  

echo "Starting Project 2..."
(cd "$PROJECT2" && go run . &)  

read -p "Press any key to stop running projects..."

trap 'kill $(jobs -p)' EXIT

echo "Stopping projects..."