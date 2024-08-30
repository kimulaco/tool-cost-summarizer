test:
	mkdir .coverage
	go test -cover ./... -race -coverprofile=.coverage/coverage.out

coverage-html:
	go tool cover -html .coverage/coverage.out -o .coverage/coverage.html
