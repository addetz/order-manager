generatescripts: 
	gopherjs build scripts/scripts.go -o scripts/scripts.js && \
	gopherjs build scripts-customers/scripts-customers.go -o scripts-customers/scripts-customers.js
build: 
	go build -o ./jobsManager
build-apple: 
	GOOS=darwin GOARCH=arm64 go build -o ./jobsManager
run:
	make generatescripts && go run server.go
