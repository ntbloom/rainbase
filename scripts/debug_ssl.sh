openssl s_client \
  -connect 68.183.114.124:8883 \
  -psk_identity rainbase \
  -psk $(cat $HOME/rainbase.psk)
