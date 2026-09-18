FROM golang:1.27rc3-alpine3.24 as builder

WORKDIR /golangApp

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .


#lightweight 
FROM alpine3.24
WORKDIR /golangApp
COPY --from=builder /golangAPp/main .

EXPOSE 8080
CMD [ "./main" ]
