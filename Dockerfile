FROM golang:alpine as builder
RUN mkdir /build 
COPY ./go.mod /build/
COPY ./go.sum /build/
COPY ./src /build/src
COPY ./main.go /build/
COPY ./config.yml /build/
WORKDIR /build 
RUN go install
RUN go build -o main .

FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/main /app/
COPY --from=builder /build/config.yml /app/
WORKDIR /app
CMD ["./main"]
