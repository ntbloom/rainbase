SERVER=68.183.114.124
PORT=8883
FULL=$SERVER:$PORT

debug()
{
  openssl s_client \
    -crlf \
    -connect $FULL \
    -servername $SERVER \
    -psk_identity rainbase \
    -psk $(cat $HOME/rainbase.psk) \
    -pass file:$HOME/rainbase.pw
}

debug
