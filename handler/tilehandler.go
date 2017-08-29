/*Package handler ...
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

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/tingold/gophertile/gophertile"
	"github.com/ruptivespatial/chopper/utils"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"github.com/ruptivespatial/chopper/tiles"
)

//Tilehandler implements the handle method and takes care of http crap while delegating to the tilemanager for actual
//tile retrieval
type Tilehandler struct {
	Manager tiles.TileManager
}

//Handle implements the httprouter method...
func (th *Tilehandler) Handle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	y,_ := strconv.Atoi(strings.TrimSuffix(ps.ByName("y"), ".pbf"))
	z,_ := strconv.Atoi(ps.ByName("z"))
	x, _:= strconv.Atoi(ps.ByName("x"))
	yInt := normalizeY(y, z)

	t, data := th.Manager.GetTile(z, x, int(yInt))


	if data == nil {
		utils.GetLogging().Warn(fmt.Sprintf("Tile not found for %v/%v/%v", ps.ByName("z"), ps.ByName("x"), yInt))

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(404)
	} else {
		w.Header().Set("Content-type", "application/x-protobuf")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		var buff = bytes.NewBuffer(data)
		r, err := gzip.NewReader(buff)
		if err != nil {
			utils.GetLogging().Error("error decompressing tile")
			w.WriteHeader(500)
			return
		}
		io.Copy(w, r)
		w.WriteHeader(200)
		if pusher, ok := w.(http.Pusher); ok {
			//log.Print("HTTP Push is OK")
			options := &http.PushOptions{}
			kids := t.Children()

			for _, kid := range kids {
				if kid == nil {
					continue
				}
				url := buildUrl(t)
				pusher.Push(url, options)
			}
		}

	}
}
func buildUrl(tile *gophertile.Tile) string{
	//return "/tiles/" + t.ZStr() + "/" + t.XStr() + "/" + t.YStr() + ".pbf"
	buf := bytes.NewBufferString("/tiles")
	buf.WriteString(strconv.Itoa(tile.Z))
	buf.WriteString("/")
	buf.WriteString(strconv.Itoa(tile.X))
	buf.WriteString("/")
	buf.WriteString(strconv.Itoa(tile.Y))
	buf.WriteString(".pbf")
	return buf.String()
}

func normalizeY(y int, z int) int32 {

	floaty := math.Pow(float64(2.0), float64(z)) - float64(y)
	floaty--

	return int32(floaty)
}
