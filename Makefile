binary = gcp-connectivitytests
version = 0.0.1

release:
	GOOS=linux GOARCH=amd64 go build -o ./bin/$(binary)-$(version)

push:
	gsutil cp bin/* gs://jbd-releases
