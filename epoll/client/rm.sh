#!/bin/bash

docker rm -f $(docker ps | grep epoll-client | awk '{print $1}')