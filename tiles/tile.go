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
package tiles

import (
	"math"
	"strconv"
	"github.com/ruptivespatial/chopper/utils"
)

type tile struct {
	z    int
	x    int
	y    int
	lat  float64
	long float64
	Data []byte
}
/** Utility functions **/
func (t *tile) ZStr() (string){
	return strconv.Itoa(t.z)
}
func (t *tile) XStr() (string){
	return strconv.Itoa(t.x)
}
func (t *tile) YStr() (string){
	return strconv.Itoa(t.y)
}
func (t *tile) GetUrl() (string){
	return "/tiles/"+t.ZStr()+"/"+t.XStr()+"/"+t.YStr()+".pbf"
}

func NewTile(z int, x int, y int) (*tile){
	t := new(tile)
	t.z = z
	t.y = y
	t.x = x

	return t
}
func NewTileStr(z string, x string, y string) (*tile){
	var err error
	t := new(tile)
	t.z, err = strconv.Atoi(z)
	t.y, err = strconv.Atoi(y)
	t.x, err = strconv.Atoi(x)
	if(err != nil){
		utils.GetLogging().Error("error parsing value from string in tile coordinate")

	}
	return t
}


type Conversion interface {
	deg2num(t *tile) (x int, y int)
	num2deg(t *tile) (lat float64, long float64)
}

func (*tile) Deg2num(t *tile) (x int, y int) {
	x = int(math.Floor((t.long + 180.0) / 360.0 * (math.Exp2(float64(t.z)))))
	y = int(math.Floor((1.0 - math.Log(math.Tan(t.lat*math.Pi/180.0)+1.0/math.Cos(t.lat*math.Pi/180.0))/math.Pi) / 2.0 * (math.Exp2(float64(t.z)))))
	return
}

func (*tile) Num2deg(t *tile) (lat float64, long float64) {
	n := math.Pi - 2.0*math.Pi*float64(t.y)/math.Exp2(float64(t.z))
	lat = 180.0 / math.Pi * math.Atan(0.5*(math.Exp(n)-math.Exp(-n)))
	long = float64(t.x)/math.Exp2(float64(t.z))*360.0 - 180.0
	return lat, long
}
