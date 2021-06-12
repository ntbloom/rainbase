#!/bin/bash
# mosquitto_pub.sh -- publish mosquitto messages

set -x

HOSTNAME=68.183.114.124
PORT=8883
SSL=/etc/rainbase/ssl
CAFILE=$SSL/ca.crt
CERT=$SSL/client.crt
KEY=$SSL/client.key


mosquitto_pub \
  -u rainbase \
  -P $(cat $HOME/rainbase.pw) \
  --cafile  $CAFILE \
  --cert $CERT \
  --key $KEY \
  --tls-version tlsv1.2 \
  -h $HOSTNAME \
  -q 0 \
  -p $PORT \
  -t $1 \
  -m $2 \
  -d


