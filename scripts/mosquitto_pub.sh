HOSTNAME=68.183.114.124
PORT=8883
SSLDIR=/etc/raincloud/ssl

CAFILE=$SSLDIR/ca.crt
CERT=$SSL/client.crt
KEY=$SSL/client.key

mosquitto_pub \
  --insecure \
  --tls-version tls1.2 \
  --cafile $CAFILE \
  --cert $CERT \
  --key $KEY \
  -h $HOSTNAME \
  -p $PORT \
  -t $1 \
  -m $2


