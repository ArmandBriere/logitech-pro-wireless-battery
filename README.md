# Logitech Pro Wireless Battery

HeadsetControl daemon listenning to headset battery state

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

