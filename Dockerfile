FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/go-etl

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go get github.com/Nhanderu/brdoc
RUN go get github.com/gorilla/mux
RUN go get github.com/lib/pq
RUN go install go-etl

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/go-etl

# Document that the service listens on port 8080.
EXPOSE 8080

