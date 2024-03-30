generatescripts: 
	gopherjs build scripts/scripts.go -o scripts/scripts.js
build: 
	go build -o ./jobsManager
run:
	make generatescripts && go run server.go
