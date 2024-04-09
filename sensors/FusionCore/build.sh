#!/bin/bash

install_dir="/usr/local/bin/"

sudo docker pull influxdb
sudo docker pull postgres


# Make directories for discern & postgres
mkdir -p /etc/discern
mkdir -p /etc/discern/postgres
mkdir -p /etc/discern/postgres/schema
# Copy postgres information over
cp ./databases/postgres/schema/* /etc/discern/postgres/schema
# cp ./databases/postgres/password /etc/discern/postgres/password
cp ./databases/postgres/Dockerfile /etc/discern/postgres/Dockerfile
# Copy config over
cp ./CoreConfig.yaml /etc/discern/CoreConfig.yaml

# Create postgres docker image
pushd .
cd /etc/discern/postgres
docker buildx build -t discernpsql .
popd


go mod tidy
go build .
mv FusionCore "$install_dir/FusionCore"

cp ./FusionCore.service /etc/systemd/system/
sudo systemctl daemon-reload

sudo systemctl enable FusionCore.service
sudo systemctl start  FusionCore.service
