build:
	docker build -t songjiayang/prometheus-exemplar:0.0.1 .
start:
	docker-compose up -d 
stop:
	docker-compose down