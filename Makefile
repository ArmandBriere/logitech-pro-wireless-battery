default: install

BUILD_TARGET = logitech-pro-wireless-notificationd
SYSTEMD_FOLDER = /etc/systemd/user

install: build disable copy enable

build:
	@echo "\nStart building the source code..."
	go build -o ${BUILD_TARGET}
	@echo "Done building.\n"

copy:
	@echo "\nStart copying the service..."
	sudo cp ${BUILD_TARGET}.service ${SYSTEMD_FOLDER}/${BUILD_TARGET}.service
	sudo cp ${BUILD_TARGET} ${SYSTEMD_FOLDER}/${BUILD_TARGET}
	@echo "Done copying.\n"
	
	@echo "\nStart copying icons..."
	sudo cp -r icons ${SYSTEMD_FOLDER}
	@echo "Done copying.\n"

disable:
	@echo "\nStart disabling the service..."
	sudo systemctl daemon-reload
	-systemctl --user disable --now ${BUILD_TARGET}
	@echo "Done disabling the service.\n"

enable:
	@echo "\nStart enabling the service..."
	sudo systemctl daemon-reload
	systemctl --user enable --now ${BUILD_TARGET}
	@echo "Done enabling the service.\n"

status:
	systemctl --user status logitech-pro-wireless-notificationd
