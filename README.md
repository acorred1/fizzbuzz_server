# README.md

## Exercise Requirements
Exercise: Write a simple fizz-buzz REST server. 

"The original fizz-buzz consists in writing all numbers from 1 to 100, and just
replacing all multiples of 3 by ""fizz"", all multiples of 5 by ""buzz"", and
all multiples of 15 by ""fizzbuzz"". 

The output would look like this:
""1,2,fizz,4,buzz,fizz,7,8,fizz,buzz,11,fizz,13,14,fizzbuzz,16,...""."

Your goal is to implement a web server that will expose a REST API endpoint that:

Accepts five parameters: three integers int1, int2 and limit, and two strings
str1 and str2.  Returns a list of strings with numbers from 1 to limit, where:
all multiples of int1 are replaced by str1, all multiples of int2 are replaced
by str2, all multiples of int1 and int2 are replaced by str1str2.
 

The server needs to be:

* Ready for production
* Easy to maintain by other developers
 

Bonus: add a statistics endpoint allowing users to know what the most frequent
request has been. This endpoint should:

* Accept no parameter
* Return the parameters corresponding to the most used request, as well as the
  number of hits for this request

## Project overview & API Endpoints

The project is a simple HTTP server that uses the Gin web framework. The main server
process is configured to serve out of port 8080. 
For statistics, the server runs alongside a Prometheus instance that's available at port 9090.

The server exposes 3 endpoints:

* `/api/v1/fizzbuzz`: computes the configurable FizzBuzz sequence
* `/api/v1/stats`: returns the top request's number of hits and params
* `/metrics`: instrumentation endpoint used by Prometheus to scrape metrics for the app 

## File Structure
The project is split into the main following files:

* `main.go`: server setup with routes and logging configuration
* `models.go`: contains structs, vars, and consts used throughout the package
* `handlers.go`: handler functions for the endpoints exposed
* `helpers.go`: helper functions. Notably includes `fizzbuzz`, which contains the main logic for computing the fizzbuzz sequences
* `Dockerfile`: for building the fizzbuzz server image
* `docker-compose.yml`: configuration for fizzbuzz server + prometheus
* `prometheus.yml`: prometheus configuration

## Running the server

### Prerequisites
In order to run the server locally you'll need:
- Go
- docker
- git

To run the server locally:

1) Clone the repo
```bash
git clone git@github.com:acorred1/fizzbuzz_server.git
cd fizzbuzz_server
```
2) Build the images
```bash
docker compose build
```

3) Run the images
```bash
docker compose up
```

## Sample requests

1) For the fizzbuzz endpoint
```bash
curl 'http://localhost:8080/api/v1/fizzbuzz?limit=100&str1=foo&int2=3&str2=bar&int1=7'

{"result":["1","2","bar","4","5","bar","foo","8","bar","10","11","bar","13","foo","bar","16","17","bar","19","20","foobar","22","23","bar","25","26","bar","foo","29","bar","31","32","bar","34","foo","bar","37","38","bar","40","41","foobar","43","44","bar","46","47","bar","foo","50"]}
```

2) For the stats endpoint
```bash
curl http://localhost:8080/api/v1/stats

{"top_fizzbuzz_request":{"request_params":{"int1":7,"int2":3,"limit":100,"str1":"foo","str2":"bar"},"num_hits":1}}
```
