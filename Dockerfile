FROM golang:1.21.6-alpine as build
WORKDIR /opt/app
COPY . .
RUN go build -o whosbest cmd/whosbest/main.go

FROM alpine:latest
WORKDIR /opt/app
COPY --from=build /opt/app/whosbest .
EXPOSE 3000
ENTRYPOINT ["./whosbest"]
