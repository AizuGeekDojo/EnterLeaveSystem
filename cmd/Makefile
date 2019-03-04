run: main.go
	cd test && go run ../main.go

build: main.go
#	GOOS=linux GOARCH=arm go build -o ../elsystem main.go
	xgo --targets=linux/arm-7 --branch go_server --pkg cmd/ github.com/AizuGeekDojo/EnterLeaveSystem

