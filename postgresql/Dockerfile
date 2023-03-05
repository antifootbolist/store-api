FROM postgres:latest

ENV POSTGRES_USER postgres
ENV POSTGRES_PASSWORD qwe123

COPY postgresql/init.sql /docker-entrypoint-initdb.d/

EXPOSE 5432