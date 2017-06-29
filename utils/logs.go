package utils

import (
	"sync"
	"time"
	"fmt"
)

var oneTime sync.Once
var logInstance *logging

const TRACE int = 0
const DEBUG int = 1
const INFO int = 2
const WARN int = 3
//const ERROR int = 4


type logging struct{

	logStream chan string
	level int
}

func GetLogging() *logging {
	oneTime.Do(func() {
		logInstance = &logging{}
		logInstance.init()
	})
	return logInstance
}

func (l *logging) init(){
		l.logStream = make(chan string)
		l.level = GetSettings().GetLevel();
		go l.logWriter()
}

func (l *logging) Trace(s string,v ...interface{}){
	if(l.level == TRACE){
		l.logStream <- "TRACE: "+fmt.Sprintf(s,v...)
	}
}
func (l *logging) Debug(s string,v ...interface{}){
	if(l.level <= DEBUG){
		l.logStream <- "DEBUG: "+fmt.Sprintf(s,v...)
	}
}
func (l *logging) Info(s string,v ...interface{}){
	if(l.level <= INFO){
		l.logStream <- "INFO: "+fmt.Sprintf(s,v...)
	}
}
func (l *logging) Warn(s string,v ...interface{}){
	if(l.level <= WARN){
		l.logStream <- "WARN: "+fmt.Sprintf(s,v...)
	}
}
func (l *logging) Error(s string,v ...interface{}){
	//Error aways gets logged
	l.logStream <- "ERROR: "+fmt.Sprintf(s,v...)

}

func (l *logging) logWriter(){

	for {
		statement :=<-l.logStream
		dt := time.Now().Format(time.RFC3339)
		println(dt+" "+statement)
	}
}