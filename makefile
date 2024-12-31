PHONY: run-local

run-local:
	go build src/winter/app.go
	LOCAL_CG=true ./app

# to run profiler on this app
profile:
	# run the app and start pprof on port 6060
	make run-local
	# run the profiler using cpu.prof file
	go tool pprof -http=:8080 app cpu.prof
