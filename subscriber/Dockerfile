FROM golang:1.13.4

RUN mkdir /subscriber
WORKDIR /subscriber
COPY . .

RUN GOPROXY=https://goproxy.io go build -v
# CMD ["/subscriber/subscriber", "$SUBJ"]
CMD ["/subscriber/subscriber", "temp", "my-queue"]