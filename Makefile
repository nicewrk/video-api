all:
	go install ./...

run: all
	video-api
