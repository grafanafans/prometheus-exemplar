# prometheus-exemplar

A random duration response app to test prometheus exemplar.

## Installation

1. clone code and build app

```bash
git clone git@github.com:songjiayang/prometheus-exemplar.git
cd prometheus-exemplar
go mod vendor
docker-compose build
```

2. start app

```
docker-compose up
```

3. stop app

```
docker-compose down
```

## How to test

1. use wrk to send requests
```
wrk -c 2 -d 3000 http://localhost:8080/v1/books
wrk -c 2 -d 3000 http://localhost:8080/v1/books/1
```

2. Query exemplar with prometheus console 

- visit `http://localhost:9090/graph`
- choose `Graph` tab
- type `histogram_quantile(0.95, rate(http_durations_histogram_seconds_bucket{}[1m]))` to the query input text.

![image](https://user-images.githubusercontent.com/1459834/158003593-7ad63d8d-0d8b-4bea-af54-947410d797a6.png)

## All in one with Grafana
