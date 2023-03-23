# prometheus-exemplar

A random duration response app to test Grafana LGTM stack.

## Installation

- git clone code

```bash
git clone git@github.com:grafanafans/prometheus-exemplar.git
cd prometheus-exemplar
```

- start app

```
docker-compose up -d
```

stop app

```
docker-compose down
```

## How to test

use wrk to send requests

```
wrk http://localhost:8080/v1/books
wrk http://localhost:8080/v1/books/1
```

Then visit `http://localhost:3000` page to see demo app dashboard.

## All in one with Grafana

### Add data sources

#### add Mimir 
   
Add exemplar configuration:  

![mimir-exemplar](https://user-images.githubusercontent.com/41465048/182307110-f9275ec3-923f-45c2-b373-5974f17ad42e.PNG)

Notes:
    
- HTTP URL：http://lb.com/metrics   

#### add Tempo  

- HTTP URL:http://lb.com/traces  

#### add Loki with Derived fields  

- HTTP URL:http://lb.com/logs
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

