
log-env:
	sudo mkdir -p /var/log/sw-api-server/
	sudo chmod -R 777 /var/log/sw-api-server/

build:
	CGO_ENABLED=0 go build -a -installsuffix cgo -o sw-api main.go

run:
	go run main.go

compose-up:
	docker-compose up -d

compose-down:
	docker-compose down

docker-build:
	docker build -t sw-api .

docker-run:
	docker run --rm --name=sw-api \
		-p 8080:8080 \
		sw-api