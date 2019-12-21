FROM golang:1.13.4

RUN mkdir /server
WORKDIR /server
COPY . .

RUN GOPROXY=https://goproxy.io go build -v

CMD ["/server/server"]

EXPOSE 20080