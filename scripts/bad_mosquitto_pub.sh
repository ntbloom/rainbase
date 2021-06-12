#!/bin/bash
# mosquitto_pub.sh -- publish mosquitto messages

set -x

HOSTNAME=68.183.114.124
PORT=8883

echo "THIS SHOULD NOT WORK!"
mosquitto_pub \
  -u rainbase \
  -h $HOSTNAME \
  -q 0 \
  -p $PORT \
  -t hello \
  -m world


mosquitto_pub \
  -u rainbase \
  -P bad_passwd \
  -h $HOSTNAME \
  -q 0 \
  -p $PORT \
  -t hello \
  -m world
