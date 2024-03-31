generatescripts: 
	gopherjs build frontend/scripts/scripts.go -o frontend/scripts/scripts.js && \
	gopherjs build frontend/scriptsCustomers/scriptsCustomers.go -o frontend/scriptsCustomers/scriptsCustomers.js
generatescripts-apple: 
	GOOS=darwin GOARCH=arm64  gopherjs build frontend/scripts/scripts.go -o frontend/scripts/scripts.js && \
	gopherjs build frontend/scriptsCustomers/scriptsCustomers.go -o frontend/scriptsCustomers/scriptsCustomers.js
build: 
	make generatescripts && go build -o ./jobsManager
build-apple: 
	make generatescripts-apple && GOOS=darwin GOARCH=arm64 go build -o ./jobsManager
run:
	make generatescripts && go run server.go
