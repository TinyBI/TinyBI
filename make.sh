#!/bin/bash
GOBIN=go
GOGET="go get"
TARGET=tinybi
MODDIR=mods
MODSOURCES=src/mod/*.go
WWWSUBDIRS=("www/public/cache")

#Build Operation;
build_mods (){
	if [ ! -d $MODDIR ];then
	        mkdir $MODDIR
	fi
	for f in $MODSOURCES; do
		$GOBIN build -buildmode=plugin -o $MODDIR/`basename -s .go $f`.so $f
	done
}

build (){
	#Build main execution;
	if [ ! -d pkg ]; then
		$GOGET github.com/go-sql-driver/mysql
		$GOGET github.com/go-xorm/xorm
		$GOGET github.com/chai2010/gettext-go/gettext
		$GOGET github.com/satori/go.uuid
		$GOGET github.com/jasonlvhit/gocron
		$GOGET github.com/jinzhu/now
	fi
	$GOBIN install $TARGET 
	#Build modules;
	build_mods
	#Runtime directories;
	for rdir in ${WWWSUBDIRS[*]}
	do
	    if [ ! -d $rdir ];then
	        mkdir $rdir
		fi
	done
}

#Clean Operation;
clean (){
	rm -rf bin
	rm -rf pkg
	rm -rf mods
}

#Show usage;
usage (){
	echo "Usage: build | build_mods | clean"
}

case "$1" in
	"build" )
		build;;
	"build_mods" )
		build_mods;;
	"clean" )
		clean;;
	* )
		usage;;
esac