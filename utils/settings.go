package utils

import (
	"flag"
	"sync"
	"strings"
	"log"
)

/**
Singleton settings type to share settings across the codebase
 */

var once sync.Once
var instance *settings

type settings struct{
	port int
	dbString string
	help bool
	ssl bool
	sslKey string
	sslcert string
	cacheSizeMB int
	dbs []string
	loglevel int

}
func GetSettings() *settings {
	once.Do(func() {
		instance =  &settings{}
		instance.init()
	})
	return instance
}
func (s *settings) GetPort() (int){
	return s.port
}
func (s *settings) GetCacheSizeMB() (int){
	return s.cacheSizeMB
}
func (s *settings) GetSsl() (bool){
	return s.ssl
}
func (s *settings) GetHelp() (bool){
	return s.help
}
func (s *settings) GetDBs() ([]string){
	return s.dbs
}
func (s *settings) GetSslKey() (string){
	return s.sslKey
}
func (s *settings) GetSslCert() (string){
	return s.sslcert
}
func (s *settings) GetLevel() (int){
	return s.loglevel
}
func (s *settings) init(){

	var level string
	flag.StringVar(&s.dbString,"db", "resources/ireland.mbtiles", "The MBTiles Database -- multiple can be specified by seperating the paths via commas")
	flag.BoolVar(&s.ssl,"ssl", true, "Whether to use SSL -- disabling SSL will also disable HTTP2 -- enabled by default")
	flag.StringVar(&s.sslKey,"key", "resources/test.key", "The ssl private key")
	flag.StringVar(&s.sslcert,"cert", "resources/test.crt", "The ssl private cert")
	flag.StringVar(&level,"loglevel", "INFO", "Log level - one of TRACE, DEBUG, INFO, WARN, ERROR")
	flag.IntVar(&s.port,"port", 8000,"The port number")
	flag.BoolVar(&s.help,"help",false,"This message")
	flag.IntVar(&s.cacheSizeMB,"cacheSize", 200,"The size of the in memeory cache (in MB)")
	flag.Parse()
	//Split DBs
	log.Printf("Log level: %v",level)
	s.dbs = make([]string,0)
	if(strings.Contains(s.dbString,",")) {
		s.dbs = strings.Split(s.dbString, ",")
	} else {
		s.dbs = append(s.dbs, s.dbString)
	}
	switch level {
	case "TRACE":
		s.loglevel = 0
		break;
	case "DEBUG":
		s.loglevel = 1
		break;
	case "INFO":
		s.loglevel = 2
		break;
	case "WARN":
		s.loglevel = 3
		break;
	case "ERROR":
		s.loglevel = 4
		break;
	}
}


