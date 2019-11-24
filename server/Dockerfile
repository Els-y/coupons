FROM golang:1.13.4

RUN mkdir /server
WORKDIR /server
COPY . .

RUN go build -v

CMD ["/server/server"]
