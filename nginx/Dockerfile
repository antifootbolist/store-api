FROM node:alpine AS apidoc

WORKDIR /app

RUN apk add --no-cache git
RUN npm install -g apidoc

COPY go-app/main.go .
RUN apidoc -i ./ -o ./apidoc

FROM nginx:alpine

WORKDIR /usr/share/nginx/html

COPY --from=apidoc /app/apidoc apidoc
COPY nginx/nginx.conf /etc/nginx/nginx.conf

EXPOSE 80
