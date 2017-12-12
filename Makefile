##This file builds golang binaries to the bin directory and also allows for cleaning them up.
all: darwin linux

darwin:
	env GOOS=darwin GOARCH=amd64 go build -o bin/container-shifter .
	chmod +x bin/container-shifter
linux:
	env GOOS=linux GOARCH=amd64 go build -o bin/container-shifter .
	chmod +x bin/container-shifter
clean:
	rm -rf ./bin
