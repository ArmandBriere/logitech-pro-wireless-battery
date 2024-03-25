# Logitech Pro Wireless Battery

HeadsetControl daemon listening to headset battery state

## Installation

```bash
git clone https://github.com/ArmandBriere/logitech-pro-wireless-battery
cd logitech-pro-wireless-battery
```

### Make

- Executing the `make` command will build the source code, copy the service to the systemd folder and then enable the new service

```bash
make install
```

### Headsetcontrol

- This package need to be [installed](https://github.com/Sapd/HeadsetControl)

```bash
git clone git@github.com:Sapd/HeadsetControl.git

cd HeadsetControl
mkdir build && cd build
cmake ..
make install
```

#### How to resolve permission denied for non-root users

Error example:

```txt
Failed to open requested device.
 HID Error: Failed to open a device with path '/dev/hidraw5': Permission denied
```

To fix the issue, add a file with the following line to `/etc/udev/rules.d` (in case of manjaro / arch you want to create a new file in /etc/udev/rules.d)

`sudo nano /etc/udev/rules.d/46-logiops-hidraw.rules`

content:
```txt
# Assigns the hidraw devices to group hidraw, and gives that group RW access:
KERNEL=="hidraw[0-9]*", GROUP="hidraw", MODE="0660"
```

- Add the group hidraw if it does not already exist:

`sudo groupadd --system hidraw`

- Add the users who will run logid to that group:

`sudo usermod -G hidraw -a $USER`

- Restart

### Manual

```bash
# Build the source code
go build -o logitech-pro-wireless-notificationd

# Copy the source code to the systemd folder
sudo cp logitech-pro-wireless-notificationd.service /etc/systemd/user/logitech-pro-wireless-notificationd.service
sudo cp logitech-pro-wireless-notificationd /etc/systemd/user/logitech-pro-wireless-notificationd

# Reload the deamon
sudo systemctl daemon-reload

# Enable the service
systemctl --user enable --now logitech-pro-wireless-notificationd
```

## Development

- The program can be run in verbose mode to get debug logs

```bash
go run main.go --verbose
```

