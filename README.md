# Android OTA Manager

A toolkit for uploading and distributing Android OTA updates.

* `ota_server` = Returns JSON array of build images. OTA clients on Android devices (or your website) can connect to this service to retrieve a list of available OTA updates

Supported clients:

* Cyanogenmod CMUpdater https://github.com/CyanogenMod/android_packages_apps_CMUpdater
* Coming soon: OpenDelta  https://github.com/omnirom/android_packages_apps_OpenDelta

How to Run:

    $ go get github.com/copperhead-security/android_ota_manager/ota_server
    $ ota_server

(Optional) Command line flags:

    ./ota_server [options]
        -env=development|production
        -config=<path to config>
        -add_user

Production Deployments:

    $ git clone github.com/copperhead-security/android_ota_manager
    $ cd android_ota_manager
    $ ./install
    $ sudo systemctl enable ota_server.service
    $ sudo systemctl start ota_server.service

Adding an admin user:

    $ ota_server -add_user -env=production -config=/path/to/config.yml
    $ > Enter a username:
    $ > Enter a password:
    $ User saved!

Routes:

    /          = JSON object with list of OTA releases retrieved by Android app
    /login     = Admin login form
    /files     = List of build image files uploaded to /builds folder
    /releases  = Publish builds by creating a Release containing a filename, version number, release notes, etc
    /users     = List of admin users (read only)
    /logout

Uploading a file:

The simplest approach is using rsync authenticated with SSH to upload files:

    rsync -zv build_image.zip server:~/ota_builds
