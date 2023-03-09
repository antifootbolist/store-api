# Store API
This is example of RESTful API for shop written on Go

## Files   
The following dirs are in this repo:   
```go-app``` - Store API written on Go  
```nginx``` - Nginx config to run Web server   
```postgresql``` - PostgreSQL as a backend   
```*/Dockerfile``` - Dockerfile to run Go, Nginx and PostgreSQL inside a container   
```Jenkinsfile``` - Jenkins pipeline to deploy Nginx and apps   

## API end points
```http://server/api/v1/product``` - List all products   
```http://server/api/v1/product/list/:id``` - Get information about a specified product  
```http://server/api/v1/product/update/:id``` - Update a product information  

Detail information about API and requests is published to GitHub Pages   

## How to deploy this app with Jenkins
1. Clone and fork the repo
2. Create variable in Jenkins (or create use as a parameter for pipeline):  
`PROD_IP` - IP address of server where we deploy containers
3. Create credentials in Jenkins:
    - GitHub login/token
    - Docker Hub login/password
4. Create pipeline in Jenkins.
5. Change variables in Jenkinsfile:   
`DOCKER_HUB_USER` - Docker Hub username   
`REPO_URL` - URL to clonned/forked repo   
`GHP_URL` - URL to GitHub Pages   
`GIT_AUTHOR_NAME`, `GIT_AUTHOR_EMAIL` - name and email for **git commit** to update apiDoc in GitHub Pages   
`GH_TOKEN_ID` - ID of credential of **GitHub** username/token (from step 3)   
`DOCKER_HUB_LOGIN` - ID of credential of **Docker Hub** username/password (from step 3)   
