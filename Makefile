gen:
	go generate ./...

dk:
	docker build . -t bob:dev