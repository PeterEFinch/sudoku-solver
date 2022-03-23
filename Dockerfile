#1st stage: build go binary
FROM golang:1.18 as builder

WORKDIR /project

# Copies and downloads dependencies
COPY /go.mod ./go.mod
RUN go mod vendor

# Copies folder
COPY /api ./api
COPY /cli ./cli
COPY /internal ./internal

# See which patch version of Go is used (important for security updates)
RUN go version

#Builds go binary
RUN CGO_ENABLED=0        \
    GOOS=linux           \
    go install           \
      -a                 \
      -installsuffix cgo \
      --ldflags="-s"     \
      ./cli/server

#2nd stage:
FROM alpine:latest

#Required for accessing
RUN apk --no-cache add ca-certificates

#Copies binary and plugins
WORKDIR /root/
COPY --from=builder /go/bin/ .

#Starting command
CMD ["./server"]