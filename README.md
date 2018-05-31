# TinyBI

[![Travis Build Status](https://travis-ci.org/TinyBI/TinyBI.svg?branch=master)](https://github.com/TinyBI/TinyBI)
[![License](http://img.shields.io/github/license/TinyBI/TinyBI.svg)](https://github.com/TinyBI/TinyBI)
[![Go Report Card](https://goreportcard.com/badge/github.com/TinyBI/TinyBI)](https://github.com/TinyBI/TinyBI)
[![Openource](http://img.shields.io/badge/opensource-BI%20Report%20System-orange.svg)](https://github.com/TinyBI/TinyBI)


Openource BI Report System, simple and powerful

### Screenshots
- Login Page

![Login Page](https://github.com/TinyBI/TinyBI/raw/master/screenshots/login.png "Login Page")

- Concurrent Tasks

![Concurrent Tasks](https://github.com/TinyBI/TinyBI/raw/master/screenshots/currentTasks.png "Concurrent Tasks")

- Version Page, loaded by module 

![Version Page](https://github.com/TinyBI/TinyBI/raw/master/screenshots/aboutModule.png "Version Page")

### Features
- Modulized Architecture (Done)
- Scheduled Concurrent Tasks (Done)
- RESTFUL APIs (In process)
- Multilingual WEB UI (Done)

### Business Modules
- General Ledger (Done, Will by splited from core)

### Installation Dependencies
github.com/go-sql-driver/mysql

github.com/go-xorm/xorm

github.com/gljubojevic/gocron

github.com/360EntSecGroup-Skylar/excelize

github.com/go-gomail/gomail

github.com/chai2010/gettext-go/gettext

### Installation process
- To build TinyBI, use make.sh to simplize the Installation process
```Bash
#Just enter make.sh to show usage;
./make.sh
#Build the core
./make.sh build
#Build the modules
./make.sh build_mods
```