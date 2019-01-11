BIN:=bin/shroo

all: rm libs
	go build -v -o ${BIN} && go vet

cross: rm cross-libs
	GOOS=linux GOARCH=arm GOARM=5 CGO_ENABLED=1 CC=/pitools/arm-bcm2708/gcc-linaro-arm-linux-gnueabihf-raspbian-x64/bin/arm-linux-gnueabihf-gcc go build -v -o ${BIN} && go vet

test: rm libs
	go build -v -o ${BIN} && go test -v && go vet

libs:
	+${MAKE} -C util
	+${MAKE} -C gpio
	+${MAKE} -C bme280
	+${MAKE} -C db
	+${MAKE} -C sheets

cross-libs:
	+${MAKE} -C util cross
	+${MAKE} -C gpio cross 
	+${MAKE} -C bme280 cross 
	+${MAKE} -C db cross
	+${MAKE} -C sheets cross
	+${MAKE} -C versionControl cross

docker-build:
	docker build . -t cxx

# docker-run: 
# 	docker run -it --rm --name "cx" -v $GOPATH:/go cxx bash

rm:
	rm -f ${BIN}
