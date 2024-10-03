
build:
	docker build -t blumbot .

up:
	docker-compose up -d

down:
	docker-compose down

delete:
	docker rmi blumbot --force

.PHONY: build up down delete