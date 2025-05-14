FROM dc1base:latest as builder
WORKDIR /app
COPY . .
RUN go build -o ./bin/dc1server ./app/dc1server/cmd

FROM alpine:latest as prod
WORKDIR /app
COPY --from=builder /app/bin/ /app
COPY --from=builder /app/app/dc1server/configs/config.yaml /app
RUN ls
EXPOSE 8000
CMD ["./dc1server", "-conf", "./config.yaml"]
