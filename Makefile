# Binary name
BINARY=chatroom
VERSION = 0.1
# Builds the project
build:
		go build -o ${BINARY}
		go test -v
# Installs our project: copies binaries
install:
		go install
release:
		cd bin
		# Clean
		go clean
		rm -rf *.gz
		# Build for mac
		cd ../client
		go build
		mv client ../bin
		cd ../server
		go build
		mv server ../bin
		cd ../bin
		tar czvf ../release/chatroom-mac64-1.0.tar.gz *
		# Build for linux
		go clean
		rm -rf *.gz
		cd ../client
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
		mv client ../bin
		cd ../server
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
		mv server ../bin
		cd ../bin
		tar czvf ../release/chatroom-linux64-1.0.tar.gz *
		# Build for win
		go clean
		rm -rf *.gz
		cd ../client
		CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build
		mv client ../bin
		cd ../server
		CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build
		mv server ../bin
		cd ../bin
		tar czvf ../release/chatroom-win64-1.0.tar.gz *
		go clean
# Cleans our projects: deletes binaries
clean:
		go clean

.PHONY:  clean build