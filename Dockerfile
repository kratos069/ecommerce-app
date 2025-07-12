# Build Stage 
FROM golang:1.24-alpine3.22 AS builder
 WORKDIR /app
#  copy from root (where Dockerfile ran) project to current work dir (/app)
COPY . .
#  output binary file (main), o is output, current main file = main.go
RUN go build -o main main.go

# Run Stage
FROM alpine:3.22
WORKDIR /app
# /app/main is the path to the file we want to copy, "." is current work dir
COPY --from=builder /app/main .
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./db/migration

EXPOSE 8080 9090
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]