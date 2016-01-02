#!/bin/sh
#
# Installs android_ota_manager & preps directories
#

echo "Installing Android OTA Server"
echo "-------------------------------------"

if [ `whoami` != "root" ]; then
	echo "This script needs to be run as root!"
	exit 1
fi

appPath="/var/lib/android_ota_manager"
dbPath="/var/lib/android_ota_manager/ota.sql"
buildsPath="/home/storage"

echo "> Copying binary to /usr/local/bin"
cp android_ota_manager /usr/local/bin/

echo "> Creating $appPath"
mkdir -pv $appPath
echo "> Touching $dbPath"
touch $dbPath

echo "> Copying static assets to $appPath"
cp -rf ./views $appPath/views
cp -f ./config.yml $appPath/config.yml

echo "> Creating build directories to $buildsPath"
mkdir -pv $buildsPath/builds
mkdir -pv $buildsPath/builds/deleted
mkdir -pv $buildsPath/builds/published

echo "> Installing systemd service"
cp -v android_ota_manager.service /etc/systemd/system/android_ota_manager.service
systemctl enable android_ota_manager.service

echo "> Adding android_ota_manager user"
useradd android_ota_manager -s /bin/false
groupadd storage 
usermod -g storage android_ota_manager

echo "> Setting permissions"
chmod 555 /usr/local/bin/android_ota_manager

chmod -R 0770 $dbPath
chmod -R 0770 $appPath
chmod -R 0770 $buildsPath

chown -R android_ota_manager:storage $dbPath
chown -R android_ota_manager:storage $appPath
chown -R uploader:storage $buildsPath/builds

echo "> Starting android_ota_manager"
systemctl enable android_ota_manager.service

echo "-----------------------------"
echo "Installed successfully!"
echo "Now run: systemctl start android_ota_manager"
