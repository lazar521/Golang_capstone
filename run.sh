#!/bin/bash

PROJECT1="./location_history"
PROJECT2="./users"

export DATA_FOLDER="`pwd`/data"

echo "Starting Project 1..."
(cd "$PROJECT1" && go run . &)  

echo "Starting Project 2..."
(cd "$PROJECT2" && go run . &)  

read -p "Press any key to stop running projects..."

trap 'kill $(jobs -p)' EXIT

echo "Stopping projects..."