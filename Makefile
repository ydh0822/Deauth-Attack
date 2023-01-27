NAME=Deauth-Attack

all: deps build

deps:

build:
	go build -o ${NAME} main.go

clean:
	go clean
	rm ${NAME}