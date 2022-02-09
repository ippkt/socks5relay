MAJOR_VER=$(shell date +%y%m%d)
MINOR_VER=1

all:
	go build -ldflags "-s -w -X main.version=${MAJOR_VER}.${MINOR_VER}"   
