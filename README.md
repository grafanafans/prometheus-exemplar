# prometheus-exemplar

A random duration response app to test prometheus exemplar.

## Installation

```bash
git clone git@github.com:grafanafans/prometheus-exemplar.git
cd prometheus-exemplar
make start
```

## How to test

use wrk to send requests

```
sudo apt install wrk

wrk -c 2 -d 3000 http://localhost:8080/v1/books
wrk -c 2 -d 3000 http://localhost:8080/v1/books/1
```


## All in one with Grafana

### Add data sources

#### add Mimir

- HTTP URL：http://mimir:8080/prometheus    
- Custom HTTP Headers: X-Scope-OrgID=demo  
- add exemplar configuration:  

![mimir-exemplar](https://user-images.githubusercontent.com/41465048/182307110-f9275ec3-923f-45c2-b373-5974f17ad42e.PNG)


#### add Tempo  

- HTTP URL: http://tempo:3200  
- Custom HTTP Headers: X-Scope-OrgID=demo

#### add Loki 

- HTTP URL: http://loki:3100 
- Custom HTTP Headers: X-Scope-OrgID=demo 
- Regex:(?:traceID|trace_id|TraceID|TraceId)=(\w+)  
- NOTICE: Regex can be modified according to your own TraceID characters.  

![Image_20220802143410](https://user-images.githubusercontent.com/41465048/182307761-7cc9ae9e-764c-48da-92e5-4692d132f7f8.png)


### Add dashboard

all in one look like:

![lgtm-all](https://user-images.githubusercontent.com/1459834/188675089-5f757184-97a3-4f8e-910e-4daeb0bf55b5.jpg)

#### add metrics panel

In "Explore" module, when you query metrics by Mimir:  

```
histogram_quantile(0.95, rate(http_durations_histogram_seconds_bucket{}[1m]))
```
then open the "Exemplars" flag, it shows a green exemplar point, then click "Query with Tempo" to jump to Tempo Explore.

![metric+tempo1](https://user-images.githubusercontent.com/41465048/182309495-17c446ca-0d0b-4a46-8192-af7eae21c5b0.PNG)

#### add `Traces` pannel 

```
{app="exemplar-demo"} |= `traceID`  
```

click "Tempo" to jump to Tempo Explore.
![log+trace](https://user-images.githubusercontent.com/41465048/182306425-a3eadfa4-60cc-41ab-ac0a-2fda7168504f.PNG)
