build:
	sudo docker build -t app:0.0.1 .
start:
	sudo docker-compose up -d 
stop:
	sudo docker-compose down