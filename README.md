# Matterhorn Bot

This is a multi-use telegram bot written in Go.

## Generate Webhook Certs

    mkdir -p ignored && openssl req -x509 -newkey rsa:2048 -keyout ignored/key.key -out ignored/cert.pem -days 3560 -subj "//O=Org\CN=Test" -nodes

