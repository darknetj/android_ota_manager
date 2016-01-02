# Android OTA Server

Golang web app that provides a simple HTTP server to deliver OTA updates.

How to Run:

    ./install.sh
    $ android_ota_manager -env=production -config=/var/lib/ota_server/config.yml

Command line flags:

    -dev=development|production
    -config=<path to config>
    -add_user

Adding an admin user:

    $ android_ota_manager -add_user -env=production -config=/var/lib/ota_server/config.yml
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

I recommend using rsync over SSH to upload files:

    rsync -zv builds/* server:/var/lib/android_ota_manager/builds
