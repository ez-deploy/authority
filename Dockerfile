FROM golang:1.17 as builder

WORKDIR /go/src/github.com/ez-deploy/authority
COPY . .

RUN go env -w GO111MODULE=on && \
    go env -w GOPROXY=https://goproxy.io && \
    CGO_ENABLED=0 go build -tags netgo -o authority ./main.go

FROM busybox

WORKDIR /

COPY --from=builder /go/src/github.com/ez-deploy/authority/authority /authority

EXPOSE 80
ENTRYPOINT [ "/authority" ]