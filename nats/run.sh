#!/bin/bash

if [ ! "$(docker network ls | grep nginx-proxy)" ]; then
      echo "Creating nginx-proxy network ..."
        docker network create nginx-proxy
    else
          echo "nginx-proxy network exists."
fi
