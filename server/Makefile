run: main.go
	export GO111MODULE=on && cd test && go run ../main.go

build: main.go
#	GOOS=linux GOARCH=arm go build -o ../elsystem main.go
	xgo --targets=linux/arm-7 --branch master --pkg cmd/ github.com/AizuGeekDojo/EnterLeaveSystem

initdb:
	cd test && rm -f database.db && sqlite3 database.db < init.sql