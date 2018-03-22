FROM golang:1.9-alpine as memorycell

WORKDIR /go/src/github.com/connorwalsh/dockerlang
COPY . .
# no need to go get since we are vendoring all our deps
RUN cd ./memorycell && ls && go build

FROM alpine

WORKDIR /usr/bin
COPY --from=memorycell  /go/src/github.com/connorwalsh/dockerlang/memorycell/memorycell .

CMD ["/usr/bin/memorycell"]
