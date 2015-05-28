#!/bin/sh
echo "Installing OTA Server"
echo "-----------------------------"
appPath="/var/lib/android_ota_server"
dbPath="/var/lib/android_ota_server/ota.sql"

echo "> Copying binary to /usr/bin"
cp android_ota_server /usr/bin/

echo "> Creating $appPath"
mkdir $appPath
echo "> Touching $dbPath"
touch $dbPath
echo "> Copying static assets to $appPath"
cp -rf ./views $appPath/views
cp -f ./config.yml $appPath/config.yml
mkdir $appPath/builds
mkdir $appPath/builds/deleted
mkdir $appPath/builds/published
echo "> Installing systemd service"
cp android_ota_server.service /usr/lib/systemd/system/android_ota_server.service
systemctl enable android_ota_server.service

echo "> Adding ota_server user"
sudo useradd android_ota_server -s /sbin/nologin

echo "> Setting permissions"
chmod -R 777 $dbPath
chmod -R 777 $appPath
chmod -R 777 $appPath/config.yml
chmod -R 777 $appPath/views
chmod -R 777 $appPath/builds
chmod -R 777 $appPath/builds/*

echo "-----------------------------"
echo "Installed successfully!"
echo "Now run: systemctl start android_ota_server"
