FROM golang:1.14-buster as develop

ENV GolangCILintVersion=1.24.0

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
        sh -s -- -b $(go env GOPATH)/bin v${GolangCILintVersion} \
        && go get github.com/smartystreets/goconvey

FROM develop as build

ENV USER=appuser
ENV UID=10001

WORKDIR /go/api
COPY ./ /go/api

RUN DEBIAN_FRONTEND=noninteractive apt-get update && apt-get install -y \
        ca-certificates tzdata \
        && update-ca-certificates \
        && adduser \
           --disabled-password \
           --gecos "" \
           --home "/nonexistent" \
           --shell "/sbin/nologin" \
           --no-create-home \
           --uid "${UID}" \
           "${USER}" \
        && make api-build

FROM scratch

COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group
COPY --from=build /go/api/dist/cmd/todo/api /bin/api

USER appuser:appuser

ENTRYPOINT ["/bin/api"]