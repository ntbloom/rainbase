#!/bin/bash
# mosquitto_pub.sh -- publish mosquitto messages

set -x

HOSTNAME=68.183.114.124
PORT=8883


mosquitto_pub \
  -u rainbase \
  -P $(cat $HOME/rainbase.pw) \
  -h $HOSTNAME \
  -q 0 \
  -p $PORT \
  -t $1 \
  -m $2


