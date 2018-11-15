build:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/vk-callback-api
clean:
	rm ./bin/vk-callback-api

compose: build 
	docker-compose build
	docker-compose up