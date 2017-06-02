#!/usr/bin/env bash
# Shamelessly copy/pasted from
# https://gist.github.com/denji/12b3a568f092ab951456


# Key considerations for algorithm "RSA" ≥ 2048-bit
openssl genrsa -out ssl/key.pem 2048

# Key considerations for algorithm "ECDSA" ≥ secp384r1
# List ECDSA the supported curves (openssl ecparam -list_curves)
openssl ecparam -genkey -name secp384r1 -out ssl/key.pem

openssl req -new -x509 -sha256 -key ssl/key.pem -out ssl/cert.crt -days 3650
