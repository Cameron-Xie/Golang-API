#!/bin/bash

set -euxo pipefail

tmpf=$(mktemp -d)
cd "$tmpf"

openssl req -new -text -passout pass:abcd -subj /CN=postgres -out server.req -keyout privkey.pem
openssl rsa -in privkey.pem -passin pass:abcd -out server.key
openssl req -x509 -in server.req -text -key server.key -out server.crt
chmod 600 server.key

cp /tmp/postgresql.conf $PGDATA/postgresql.conf
cp "$tmpf"/server.crt $PGDATA/server.crt
cp "$tmpf"/server.key $PGDATA/server.key
chown -R postgres: $PGDATA

