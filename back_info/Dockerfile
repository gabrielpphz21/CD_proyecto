FROM golang:1.27rc3-alpine3.24 AS builder

WORKDIR /golangApp

COPY go.work go.work.sum ./
COPY protoc/ ./protoc/
COPY back_info/ ./back_info/

RUN go work edit -dropuse=./back_satelites -dropuse=./back_cohetes

WORKDIR /golangApp/back_info
RUN CGO_ENABLED=0 GOOS=linux go build -o /golangApp/main
#lightweight 
FROM alpine:3.24
WORKDIR /golangApp
COPY --from=builder /golangApp/main .


EXPOSE 9000
CMD [ "./main" ]