
build:
	CGO_ENABLED=0 go build -a -installsuffix cgo -o sw-api main.go

run:
	go run main.go

docker-build:
	docker build -t star-wars-api .

docker-run:
	docker run --rm --name=sw-api \
		-p 808:8080 \
		sw-api