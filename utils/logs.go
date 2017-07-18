/*
 * Copyright 2017-present Tom Ingold / Ruptive.io
 * Licensed under the Apache License, Version 2.0 (the "License").
 * You may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */
package utils


import (
	"sync"
	"time"
	"fmt"
	"os"
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
func (l *logging) Fatal(s string,v ...interface{}){

	dt := time.Now().Format(time.RFC3339)
	println(dt+"  FATAL: "+fmt.Sprintf(s,v...))
	os.Exit(1)
}

func (l *logging) logWriter(){

	for {
		statement :=<-l.logStream
		dt := time.Now().Format(time.RFC3339)
		println(dt+" "+statement)
	}
}