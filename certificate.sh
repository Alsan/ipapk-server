#!/bin/bash

ip=$1

mkdir -p .ca && cd .ca
openssl genrsa -out install.key 2048 2> "/dev/null"
openssl req -x509 -new -key install.key -out install.cer -days 730 -subj /CN="IPAPK Generated CA "$ip"" 2> "/dev/null"
openssl genrsa -out server.key 2048 2> "/dev/null"
openssl req -new -out server.req -key server.key -subj /CN=$ip 2> "/dev/null"
openssl x509 -req -in server.req -out server.cer -CAkey install.key -CA install.cer -days 365 -CAcreateserial -CAserial serial 2> "/dev/null"

