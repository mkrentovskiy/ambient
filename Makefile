TOOLCHAINS_DIR = /home/mkrentovskiy/develop/virt2real/install-sdk/codesourcery/arm-2013.05
SYSROOT = /home/mkrentovskiy/develop/virt2real/virt2real-sdk.lxc/fs/output/host/usr/arm-buildroot-linux-gnueabi/sysroot

OPTS = GOARCH=arm
OPTS += GOARM=5 
OPTS += GOOS=linux
OPTS += CGO_ENABLED=1 
OPTS += CC="arm-none-linux-gnueabi-gcc"
OPTS += CXX=arm-none-linux-gnueabi-g++
OPTS += PATH=$(PATH):$(TOOLCHAINS_DIR)/bin 
OPTS += PKG_CONFIG_DIR=""
OPTS += PKG_CONFIG_LIBDIR=$(SYSROOT)/usr/lib/pkgconfig
OPTS += PKG_CONFIG_SYSROOT_DIR=$(SYSROOT)
OPTS += CGO_CFLAGS="-I$(SYSROOT)/usr/include"
OPTS += CGO_LDFLAGS="-L$(SYSROOT)/usr/lib -Wl,-rpath,$(SYSROOT)/usr/lib --sysroot=$(SYSROOT)"

all:
	$(OPTS) go env
	$(OPTS) go build
	scp ambient root@192.168.3.1:/root
