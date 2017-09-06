BuildStamp=main.BuildStamp=`date '+%Y-%m-%d_%I:%M:%S%p'`
GitHash=main.GitHash=`git rev-parse HEAD`

all:
	go run main.go

test:
	go test .

authors:
	echo "Authors\n=======\n\nProject's contributors:\n" > AUTHORS.md
	git log --raw | grep "^Author: " | cut -d ' ' -f2- | cut -d '<' -f1 | sed 's/^/- /' | sort | uniq >> AUTHORS.md

build: clean
	mkdir release
	go build -o release/server -ldflags "-s -w -X ${BuildStamp} -X ${GitHash}" main.go

clean:
	-rm -rf release
