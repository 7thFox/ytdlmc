
build: main.go
	go build -o bin/ytdlmc main.go

run: build
	bin/ytdlmc --config example.json --simulate

realrun: build
	bin/ytdlmc --config example.json

docker:
	docker build -t ytdlmc .

docker-test: docker-destroy docker
	sudo rm -rf $(CURDIR)/test/downloads $(CURDIR)/test/log/ytdlmc.log $(CURDIR)/test/log/ytdlmc.lock
	mkdir $(CURDIR)/test/downloads

	docker create \
		--name ytdlmc \
		-e CRONTIME="*/1 * * * *" \
		-v $(CURDIR)/test/appdata:/opt/appdata \
		-v $(CURDIR)/test/log:/opt/log \
		-v $(CURDIR)/test/downloads:/opt/downloads \
		ytdlmc

	docker start ytdlmc

docker-destroy:
	docker stop ytdlmc || true
	docker rm ytdlmc || true
