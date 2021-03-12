# select_star.sh - watch the database in real time

DB=/tmp/rainbase.db

while true; do 
  clear
  date
  sqlite3 $DB "SELECT * FROM log;"
  sleep 2
done
