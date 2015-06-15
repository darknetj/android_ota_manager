#!/bin/sh
#
# Installs android_ota_server & preps directories
#

echo "Installing Android OTA Server"
echo "-------------------------------------"

if [ `whoami` != "root" ]; then
	echo "This script needs to be run as root!"
	exit 1
fi

appPath="/var/lib/android_ota_server"
dbPath="/var/lib/android_ota_server/ota.sql"
buildsPath="/home/storage"

echo "> Copying binary to /usr/local/bin"
cp android_ota_server /usr/local/bin/

echo "> Creating $appPath"
mkdir $appPath
echo "> Touching $dbPath"
touch $dbPath

echo "> Copying static assets to $appPath"
cp -rf ./views $appPath/views
cp -f ./config.yml $appPath/config.yml

echo "> Creating build directories to $buildsPath"
mkdir $buildsPath/builds
mkdir $buildsPath/builds/deleted
mkdir $buildsPath/builds/published

echo "> Setting permissions"
chmod 555 /usr/local/bin/android_ota_server 

chmod -R 0770 $dbPath
chmod -R 0770 $appPath
chmod -R 0770 $buildsPath

chown -R android_ota_server:storage $dbPath
chown -R android_ota_server:storage $appPath
chown -R uploader:storage $buildsPath/builds

echo "> Installing systemd service"
cp android_ota_server.service /etc/systemd/system/android_ota_server.service
systemctl enable android_ota_server.service

echo "> Adding android_ota_server user"
useradd android_ota_server -s /bin/false
groupadd storage 
usermod -g storage android_ota_server

echo "> Starting android_ota_server"
systemctl enable android_ota_server.service

echo "-----------------------------"
echo "Installed successfully!"
echo "Now run: systemctl start android_ota_server"
