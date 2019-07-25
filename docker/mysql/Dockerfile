FROM mysql:8.0.16

ENV MYSQL_DATABASE=van \
    MYSQL_ROOT_PASSWORD=root

ADD sql/schema.sql /docker-entrypoint-initdb.d

EXPOSE 3306