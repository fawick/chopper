package handler

/* Licensed under the Apache License, Version 2.0 (the "License").
 * You may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */
import (
	"encoding/json"
	"github.com/boundlessgeo/chopper/tiles"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

//InfoHandler exposes information about how the server is currently configured
type InfoHandler struct {
	Tm *tiles.TileManager
}

//Handle implements the http server interface
func (infoHandler *InfoHandler) Handle(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	out, _ := json.Marshal(infoHandler.Tm.Metadatas)
	w.Header().Add("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(200)
	w.Write(out)

}
