run:
    pipenv run python server.py
migrate-up:
    sqlite3 database.db < inittable
migrate-insert:
    sqlite3 database.db < initusers
migrate-down:
    echo "" > database.db


