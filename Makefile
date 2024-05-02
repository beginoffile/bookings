WEB_BINARY_LINUX=bookings
WEB_BINARY_WIN=bookings.exe


## up: starts all containers in the background without forcing build
up:
	@echo Starting Docker images...
	docker-compose up -d
	@echo Docker images started!


## build_front: builds the web end binary
build_windows:
	@echo Building front end binary windows...
	set CGO_ENABLED=0&& set GOOS=windows&& go build -o ${WEB_BINARY_WIN} ./cmd/web
	@echo Done!

## build_backend: builds the web end binary
build_linux:
	@echo Building front end binary linux...
	set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o ${WEB_BINARY_LINUX} ./cmd/web
	@echo Done!

## start: starts the front end
start: build_windows
	@echo Starting web
	${WEB_BINARY_WIN} -dbname=bookings -dbuser=postgres -dbpass="12345678" -cache=false -production=false

## stop: stop the front end
stop:
	@echo Stopping Web...
	@taskkill /IM "${WEB_BINARY_WIN}" /F
	@echo "Stopped Web!"
