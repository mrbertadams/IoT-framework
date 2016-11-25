# DeviceHive D-Bus Framework

[DeviceHive]: http://devicehive.com "DeviceHive framework"
[DataArt]: http://dataart.com "DataArt"

Develop apps connected to devices, peripherials and cloud using commonly available linux dbus interface.

![](framework.png?raw=true)

Documentation is here: http://docs.devicehive.com/docs/iot-toolkit-overview

## Contents

Currently this framework contains few components:

`devicehive-alljoyn` - our custom API to communicate with AllJoyn devices with the easiest way

`devicehive-ble` - our custom API for access BLE(Bluetooth low energy) GATT services

`devicehive-cloud` - our API for access to DeviceHive cloud services

`devicehive-enocean` - our simple API for access to EnOcean devices

`devicehive-gpio` - our API for access to device's pins

`examples` - samples of usage out IoT framework

`build` - build system for our IoT framework

Each directory contains more detailed information about component, see details there in readme files.


## Cloning

This repository uses submodules so to clonse fresh copy use command
```
git clone --recursive git@github.com:devicehive/IoT-framework.git
```

Or if doing pull to existing repo use
```
git pull
git submodule update --init --recursive
```


## Building 

To build `Ubuntu Snappy` package navigate to `build/snappy` and run 
```
# for framework
./build-framework.sh

#for sample apps
./build-apps.sh
```

DeviceHive license
------------------

[DeviceHive] is developed by [DataArt] Apps and distributed under Open Source
[Apache 2.0](https://en.wikipedia.org/wiki/Apache_License). This basically means
you can do whatever you want with the software as long as the copyright notice
is included. This also means you don't have to contribute the end product or
modified sources back to Open Source, but if you feel like sharing, you are
highly encouraged to do so!

&copy; Copyright 2015-2016 DataArt Apps &copy; All Rights Reserved