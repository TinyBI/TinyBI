#!/bin/bash
GOBIN=go
GOGET="go get"
GOBUILD="go build"
TARGET=tinybi
MODDIR="mods"
MODSOURCES="src/mod"
WWWSUBDIRS=("www/public/cache")
GOPATH=`pwd`

#Build Operation;
build_mods(){
    #Build web mods;
    export GOPATH
    if [ ! -d $MODDIR ];then
        mkdir $MODDIR
    fi
    for mDir in $MODSOURCES/*
    do
        if test -d $mDir; then
            $GOBUILD -buildmode=plugin -o $MODDIR/`basename $mDir`.so $mDir/main.go
        fi
    done
}

build (){
	#Build main execution;
	export GOPATH
	if [ ! -d pkg ]; then
		$GOGET github.com/go-sql-driver/mysql
		$GOGET github.com/go-xorm/xorm
		$GOGET github.com/chai2010/gettext-go/gettext
		$GOGET github.com/satori/go.uuid
		$GOGET github.com/jasonlvhit/gocron
		$GOGET github.com/jinzhu/now
		$GOGET github.com/go-gomail/gomail
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