/*Package utils ...
 *
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
	"github.com/namsral/flag"
	"log"
	"strings"
	"sync"
)

/**
Singleton settings type to share settings across the codebase
*/

var once sync.Once
var instance *settings

//Settings holds the values parsed out from flags
type settings struct {
	port                int
	dbString            string
	help                bool
	ssl                 bool
	sslKey              string
	sslcert             string
	cacheSizeMB         int
	dbs                 []string
	loglevel            int
	proxyHostname       string
	proxyPort           string
	proxyScheme         string
	enableProxySettings bool
}

//GetSettings returns the setting singleton
func GetSettings() *settings {
	once.Do(func() {
		instance = &settings{}
		instance.setup()

	})
	return instance
}
func (s settings) GetEnableProxySettings() bool {
	return s.enableProxySettings
}
func (s settings) GetHostname() string {
	return s.proxyHostname
}
func (s settings) GetProxyPort() string {
	return s.proxyPort
}
func (s settings) GetProxyScheme() string {
	return s.proxyScheme
}
func (s settings) GetPort() int {
	return s.port
}
func (s settings) GetCacheSizeMB() int {
	return s.cacheSizeMB
}
func (s settings) GetSsl() bool {
	return s.ssl
}
func (s settings) GetHelp() bool {
	return s.help
}
func (s settings) GetDBs() []string {
	return s.dbs
}
func (s settings) GetSslKey() string {
	return s.sslKey
}
func (s settings) GetSslCert() string {
	return s.sslcert
}
func (s settings) GetLevel() int {
	return s.loglevel
}
func (s *settings) setup() {

	if flag.Parsed() {
		return
	}
	var level string

	flag.StringVar(&s.dbString, "db", "resources/ireland.mbtiles", "The MBTiles Database -- multiple can be specified by separated the paths via commas")
	flag.BoolVar(&s.ssl, "ssl", true, "Whether to use SSL -- disabling SSL will also disable HTTP2 -- enabled by default")
	flag.StringVar(&s.sslKey, "key", "resources/test.key", "The ssl private key")
	flag.StringVar(&s.sslcert, "cert", "resources/test.crt", "The ssl private cert")
	flag.StringVar(&level, "loglevel", "INFO", "Log level - one of TRACE, DEBUG, INFO, WARN, ERROR")
	flag.IntVar(&s.port, "port", 8000, "The port number")
	flag.BoolVar(&s.help, "help", false, "This message")
	flag.IntVar(&s.cacheSizeMB, "cacheSize", 200, "The size of the in memory cache (in MB)")
	flag.BoolVar(&s.enableProxySettings, "proxy", false, "For Proxies -- Whether to enable proxy settings or just use whatever hostname is found in the http request headers")
	flag.StringVar(&s.proxyHostname, "proxyhostname", "localhost", "For Proxies -- The hostname that should be advertised")
	flag.StringVar(&s.proxyPort, "proxyport", "8000", "For Proxies -- The port that should be advertised")
	flag.StringVar(&s.proxyScheme, "proxyscheme", "https", "For Proxies -- The hostname that should be advertised")

	flag.Parse()
	//Split DBs
	log.Printf("Log level: %v", level)
	s.dbs = make([]string, 0)
	if len(s.dbString) == 0 {
		log.Fatal("No Database specified! use -db to set one....")
	}
	if strings.Contains(s.dbString, ",") {
		log.Printf("here")
		s.dbs = strings.Split(s.dbString, ",")
	} else {
		s.dbs = append(s.dbs, s.dbString)
	}
	switch level {
	case "TRACE":
		s.loglevel = 0
		break
	case "DEBUG":
		s.loglevel = 1
		break
	case "INFO":
		s.loglevel = 2
		break
	case "WARN":
		s.loglevel = 3
		break
	case "ERROR":
		s.loglevel = 4
		break
	}

}
