#!/bin/bash

PROJECT1="./src/location_history"
PROJECT2="./src/users"

source ./.env

echo "Starting Project 1..."
(cd "$PROJECT1" && go run ./main &)  

echo "Starting Project 2..."
(cd "$PROJECT2" && go run ./main &)  

read -p "Press any key to stop running projects..."

trap 'kill $(jobs -p)' EXIT INT

echo "Stopping projects..."