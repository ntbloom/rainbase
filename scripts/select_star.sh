# select_star.sh - watch the database in real time

DB=/tmp/rainbase.db

while true; do 
  clear
  date
  sqlite3 $DB "SELECT log.id, mappings.longname AS event, log.value, log.timestamp FROM log INNER JOIN mappings ON log.tag = mappings.id;"
  sleep 2
done
