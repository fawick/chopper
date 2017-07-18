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

package handler

import ("github.com/julienschmidt/httprouter"
        "net/http"
	"github.com/ruptivespatial/chopper/utils"
	"fmt"
	"strings"
	"github.com/elazarl/go-bindata-assetfs"

)

type proxyHostHandler struct{
	bdfs *assetfs.AssetFS
}
func NewProxyHostHandler(fs *assetfs.AssetFS) (*proxyHostHandler) {
	phh := new(proxyHostHandler)
	phh.bdfs = fs
	return phh
}

func (phh *proxyHostHandler) Handle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {


	logger := utils.GetLogging()
	//data, err := phh.bdfs.Asset(strings.TrimPrefix(r.URL.Path,"/"))
	data, err := phh.bdfs.Asset("static_source"+r.URL.Path)
	if(err != nil){
		logger.Warn("File not found: %v",err);
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(404)
		fmt.Fprintln(w, "Not Found")
		return
	}
	n := len(data)
	fileString := string(data[:n])

	var newhostname string
	if(utils.GetSettings().GetEnableProxySettings()) {
		newhostname = utils.GetSettings().GetProxyScheme() + "://" + utils.GetSettings().GetHostname() + ":" + utils.GetSettings().GetProxyPort()
		logger.Debug("Using PROXY values to rewrite json")
	} else{
		logger.Debug("Using request values to rewrite json")
		if(utils.GetSettings().GetSsl()){
			newhostname = "https://"
		} else{
			newhostname = "http://"
		}
		newhostname += r.Host
	}

	logger.Debug("Setting base URLs to %v",newhostname)
	reformattedString := strings.Replace(fileString,"https://localhost:8000",newhostname,-1);
	w.Write([]byte(reformattedString))


}