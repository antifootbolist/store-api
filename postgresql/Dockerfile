FROM postgres:latest

ENV POSTGRES_USER postgres
ENV POSTGRES_PASSWORD qwe123
ENV POSTGRES_DB dbname

COPY postgresql/init.sql /docker-entrypoint-initdb.d/

EXPOSE 5432