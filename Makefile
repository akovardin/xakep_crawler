BIN=xacrawler
${BIN}: main.go src/4gophers.com/*/*.go gpm
	go build -o ./bin/${BIN}

bin: main.go src/4gophers.com/*/*.go
	go build -o ./bin/${BIN}

gpm: main.go
	gpm install

clean:
	rm -rf ./bin/${BIN}
	rm -rf ./pkg/*