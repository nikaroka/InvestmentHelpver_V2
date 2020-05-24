FROM golang:1.13
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build -o main .
EXPOSE 8090
CMD ["/app/main"]