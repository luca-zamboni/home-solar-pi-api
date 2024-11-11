FROM alpine:3.20.3

WORKDIR /app

COPY ./dist/home-solar-pi-server ./
COPY ./dist/.env ./

EXPOSE 5000

ADD ./solar-home ./solar-home

ENTRYPOINT ["./home-solar-pi-server"]