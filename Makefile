build:
	docker build -t app:0.0.1 .
start:
	docker-compose up -d 
down:
	docker-compose down