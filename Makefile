BINARY=trello-backup

build:
	go build -o bin/${BINARY} main.go

clean:
	rm -rf dist/*

compile:
	GOOS=darwin GOARCH=amd64 go build -o ./dist/darwin-amd64/${BINARY} main.go
	GOOS=linux GOARCH=amd64 go build -o ./dist/linux-amd64/${BINARY} main.go
	GOOS=windows GOARCH=amd64 go build -o ./dist/windows-amd64/${BINARY}.exe main.go

pack:
	cp ./templates/trello-backup.yaml ./dist/darwin-amd64/
	zip -j ./dist/darwin-amd64.zip ./dist/darwin-amd64/*

	cp ./templates/trello-backup.yaml ./dist/linux-amd64/
	zip -j ./dist/linux-amd64.zip ./dist/linux-amd64/*

	cp ./templates/trello-backup.yaml ./dist/windows-amd64/
	cp ./scripts/windows/*.bat ./dist/windows-amd64/
	zip -j ./dist/windows-amd64.zip ./dist/windows-amd64/*

all: compile pack