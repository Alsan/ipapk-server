#!/bin/bash

ip=$1

if [[ ! -d ".ca" ]]; then
    mkdir -p .ca && cd .ca
    openssl genrsa -out myCA.key 2048 2> "/dev/null"
    openssl req -x509 -new -key myCA.key -out myCA.cer -days 730 -subj /CN="ipapk-server "$ip" Custom CA" 2> "/dev/null"
    openssl genrsa -out mycert.key 2048 2> "/dev/null"
    openssl req -new -out mycert.req -key mycert.key -subj /CN=$ip 2> "/dev/null"
    openssl x509 -req -in mycert.req -out mycert.cer -CAkey myCA.key -CA myCA.cer -days 365 -CAcreateserial -CAserial serial 2> "/dev/null"
fi

