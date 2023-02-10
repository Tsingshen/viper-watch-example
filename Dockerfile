FROM ubuntu:22.04

WORKDIR /app

COPY viper-watch /app/viper-watch

CMD ["/app/viper-watch"]
