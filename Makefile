
build: main.go
	go build -o bin/youtube-dl-multiconfig main.go

run: build
	bin/youtube-dl-multiconfig --config example.json --simulate

realrun: build
	bin/youtube-dl-multiconfig --config example.json