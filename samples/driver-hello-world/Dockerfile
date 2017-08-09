FROM golang:1.8.3-alpine3.6 as gobuild
WORKDIR /
ENV GOPATH="/"
RUN apk update && apk add git
RUN go get -u github.com/gorilla/mux 
RUN go get -u golang.org/x/net/websocket
RUN go get -u github.com/me-box/lib-node-databox
COPY . .
RUN GGO_ENABLED=0 GOOS=linux go build -a -tags netgo -installsuffix netgo -ldflags '-d -s -w -extldflags "-static"' -o app /src/main.go

FROM alpine:latest  
WORKDIR /root/
COPY --from=gobuild /app .
COPY --from=gobuild /www/ /root/www/
LABEL databox.type="driver"
EXPOSE 8080

CMD ["./app"]  