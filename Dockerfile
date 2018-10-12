FROM golang:1.11.0 AS builder
WORKDIR /go/src/github.com/atmiguel/cerealnotes
COPY . .
RUN go get github.com/golang/dep/cmd/dep
RUN dep ensure

# Cross compile cerealnotes to work in a minimal alpine image. CGO must be
# disabled for cross compilation. See https://github.com/golang/go/issues/5104
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cerealnotes


# Use alpine docker image for production for the small image size (5MB)
FROM alpine:3.8
WORKDIR /root/
COPY --from=builder /go/src/github.com/atmiguel/cerealnotes/cerealnotes .
COPY --from=builder /go/src/github.com/atmiguel/cerealnotes/templates ./templates
COPY --from=builder /go/src/github.com/atmiguel/cerealnotes/static ./static
CMD ["./cerealnotes"]
EXPOSE 8080
