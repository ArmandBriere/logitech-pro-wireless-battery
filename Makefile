default: install

BUILD_TARGET = logitech-pro-wireless-notificationd
SYSTEMD_FOLDER = /etc/systemd/user

install: build disable copy enable

build:
	@echo -e "\nStart building the source code..."
	go build -o ${BUILD_TARGET}
	@echo -e "Done building.\n"

copy:
	@echo -e "\nStart copying the service..."
	sudo cp ${BUILD_TARGET}.service ${SYSTEMD_FOLDER}/${BUILD_TARGET}.service
	sudo cp ${BUILD_TARGET} ${SYSTEMD_FOLDER}/${BUILD_TARGET}
	@echo -e "Done copying.\n"
	
	@echo -e "\nStart copying icons..."
	sudo cp -r icons ${SYSTEMD_FOLDER}
	@echo -e "Done copying.\n"

disable:
	@echo -e "\nStart disabling the service..."
	sudo systemctl daemon-reload
	systemctl --user disable --now ${BUILD_TARGET}
	@echo -e "Done disabling the service.\n"

enable:
	@echo -e "\nStart enabling the service..."
	sudo systemctl daemon-reload
	systemctl --user enable --now ${BUILD_TARGET}
	@echo -e "Done enabling the service.\n"
