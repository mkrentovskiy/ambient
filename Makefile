all:
	GOARCH=arm GOARM=5 CGO_ENABLED=1 go build
	scp ambient root@192.168.222.125:/root
