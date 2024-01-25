FROM golang:alpine3.17 as build
WORKDIR /opt/app
COPY . .
RUN go build -o whosbest cmd/whosbest/main.go

FROM alpine:latest
WORKDIR /opt/app
COPY --from=build /opt/app/whosbest .
COPY --from=build /opt/app/.env.example .env
ENTRYPOINT ["./whosbest"]
