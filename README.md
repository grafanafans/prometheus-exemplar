# prometheus-exemplar

A random duration response app to test prometheus exemplar.

## Installation

1. clone code and build app

```bash
git clone git@github.com:grafanafans/prometheus-exemplar.git
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
2. If you failed to install wrk tool, you can modify Metrics function in middleware.go by adding sleep time to simulate an timeout.  

```go
func Metrics(metricPath string, urlMapping func(string) string) gin.HandlerFunc {
  ....
  method := c.Request.Method
  url := urlMapping(c.Request.URL.Path)
  // add sleep time
  time.Sleep(time.Second)

  elapsed := float64(time.Since(start)) / float64(time.Second)
  ...
}
```

then curl:

```
curl -v http://0.0.0.0:8080/v1/books
```

## All in one with Grafana

### Add data sources

#### add Mimir  
- HTTP URL：http://load-balancer:9009/prometheus      
- add exemplar configuration:  

![mimir-exemplar](https://user-images.githubusercontent.com/41465048/182307110-f9275ec3-923f-45c2-b373-5974f17ad42e.PNG)


#### add Tempo  

- HTTP URL:http://tempo:3200  

#### add Loki with Derived fields  

- HTTP URL:http://loki:3100  
- Regex:(?:traceID|trace_id|TraceID|TraceId)=(\w+)  
- NOTICE: Regex can be modified according to your own TraceID characters.  

![Image_20220802143410](https://user-images.githubusercontent.com/41465048/182307761-7cc9ae9e-764c-48da-92e5-4692d132f7f8.png)


### Query exemplar metrics

#### query mimir with tempo

In "Explore" module, when you query metrics by Mimir:  

```
histogram_quantile(0.95, rate(http_durations_histogram_seconds_bucket{}[1m]))
```
then open the "Exemplars" flag, it shows a green exemplar point, then click "Query with Tempo" to jump to Tempo Explore.

![metric+tempo1](https://user-images.githubusercontent.com/41465048/182309495-17c446ca-0d0b-4a46-8192-af7eae21c5b0.PNG)

#### query loki with tempo  

```
{app="exemplar-demo"} |= `traceID`  
```

click "Tempo" to jump to Tempo Explore.
![log+trace](https://user-images.githubusercontent.com/41465048/182306425-a3eadfa4-60cc-41ab-ac0a-2fda7168504f.PNG)

