#!/usr/bin/env bash

BOT_TOKEN="$1"
WEBHOOK_URL="$2"

curl -v -XGET "https://api.telegram.org/bot${BOT_TOKEN}/setWebhook?url=${WEBHOOK_URL}"
