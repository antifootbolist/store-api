# Описание проекта    
В проекте используется 3 контейнера:
- для запуска приложения, написанного на Go   
- для запуска СУБД PostgreSQL   
- для запуска Nginx, который проксирует запросы на http://hosname/api/v1 на контейнер с Go   

Также для выполнения дополнительных операций используются вспомогательные контейнеры:
- для генерации документации ApiDoc  
- запуска миграции схемы СУБД с помощью Flyway   

Приложение на Go фреймворк GORM (https://gorm.io/index.html), который позволяет сокращать код для работы с СУБД. Для миграции схемы используется flyway (https://flywaydb.org). Тестовые данные описаны в виде кода и загружаются в зависимости от значения параметра TEST_DATA при старте приложения.    
# Структура проекта
1. `flyway` - каталог миграций схемы СУБД:  
	- `conf/flyway.conf` - конфиг Flyway  
	- `sql/*.sql` - файлы миграций  
2. `go-app` - каталог приложения на Go:  
	- `apidoc.json` - конфиг apidoc
	- `Dockerfile` - для сборки контейнера с приложением
	- `Dockerfile.apidoc` - для генерации apidoc  
	- `env.list` - параметры подключения к БД, которые передаются при запуске контейнера (docker run)  
	- `go.mod` - для выгрузки требуемых пакетов Go  
	- `go.sum` - для выгрузки требуемых пакетов Go  
	- `main.go` - исходный код приложения с инструкциями apidoc   
3. `nginx` - каталог веб сервера Nginx:
	- `Dockerfile` - для сборки контейнера с Nginx
	- `nginx.conf` - конфиг Nginx
4. `postgresql` - каталог БД на PostgreSQL:
	- `Dockerfile` - для сборки контейнера с БД
	- `init.sql` - первоначальные инструкции для настройки СУБД для работы GORM
5. `Jenkinsfile*` - файлы для автоматизации деплоя описанного выше функционала. Инструкции по настройке Jenkins находятся внутри самого файла.   

# Запуск проекта
## Первоначальные настройки
1. Сконировать репозиторий на сервер.
2. Перейти в склонированный каталог.
3. Установить Docker Engine (например, для CentOS делается вот так - https://docs.docker.com/engine/install/centos/).
4. Создать общую сеть для контейнеров:
```
 docker network create app-net 
```
## Запуск контейнера с БД
1. [Опционально] изменить пароль пользователя `postgres` в `postgresql/Dockerfile`    
2. [Опционально] изменить пароль пользователя `user-api` в `postgresql/init.sql`. После этого требуется изменить данный пароль также в параметрах подключения приложения к БД `go-app/env.list`.  
3. Выполнить команды по сборке и запуску образа:
```
docker build -t postgresql -f postgresql/Dockerfile .
docker run -d --restart always --name postgresql --network app-net -p 5432:5432 postgresql
```  
После выполнения данного шага будет задан пароль пользователю postgres. Создан пользователь `user-api` и DB `store_api`, владельцем которой является пользователь `user-api`.

## Запуск миграции с помощью Flyway
1. Установить Flyway Community Edition в соответствии с инструкцией:  
https://documentation.red-gate.com/fd/command-line-184127404.html   
2. Запустить миграцию схемы СУБД:  
```flyway -configFiles=./flyway/conf -locations=filesystem:./flyway/sql migrate```
3. Проверить что миграции были успешно применены:  
```flyway -configFiles=./flyway/conf/flyway.conf -locations=filesystem:./flyway/sql info```

## Запуск контейнера с приложением
1. [Опционально] изменить пароль от пользователя *user-api* (если был изменен при запуске БД см. п.2)
2. Выполнить команды по сборке и запуску образа:   
```
docker build -t go-app -f go-app/Dockerfile .
docker run -d --restart always --name go-app --network app-net -p 8080:8080 --env-file go-app/env.list -e TEST_DATA=False go-app
```
После выполнение данного шага будет скомплирован main.go и запущен как go-app приложение внутри контейнера. При старте приложения он подключится с БД (параметры указаны в файле `go-app/env.list`) и загрузит тестовые данные, если параметр `TEST_DATA` задан в значение `True`.

## Запуск контейнера с Nginx  
1. Выполнить команды по сборке и запуску образа:   
```
docker build -t antifootbolist/nginx -f nginx/Dockerfile .
docker run -d --restart always --name nginx --network app-net -p 80:80 antifootbolist/nginx
```
## Проверка работоспособности приложения  
Сервис должен быть доступен по следующим endpoints:  
`http://localhost/api/v1/product` - для вывода все продуктов в БД (метод Get)  
`http://localhost/api/v1/product/list/:id` - для вывода информации по конкретному id продукта (метод Get)  
`http://localhost/api/v1/product/update/:id` - для изменения информации о продукте (метод POST)


# Автоматическое обновление документации
Detail information about API and requests is published to GitHub Pages   

# How to deploy this project with Jenkins
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
