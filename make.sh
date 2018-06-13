#!/bin/bash
GOBIN=go
GOGET="go get"
GOBUILD="go build"
TARGET=tinybi_exec
BINDIR="bin"
PKGDIR="pkg"
SCRIPTSDIR="scripts"
GOPATH=`pwd`
DEBUG="false"
MODDIR="mods"
MODSOURCES="src/tinybi_mods"
SRC_ARCHIVE="tinybi.src.tar.gz"

#Build Operation;
#Production build;
build(){
	#Build main execution;
	export GOPATH
	if [ ! -d $PKGDIR ]; then
		$GOGET github.com/go-sql-driver/mysql
		$GOGET github.com/go-xorm/xorm
		$GOGET github.com/jasonlvhit/gocron
		#$GOGET github.com/gljubojevic/gocron
		$GOGET github.com/360EntSecGroup-Skylar/excelize
		$GOGET github.com/go-gomail/gomail
		$GOGET github.com/chai2010/gettext-go/gettext
	fi
	if [ $DEBUG == "true" ]; then
		$GOBIN install $TARGET 
	else
		$GOBIN install -ldflags="-w -s" $TARGET
	fi 
	#Copy scripts into bin
	#cp $SCRIPTSDIR/* $BINDIR
}

build_mods(){
    #Build web mods;
    export GOPATH
    if [ ! -d $MODDIR ];then
        mkdir $MODDIR
    fi
    for mDir in $MODSOURCES/*
    do
        if test -d $mDir; then
            if [ $DEBUG == "true" ]; then
                $GOBUILD -buildmode=plugin -o $MODDIR/`basename $mDir`.so $mDir/main.go
            else
                $GOBUILD -ldflags="-w -s" -buildmode=plugin -o $MODDIR/`basename $mDir`.so $mDir/main.go
            fi
        fi
    done
}

#Clean Operation;
clean (){
	rm -rf $BINDIR
	rm -rf $PKGDIR
}

#Show usage
help() {
	echo "Build script for tinybi"
	echo "Usage:"
	echo "    make.sh command"
	echo "The commands are:"
	echo "    build        :create binary with debug information"
	echo "    dist         :create binary without debug information"
	echo "    build_mods   :create mod binaries with debug information"
    echo "    dist_mods    :create mod binaries without debug information"
	echo "    src_archive  :create archive for source which is used for installation"
	echo "    install_src  :install source file from archive"
	echo "    clean        :clean built files"
	echo "    help         :show this usage"
}

#Compress source to tar.gz;
src_archive() {
	tar -zcvf $SRC_ARCHIVE src
}

#install source;
install_src() {
	if [ ! -f $SRC_ARCHIVE ]; then
		echo "Cannot find archive of source:"$SRC_ARCHIVE
	else
		rm -rf src
		tar -zxvf $SRC_ARCHIVE
	fi
}

case "$1" in
	"build" )
		DEBUG="true"
		build;;
	"dist" )
		DEBUG="false"
		build;;
	"build_mods" )
	    DEBUG="true"
	    build_mods;;
	"dist_mods" )
	    DEBUG="false"
	    build_mods;;
	"src_archive" )
		src_archive;;
	"install_src" )
		install_src;;
	"clean" )
		clean;;
	* )
		help;;
esac
