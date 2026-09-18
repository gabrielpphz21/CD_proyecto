FROM golang:1.27rc3-alpine3.24 AS builder

WORKDIR /golangApp

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .


#lightweight 
FROM alpine:3.24
WORKDIR /golangApp
COPY --from=builder /golangApp/main .
COPY --from=builder /golangApp/routes.json .


EXPOSE 8080
CMD [ "./main" ]