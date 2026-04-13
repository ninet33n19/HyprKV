#!/bin/bash

# Configuration
HOST="127.0.0.1"
PORT=7379
CLIENT_COUNT=10
MESSAGE="Hello from client"

echo "Starting $CLIENT_COUNT clients connecting to $HOST:$PORT..."

# Loop to spawn clients
for i in {1..10}
do
   # Use a subshell to keep the connection open
   (
     echo "$MESSAGE $i"
     # Keep the connection alive for 5 seconds so you can see
     # the 'concurrent_clients' log in your Go server
     sleep 5
   ) | nc $HOST $PORT &

   echo "Client $i launched."
done

# Wait for background jobs to finish
wait
echo "All clients disconnected."
