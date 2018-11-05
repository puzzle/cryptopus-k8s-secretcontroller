FROM golang:1.11 AS builder

# Download and install the latest release of dep
ADD https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep


WORKDIR $GOPATH/src/github.com/puzzle/cryptopus-k8s-secretcontroller
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /app .

FROM scratch
COPY --from=builder /etc/passwd /etc/passwd
USER nobody
COPY --from=builder /app ./
ENTRYPOINT ["./app", "-logtostderr"]
