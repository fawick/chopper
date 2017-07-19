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

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/ruptivespatial/chopper/tiles"
	"github.com/ruptivespatial/chopper/utils"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
)

type Tilehandler struct {
	Manager tiles.TileManager
}

func (th *Tilehandler) Handle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var y string = strings.TrimSuffix(ps.ByName("y"), ".pbf")
	var z string = ps.ByName("z")
	yInt := normalizeY(y, z)

	t := th.Manager.GetTile(z, ps.ByName("x"), strconv.Itoa(int(yInt)))

	if t.Data == nil {
		utils.GetLogging().Warn(fmt.Sprintf("Tile not found for %v/%v/%v", ps.ByName("z"), ps.ByName("x"), yInt))

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(404)
	} else {
		w.Header().Set("Content-type", "application/x-protobuf")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		var buff = bytes.NewBuffer(t.Data)
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
			adjecentArray := tiles.GetZoomLevelManager().GetAdjacentTiles(t)

			for _, adjtile := range adjecentArray {
				if adjtile == nil {
					continue
				}
				url := adjtile.GetUrl()
				//log.Printf("Pushing tile %v",url)
				pusher.Push(url, options)
			}
		}

	}

}
func normalizeY(whyStr string, zStr string) int32 {

	z, err := strconv.Atoi(zStr)
	if err != nil {
		log.Printf("error converting val: %v to int", zStr)
		return -1
	}
	y, error := strconv.Atoi(whyStr)
	if error != nil {
		log.Printf("error converting val: %v to int", whyStr)
		return -1
	}

	floaty := math.Pow(float64(2.0), float64(z)) - float64(y)
	floaty--

	return int32(floaty)
}
