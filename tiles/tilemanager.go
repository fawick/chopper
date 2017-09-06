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
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/allegro/bigcache"
	"github.com/ruptivespatial/chopper/utils"
	"github.com/tingold/gophertile/gophertile"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

//TileManager manages the retrieval and caching of tiles
type TileManager struct {
	cache       *bigcache.BigCache
	Metadatas   []DBMetadata
	fullyCached bool
}

//DBMetadata holds all info and object for the various mbtile files
type DBMetadata struct {
	Id        string
	Fields    map[string]string
	LayerInfo interface{}
	Conn      *sql.DB   `json:"-"`
	Prep      *sql.Stmt `json:"-"`
}

//GetTile returns a tile based on a z/x/y request
func (tm *TileManager) GetTile(z int, x int, y int) (*gophertile.Tile, []byte) {
	//func (tm *TileManager) GetTile(z string, x string, y string) *T"ile {

	tile := gophertile.Tile{x, y, z}
	var tiledata []byte
	key := buildKey(z, x, y)
	tiledata, err := tm.cache.Get(key)

	//if tile is empty we can check the DBs unless we know everything is already loaded in the cache
	if tiledata == nil || err != nil {

		if !tm.fullyCached {
			//iterate over databases and look for a tile
			for _, metadata := range tm.Metadatas {
				row := metadata.Prep.QueryRow(z, x, y)
				row.Scan(&tiledata)
				if tiledata != nil {
					break
				}
			}
		}
	}

	return &tile, tiledata
}

//NewTileManager creates an instance of a tile manager based on a list of mbtile files
func NewTileManager(mbtilePath []string, useCache bool) *TileManager {

	tm := TileManager{}
	utils.GetLogging().Info("Initializing tile manager...")
	tm.Metadatas = make([]DBMetadata, 0)
	//initialize cache....100mb by default
	config := bigcache.Config{Shards: 1024, Verbose: false, HardMaxCacheSize: utils.GetSettings().GetCacheSizeMB() * 1000}
	cache, initErr := bigcache.NewBigCache(config)
	tm.cache = cache

	if initErr != nil {
		utils.GetLogging().Error("Error creating cache!")
	}
	utils.GetLogging().Debug("Cache initialized")

	for _, connStr := range mbtilePath {

		fi, err := os.Stat(connStr)
		if err != nil {
			utils.GetLogging().Error("Database %v does not exist...exiting", connStr)
			continue
		}

		// Open database file
		db, err := sql.Open("sqlite3", connStr)
		if err != nil {
			utils.GetLogging().Error("Error opening database!")
			continue
		}
		//initialize database info
		dbMetadata := DBMetadata{}
		_, dbMetadata.Id = filepath.Split(connStr)
		dbMetadata.Fields = make(map[string]string)
		dbMetadata.Conn = db

		//load metadata
		metadataRows, err := db.Query("Select name, value FROM metadata")
		for metadataRows.Next() {
			var name string
			var val string
			metadataRows.Scan(&name, &val)
			if name == "json" {
				json.Unmarshal(bytes.NewBufferString(val).Bytes(), &dbMetadata.LayerInfo)
			} else {
				dbMetadata.Fields[name] = val
			}
			utils.GetLogging().Warn("key %v val: %v", name, val)
		}

		//see if we can fit the whole thing...
		if fi.Size() < int64(utils.GetSettings().GetCacheSizeMB())*1000000 {
			utils.GetLogging().Info("Database %v is %v MB....going to try to fit it into RAM", connStr, fi.Size()/1000000)

			for i := 0; i < 15; i++ {
				loadTileLevelIntoCache(i, db, cache)
			}
		} else {
			utils.GetLogging().Info("Database %v is %v MB - too big to cache", connStr, fi.Size()/1000000)
		}

		var count int
		rows := db.QueryRow("SELECT COUNT(*) as count from tiles")
		rows.Scan(&count)
		utils.GetLogging().Info("Found %v tiles in db", count)

		////prepare statement
		prepStmt, _ := db.Prepare("SELECT tile_data as tile FROM tiles where zoom_level=? AND tile_column=? AND tile_row=?")
		dbMetadata.Prep = prepStmt
		tm.Metadatas = append(tm.Metadatas, dbMetadata)

	}
	return &tm

}

func buildKey(z int, x int, y int) string {
	buf := bytes.NewBufferString(strconv.Itoa(z))
	buf.WriteString("_")
	buf.WriteString(strconv.Itoa(x))
	buf.WriteString("_")
	buf.WriteString(strconv.Itoa(z))
	return buf.String()
}

func loadTileLevelIntoCache(zoom int, database *sql.DB, cache *bigcache.BigCache) {

	start := time.Now()

	rows, err := database.Query("SELECT zoom_level, tile_column, tile_row, tile_data FROM tiles where zoom_level=" + strconv.Itoa(zoom))

	defer rows.Close()
	var tilecount int
	for rows.Next() {
		tilecount++
		var zoomLevel int
		var tileColumn int
		var tileRow int
		var tile []byte
		err = rows.Scan(&zoomLevel, &tileColumn, &tileRow, &tile)
		if err != nil {
			utils.GetLogging().Error("Error loading row: %v", err)
			continue
		}
		var key = buildKey(tileRow, tileColumn, zoomLevel)

		cache.Set(key, tile)
	}
	elapsed := time.Since(start)

	utils.GetLogging().Info("Finished loading level %v (%v tiles) in %s", zoom, tilecount, elapsed)

}
