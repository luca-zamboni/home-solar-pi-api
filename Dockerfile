FROM alpine:3.20.3

WORKDIR /app

COPY ./dist/home-solar-pi-server ./
COPY ./dist/.env ./

EXPOSE 5000

ENTRYPOINT ["./home-solar-pi-server"]