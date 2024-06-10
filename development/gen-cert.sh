#!/bin/bash
rm *.pem

SUBJECT="/C=VN/ST=Vietnam/L=Hochiminh/O=GXBank/OU=Banking/CN=*.gft.com/emailAddress=dat.le@gft.com"

# Files
CA_CERT=./certs/ca-cert.pem
CA_KEY=./certs/ca-key.pem

CLIENT_REQ="./certs/client-req.pem"
CLIENT_KEY="./certs/client-key.pem"
CLIENT_CERT="./certs/client-cert.pem"

SERVER_REQ=./certs/server-req.pem
SERVER_KEY=./certs/server-key.pem
SERVER_CERT=./certs/server-cert.pem

SERVER_CONF=./certs/server-ext.conf
CLIENT_CONF=./certs/client-ext.conf

# 1. Generate CA's private key and self-signed certificate
openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout ${CA_KEY} -out ${CA_CERT} -subj ${SUBJECT}

echo "CA's self-signed certificate"
openssl x509 -in ${CA_CERT} -noout -text

# 2. Generate web server's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout ${SERVER_KEY} -out ${SERVER_REQ} -subj ${SUBJECT}

# 3. Use CA's private key to sign web server's CSR and get back the signed certificate
openssl x509 -req -in ${SERVER_REQ} -days 60 -CA ${CA_CERT} -CAkey ${CA_KEY} -CAcreateserial -out ${SERVER_CERT} -extfile ${SERVER_CONF}

echo "Server's signed certificate"
openssl x509 -in ${SERVER_CERT} -noout -text

# 4. Generate client's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout ${CLIENT_KEY} -out ${CLIENT_REQ} -subj ${SUBJECT}

# 5. Use CA's private key to sign client's CSR and get back the signed certificate
openssl x509 -req -in ${CLIENT_REQ} -days 60 -CA ${CA_CERT} -CAkey ${CA_KEY} -CAcreateserial -out ${CLIENT_CERT} -extfile ${CLIENT_CONF}

echo "Client's signed certificate"
openssl x509 -in ${CLIENT_CERT} -noout -text
