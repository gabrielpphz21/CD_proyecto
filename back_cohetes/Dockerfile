FROM golang:1.27rc3-alpine3.24 AS builder

WORKDIR /golangApp

COPY go.work go.work.sum ./
COPY protoc/ ./protoc/
COPY back_cohetes/ ./back_cohetes/

RUN go work edit -dropuse=./back_satelites -dropuse=./back_info

WORKDIR /golangApp/back_cohetes
RUN CGO_ENABLED=0 GOOS=linux go build -o /golangApp/main
#lightweight 
FROM alpine:3.24
WORKDIR /golangApp
COPY --from=builder /golangApp/main .

EXPOSE 9000
CMD [ "./main" ]