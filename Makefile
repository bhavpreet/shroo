BIN:=bin/shroo
CWD:=$(shell pwd)
DIST:=${CWD}/dist
DIST_PLUG:=${DIST}/plugins

all: rm dist libs
	go build -v -o ${DIST}/${BIN} && go vet
	cp -v sample-config.yml ${DIST}/config.yml


cross: rm dist cross-libs
	GOOS=linux GOARCH=arm GOARM=6 CGO_ENABLED=1 CC=/pitools/arm-bcm2708/gcc-linaro-arm-linux-gnueabihf-raspbian-x64/bin/arm-linux-gnueabihf-gcc go build -v -o ${DIST}/${BIN} && go vet
	cp -v config.yml ${DIST}/config.yml

dist:
	mkdir -p ${DIST_PLUG}

test: rm dist libs
	go build -v -o ${DIST}/${BIN} && go test -v && go vet

libs: dist
	+${MAKE} DIST_PLUG=${DIST_PLUG} -C util
	+${MAKE} DIST_PLUG=${DIST_PLUG} -C gpio
	+${MAKE} DIST_PLUG=${DIST_PLUG} -C bme280
	+${MAKE} DIST_PLUG=${DIST_PLUG} -C aht10
	+${MAKE} DIST_PLUG=${DIST_PLUG} -C sheets
	# +${MAKE} -C db

cross-libs: dist
	+${MAKE} DIST_PLUG=${DIST_PLUG} -C util cross
	+${MAKE} DIST_PLUG=${DIST_PLUG} -C gpio cross
	+${MAKE} DIST_PLUG=${DIST_PLUG} -C bme280 cross
	+${MAKE} DIST_PLUG=${DIST_PLUG} -C aht10 cross
	+${MAKE} DIST_PLUG=${DIST_PLUG} -C sheets cross
	# +${MAKE} -C db cross
	# +${MAKE} -C versionControl cross

docker-build:
	docker build . -t cxx

docker-run: 
	docker run -it --rm --name "cx" -v ${GOPATH}:/go cxx bash

rm:
	rm -rf ${DIST}
