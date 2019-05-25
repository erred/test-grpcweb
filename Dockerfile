FROM golang AS build

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -o grpc-server ./server

FROM scratch

WORKDIR /app
COPY --from=build /app/grpc-server .
EXPOSE 8080/tcp
ENTRYPOINT ["/app/grpc-server"]
