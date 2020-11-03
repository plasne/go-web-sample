FROM golang:alpine as build
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .

FROM scratch as run
WORKDIR /app
COPY --from=build /build/main .
EXPOSE 8080
CMD ["./main"]