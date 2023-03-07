# Store API
This is example of RESTful API for shop written on Go

## Files   
The following dirs are in this repo:   
```go-app``` - Store API written on Go  
```nginx``` - Nginx config to run Web server   
```postgresql``` - PostgreSQL as a backend   
```*/Dockerfile``` - Dockerfile to run apps inside container   
```Jenkinsfile``` - Jenkins pipeline to deploy Nginx and apps   

## Nginx configuration
When Go API container are up and running, it is available at the following URL:  
```http://server:8080```  

All requests to Nginx (frontend) are proxied according to the following algorithm:   
```http://server/api``` - to Go application   
```http://server/apidoc``` - to apiDoc documentation   


## Example how to launch app written on Go   

### From VM
1. Pull source code from repository:   
```git pull https://github.com/antifootbolist/store-api.git```
2. Run the application:   
```go run main.go```
4. Check a status of the application:  
```curl -X GET http://localhost:8080```

### Inside Docker container
1. Pull source code from repository:   
```git pull https://github.com/antifootbolist/store-api.git```
2. Execute docker image build by using Dockerfile from repo:   
```docker build -t go-app .```
3. Run docker container:  
```docker run -d --name go-app -p 8080:8080```
4. Check a status of the application:  
```curl -X GET http://localhost:8080```
