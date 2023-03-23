build:
	docker build -t songjiayang/prometheus-exemplar:0.0.2 .
start:
	docker-compose up -d 
down:
	docker-compose down