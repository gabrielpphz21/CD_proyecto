FROM golang:1.27rc3-alpine3.24 AS builder

WORKDIR /golangApp

COPY go.work go.work.sum ./
COPY satelites_consumer/ ./satelites_consumer/
COPY kafkax/ ./kafkax/
COPY events/ ./events/

RUN go work edit -dropuse=./back_satelites -dropuse=./back_info -dropuse=./back_cohetes -dropuse=./protoc -dropuse=./process_path -dropuse=./cohetes_consumer


WORKDIR /golangApp/satelites_consumer
RUN CGO_ENABLED=0 GOOS=linux go build -o /golangApp/main
#lightweight 
FROM alpine:3.24
WORKDIR /golangApp
COPY --from=builder /golangApp/main .

EXPOSE 9000
CMD [ "./main" ]