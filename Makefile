default: install

BUILD_TARGET = logitech-pro-wireless-notificationd
SYSTEMD_FOLDER = /etc/systemd/user

install: build disable copy enable

test:
	@echo "Start testing the source code..."
	go test -coverprofile=c.out
	@echo "Done testing."

coverage: test
	@echo "Start generating coverage report..."
	go tool cover -html=c.out
	@echo "Done generating coverage report."

build:
	@echo "Start building the source code..."
	go build -o ${BUILD_TARGET}
	@echo "Done building."

copy:
	@echo "Start copying the service..."
	sudo cp ${BUILD_TARGET}.service ${SYSTEMD_FOLDER}/${BUILD_TARGET}.service
	sudo cp ${BUILD_TARGET} ${SYSTEMD_FOLDER}/${BUILD_TARGET}
	@echo "Done copying."
	
	@echo "Start copying icons..."
	sudo cp -r icons ${SYSTEMD_FOLDER}
	@echo "Done copying."

disable:
	@echo "Start disabling the service..."
	sudo systemctl daemon-reload
	-systemctl --user disable --now ${BUILD_TARGET}
	@echo "Done disabling the service."

enable:
	@echo "Start enabling the service..."
	sudo systemctl daemon-reload
	systemctl --user enable --now ${BUILD_TARGET}
	@echo "Done enabling the service."

status:
	systemctl --user status logitech-pro-wireless-notificationd
