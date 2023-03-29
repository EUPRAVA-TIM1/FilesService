FROM golang:alpine as build_container
WORKDIR /app
COPY ./FileService/go.mod ./FileService/go.sum ./
RUN go mod download
COPY ./FileService .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o filesrv .

FROM alpine
WORKDIR /root/
COPY --from=build_container /app/filesrv .
EXPOSE 8000
ENTRYPOINT ["./filesrv"]