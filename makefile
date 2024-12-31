PHONY: run-local

run-local:
	go build src/winter/app.go
	LOCAL_CG=true ./app
