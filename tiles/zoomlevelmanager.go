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
	"sync"
	"math"
	"github.com/tingold/squirrelchopper/utils"
)


var instance *zoomLevelManager
var once sync.Once

type zoomLevelManager struct{

	level []*zoomLevel

}
type zoomLevel struct{

	maxY int
	maxX int
	numTiles int
	zoom int
}
func (zm *zoomLevel) Zoom() (int){
	return zm.zoom
}
func (zm *zoomLevel) MaxX() (int){
	return zm.maxX
}
func (zm *zoomLevel) MaxY() (int){
	return zm.maxY
}
func (zm *zoomLevelManager) GetLevel(z int) (*zoomLevel){
	return zm.level[z]
}

//TODO this whole thing is a little redundant...
func GetZoomLevelManager() *zoomLevelManager {
	once.Do(func() {
		instance = &zoomLevelManager{}
		instance.level = make([]*zoomLevel,18)
		for i := 0; i < 18; i++{
			square := math.Pow(float64(2),float64(i))
			sqInt := int(square)
			zl := new(zoomLevel)
			zl.zoom = i
			zl.numTiles = sqInt * sqInt
			zl.maxX = sqInt
			zl.maxY =  sqInt
			instance.level[i] = zl
		}
	})
	return instance
}
func (zm *zoomLevelManager) ValidTile(t *tile) (bool) {
	zl := zm.GetLevel(t.z)
	//utils.GetLogging().Debug("Zoom level has max x and y of %v %v",zl.MaxX(), zl.MaxY())
	//utils.GetLogging().Debug("Tile has x and y of %v %v",t.x, t.y)
	return t.y <= zl.MaxY() && t.y <= zl.MaxY()
}

func (zm *zoomLevelManager) GetAdjacentTiles(t *tile) ([]*tile) {

	tiles := make([]*tile,8)
	for x := -1; x < 1; x++{
		for y := -1; y < 1; y++ {
			if(x == 0 && y == 0){continue}
			nt := NewTile(t.z,t.x+x, t.y+y)
			utils.GetLogging().Debug("Made new tile")
			if zm.ValidTile(nt){
				utils.GetLogging().Debug("Tile is valid")
				tiles =	append(tiles,nt)
				utils.GetLogging().Debug("Found valid adjacent tile")
			}
		}

	}
	return tiles
}
