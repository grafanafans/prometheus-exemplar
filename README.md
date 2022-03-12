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

### Add data sources

- add Prometheus

![image](https://user-images.githubusercontent.com/1459834/158003719-87caf71b-15cd-4faf-91ed-b2d535053d49.png)

- add Jaeger
![image](https://user-images.githubusercontent.com/1459834/158003756-1dbee018-ca29-4f9c-b1a6-c3d70833be58.png)

- add Loki with Derived fields

![image](https://user-images.githubusercontent.com/1459834/158003785-5f7acec8-f6a3-4086-96a8-5af2935343cd.png)

![image](https://user-images.githubusercontent.com/1459834/158003792-2a573888-b5bf-48a3-9b24-57e37fabaee5.png)


### Create dashboard

- create http duration panel

![image](https://user-images.githubusercontent.com/1459834/158003892-ef807fc8-6870-4a77-b296-05557e2989b9.png)

- create log panel with `traceID` variable

![image](https://user-images.githubusercontent.com/1459834/158004080-b76ce802-342d-4235-bc29-3804ab036255.png)


![image](https://user-images.githubusercontent.com/1459834/158004049-0c6ad080-700e-48e6-85cf-53ed8849b7ea.png)


### Query with traceID

- query with traceID 

![image](https://user-images.githubusercontent.com/1459834/158004110-bbdec2c2-376e-4a06-8cea-27766c089df6.png)

- jump to jaeger trace view
![image](https://user-images.githubusercontent.com/1459834/158004123-cbb4cfb8-2b45-488b-8871-34645b96b413.png)

- spit view with loki and jaeger

![image](https://user-images.githubusercontent.com/1459834/158004183-d7c3857d-e26a-4c27-b86e-791ea85dd7ab.png)
