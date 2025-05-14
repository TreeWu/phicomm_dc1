FROM dc1base:latest as builder
WORKDIR /app
COPY . .
RUN go build -o ./bin/gateway ./app/gateway/cmd

FROM alpine:latest as prod
WORKDIR /app
COPY --from=builder /app/bin/ /app
COPY --from=builder /app/app/gateway/configs/config.yaml /app
EXPOSE 8002
CMD ["./gateway" , "-conf" , "./config.yaml"]
