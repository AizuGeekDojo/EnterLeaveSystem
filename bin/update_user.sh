SERVER=""
printf "Input your ID: "
read USERID

ssh ${USERID}@${SERVER} getent passwd |awk -F: '{print $1 "|" $5 "|0"}' >users_sqlite

sqlite3 ../database.db << EOF
delete from users;
.import users_sqlite users
EOF
