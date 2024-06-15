# BEAVER-IDP

## Create Key for signing JWT

`openssl genpkey -algorithm RSA -out private.pem`

` openssl rsa -pubout -in private.pem -out public.key`

## ENV Variables

- PORT
- DB
- PRIVATE_KEY
