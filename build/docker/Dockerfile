FROM golang:1.10.1-alpine3.7 as builder
WORKDIR /go/src/kubesphere.io/devops/
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v  -a -installsuffix cgo -ldflags '-w'  -o cmd/ks-devops pkg/server/main.go



FROM alpine
COPY --from=builder /go/src/kubesphere.io/devops/cmd/* /usr/local/bin/
EXPOSE 8080
CMD ["/usr/local/bin/ks-devops"]
