all:
	go build -buildmode=plugin -o ${DIST_PLUG}/gpio.so 
cross:
	GOOS=linux GOARCH=arm GOARM=6 CGO_ENABLED=1 CC=/pitools/arm-bcm2708/gcc-linaro-arm-linux-gnueabihf-raspbian-x64/bin/arm-linux-gnueabihf-gcc go build -buildmode=plugin -o ${DIST_PLUG}/gpio.so 
