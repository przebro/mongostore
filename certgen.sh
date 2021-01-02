#!/bin/bash
#generate keys for rootca server and client
openssl genrsa -out ./docker/cert/root_ca.key 4096
openssl genrsa -out ./docker/cert/bazaarsrv.key 2048
openssl genrsa -out ./docker/cert/client.key 2048

#create root certficiate
openssl req -x509 -new -nodes -key ./docker/cert/root_ca.key -sha256 -days 1024 \
-out ./docker/cert/root_ca.crt \
-subj "/C=US/ST=State/L=City/O=Couchdbdrv/OU=DEV/CN=localuser/emailAddress=localuser@localhost.com"
#create signing request
openssl req -new -key ./docker/cert/bazaarsrv.key -out ./docker/cert/bazaarsrv.csr \
-subj "/C=US/ST=State/L=City/O=Couchdbsrv/OU=DEV/CN=localuser/emailAddress=localuser@localhost.com"
#create certificate for server
openssl x509 -req -in ./docker/cert/bazaarsrv.csr -CA ./docker/cert/root_ca.crt -CAkey ./docker/cert/root_ca.key \
-CAcreateserial -out ./docker/cert/bazaarsrv.crt -days 1024 -sha256 \

#create certificate for client
openssl req -new -key ./docker/cert/client.key -out ./docker/cert/client.req \
-subj "/C=US/ST=State/L=City/O=couchdbcli/OU=DEV/CN=localuser/emailAddress=localuser@localhost.com"
openssl x509 -req -in ./docker/cert/client.req -CA ./docker/cert/root_ca.crt -CAkey ./docker/cert/root_ca.key \
-set_serial 101010 -extensions client -days 1024 -out ./docker/cert/client.crt

cat ./docker/cert/bazaarsrv.key ./docker/cert/bazaarsrv.crt > ./docker/cert/bazaarsrv.pem

#Subject: C = PL, ST = State, L = Lodz, O = Couchdrv, OU = section, CN = localuser, emailAddress = localuser@localhost.com
#Subject: C = US, ST = State, L = City, O = CouchDB driver ltd., OU = DEV, CN = localuser, emailAddress = localuser@localhost.com

#Subject: emailAddress = user@host.com, C = US, ST = State, L = City, O = CouchDB driver ltd., OU = DEV, CN = localhost
#Subject: C = PL, ST = State, L = City, O = Company ltd, OU = section, CN = localuser, emailAddress = user@localhost.com