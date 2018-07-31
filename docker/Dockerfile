FROM centos:7 

WORKDIR /app
COPY ./bin/wss-exporter /app/wss-exporter
COPY ./config.yml.sample /app/config.yml

CMD ["/app/wss-exporter"]
