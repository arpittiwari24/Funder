run: 
	go mod tidy
	nodemon --exec go run main.go --signal SIGTERM

test:
	go test -v ./...

push:
	git add .
	git commit -m "update"
	git push origin main