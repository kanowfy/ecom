#!/usr/bin/bash

docker compose down
docker rmi -f $(docker image ls -q ecom*)
