FROM golang


ADD . /go/src/github.com/LuhanM/go-etl
WORKDIR /go/src/github.com/LuhanM/go-etl

RUN go get github.com/Nhanderu/brdoc
RUN go get github.com/gorilla/mux
RUN go get github.com/lib/pq
RUN go build

CMD [ "./go-etl" ]

EXPOSE 8080

