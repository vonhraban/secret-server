FROM golang:latest AS builder
ADD . /app
WORKDIR /app
RUN go mod download
RUN CGO_ENABLED=0 go build -a -o /main .

# final stage
FROM alpine:latest
COPY --from=builder /main ./
ADD ./config.yaml ./config.yaml
RUN chmod +x ./main
ENTRYPOINT ["./main"]
EXPOSE 80