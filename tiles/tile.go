/*Package tiles ...
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
	"github.com/ruptivespatial/chopper/utils"
	"strconv"
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
func (t *tile) ZStr() string {
	return strconv.Itoa(t.z)
}
func (t *tile) XStr() string {
	return strconv.Itoa(t.x)
}
func (t *tile) YStr() string {
	return strconv.Itoa(t.y)
}
func (t *tile) GetURL() string {
	return "/tiles/" + t.ZStr() + "/" + t.XStr() + "/" + t.YStr() + ".pbf"
}
//NewTile creates a tile with the given z/x/y
func NewTile(z int, x int, y int) *tile {
	t := new(tile)
	t.z = z
	t.y = y
	t.x = x

	return t
}
//NewTileStr creates a tile with the given z/x/y but in string format
func NewTileStr(z string, x string, y string) *tile {
	var err error
	t := new(tile)
	t.z, err = strconv.Atoi(z)
	t.y, err = strconv.Atoi(y)
	t.x, err = strconv.Atoi(x)
	if err != nil {
		utils.GetLogging().Error("error parsing value from string in tile coordinate")

	}
	return t
}

