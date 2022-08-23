#!/bin/bash -xe

# Read the first parameter into $APP
if [[ "$1" != "" ]]; then
    APP="$1"
else
    echo "Please specify the location of application you are trying to deploy."
    exit 1
fi
# Install the app into the /var/www/ff_src directory
mkdir -p /var/www/ff_src
# Technically I just need the docker-compose & .env file :PP
cp $APP /var/www/ff_src/ff_src.zip
cd /var/www/ff_src
unzip ff_src.zip
rm ff_src.zip

# Installing docker & docker-compose
sudo yum update -y
sudo amazon-linux-extras install docker
pip3 install docker-compose
# Adding group membership for default ec2-user
# : to run all docker commands w/o sudo
sudo usermod -a -G docker ec2-user

# Enable docker service @ AMI boot time
sudo systemctl enable docker
#sudo systemctl enable docker.service
sudo service docker start
#sudo systemctl start docker.service

# BUG: somehow I still had to manually switch usergroup before executing
# docker-compose command, even though I already have the following??
# Change to the newly created "docker" group
newgrp docker
# Build and start backend via docker-compose
docker-compose up
