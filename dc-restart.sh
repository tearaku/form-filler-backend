#!/bin/bash

docker-compose down
docker image rm form-filler-backend
docker-compose up
