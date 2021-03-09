#!/bin/bash
## This script helps regenerating multiple client instances in different network namespaces using Docker
## This helps to overcome the ephemeral source port limitation
## Usage: ./connect <connections> <number of clients> <server ip>
## Server IP is usually the Docker gateway IP address, which is 172.17.0.1 by default
## Number of clients helps to speed up connections establishment at large scale, in order to make the demo faster

# CONNECTIONS=$1
# [[ CONNECTIONS -e ]] CONNECTIONS=28000

REPLICAS=$1
[ -z "$REPLICAS" ] && REPLICAS=1;

for (( c=0; c<${REPLICAS}; c++ ))
do
    docker run --rm --network="host" --name epoll-client-$c -d client
done