#!/bin/sh
echo "Installing OTA Server"
echo "-----------------------------"
appPath="/var/lib/ota_server"
dbPath="/var/lib/ota_server/ota.sql"

echo "> Copying binary to /usr/bin"
cp ota_server /usr/bin/

echo "> Creating $appPath"
mkdir $appPath
echo "> Touching $dbPath"
touch $dbPath
echo "> Copying static assets to $appPath"
cp -rf ./templates $appPath/templates
# cp -rf ./assets $appPath/assets
mkdir $appPath/builds
echo "> Installing systemd service"
cp lib/ota_server.service /usr/lib/systemd/system/ota_server.service
systemctl enable ota_server.service

echo "> Adding ota_server user"
sudo useradd ota_server -s /sbin/nologin

echo "> Setting permissions"
chmod -R 777 $dbPath
chmod -R 777 $appPath
chmod -R 777 $appPath/templates
# chmod -R 777 $appPath/assets
# chmod -R 777 $appPath/assets/*
chmod -R 777 $appPath/builds
chmod -R 777 $appPath/builds/*

echo "-----------------------------"
echo "Installed successfully!"
echo "Now run: systemctl start ota_server"
