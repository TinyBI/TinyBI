#!/bin/bash
GOBIN=go
TARGET=tinybi
MODDIR=mods
MODSOURCES=src/mod/*.go
WWWSUBDIRS=("www/public/cache")

#Build Operation;
build (){
	#Build main execution;
	$GOBIN install $TARGET 
	#Build modules;
	if [ ! -d $MODDIR ];then
	        mkdir $MODDIR
	fi
	for f in $MODSOURCES; do
		$GOBIN build -buildmode=plugin -o $MODDIR/`basename -s .go $f`.so $f
	done
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
	echo "Usage: build | clean"
}

case "$1" in
	"build" )
		build;;
	"clean" )
		clean;;
	* )
		usage;;
esac