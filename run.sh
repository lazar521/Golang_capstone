#!/bin/bash

TEST=0

while getopts "t" opt; do
  case ${opt} in
    t )
      TEST=1
      ;;
    * )
      echo "Invalid option $opt"
      exit 1
      ;;
  esac
done


PROJECT1="$(pwd)/src/location_history"
PROJECT2="$(pwd)/src/users"
COMMON="$(pwd)/src/common"

source ./.env

if [ "$TEST" -eq 0 ]; then
    echo "Starting Project 1..."
    (cd "$PROJECT1" && go run ./main &)  

    echo "Starting Project 2..."
    (cd "$PROJECT2" && go run ./main &)  

    read -p "Press any key to stop running projects..."

    trap 'kill $(jobs -p)' EXIT INT

    echo "Stopping projects..."
else
    echo "Testing Project 1..."
    cd "$PROJECT1"
    go test ./main -v

    echo "Testing Project 2..."
    cd "$PROJECT2"
    go test ./main -v

    echo "Testing Common module..."
    cd "$COMMON"
    go test ./utils -v

    echo "Finished..."
fi