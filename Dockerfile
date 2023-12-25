FROM golang:1.21

RUN mkdir /go/src/translator

WORKDIR /go/src/translator

COPY ./go.mod ./
COPY ./go.sum ./
COPY ./conf.json ./

RUN go mod download
RUN go mod tidy
RUN go get github.com/gilang-as/google-translate
COPY ./src/*.go ./

#build main
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build

CMD ["/go/src/translator/translator"]