FROM --platform=linux/arm64  golang:1.14.10-buster as build


ENV GOBIN=$GOPATH/bin

RUN apt install git -y

RUN mkdir /app
ADD main.go /app

WORKDIR /app

RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build  -o /app/main 

FROM scratch 
LABEL maintainer=developers@gepardec.com

COPY --from=build /app/main /main

EXPOSE 8082

ENTRYPOINT ["/main"]
