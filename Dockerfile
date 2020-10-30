FROM golang:1.14.10-buster

LABEL maintainer=developers@gepardec.com

ENV GOBIN=$GOPATH/bin

RUN apt install git -y

RUN mkdir /app
ADD main.go /app

WORKDIR /app

RUN go get
RUN go build -o main .

EXPOSE 8082

CMD /app/main