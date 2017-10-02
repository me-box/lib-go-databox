FROM golang:1.8.3-alpine3.6 as gobuild
WORKDIR /
ENV GOPATH="/"
RUN apk update && apk add git

COPY . .

#get dependencies after copy to stop them from being cached
RUN go get -u golang.org/x/net/websocket
RUN go get -u github.com/me-box/lib-go-databox
RUN addgroup -S databox && adduser -S -g databox databox
RUN GGO_ENABLED=0 GOOS=linux go build -a -tags netgo -installsuffix netgo -ldflags '-d -s -w -extldflags "-static"' -o app /src/main.go

FROM scratch
USER databox
WORKDIR /root/
COPY --from=gobuild /app .
LABEL databox.type="app"
EXPOSE 8080

CMD ["./app"]