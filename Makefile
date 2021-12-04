
build: main.go
	go build -o bin/ytdlmc main.go

run: build
	bin/ytdlmc --config example.json --simulate

realrun: build
	bin/ytdlmc --config example.json