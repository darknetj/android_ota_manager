#!/bin/sh
#
# Un-installs android_ota_manager & preps directories
#

echo "Un-installing Android OTA Server"
echo "-------------------------------------"

if [ `whoami` != "root" ]; then
	echo "This script needs to be run as root!"
	exit 1
fi

appPath="/var/lib/android_ota_manager"
dbPath="/var/lib/android_ota_manager/ota.sql"
buildsPath="/home/storage"

echo "> Removing binary from /usr/local/bin"
rm -fv /usr/local/bin/android_ota_manager

echo "> Creating $appPath"
rm -rfv $appPath
echo "> Removing $dbPath"
rm -fv $dbPath

echo "> Removing build directories to $buildsPath"
rm -rfv $buildsPath

echo "> Removing systemd service"
systemctl disable android_ota_manager.service
rm -fv /etc/systemd/system/android_ota_manager.service

echo "> Removing android_ota_manager user"
userdel --force android_ota_manager
delgroup storage
delgroup android_ota_manager

echo "-----------------------------"
echo "Un-installed successfully!"
