#1st stage: build go binary
FROM golang:1.14 as builder

WORKDIR /project

# Copies folder
COPY /api ./api
COPY /cli ./cli
COPY /internal ./internal

# Copies and downloads dependencies
COPY /go.mod ./go.mod
RUN go mod vendor

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