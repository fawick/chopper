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
	"database/sql"
	"github.com/allegro/bigcache"
	"github.com/ruptivespatial/chopper/utils"
	"os"
	"strconv"
	"time"
)
//TileManager manages the retrieval and caching of tiles
type TileManager struct {
	cache       *bigcache.BigCache
	dbs         []*sql.DB
	prepStmts   []*sql.Stmt
	fullyCached bool
}
//GetTile returns a tile based on a z/x/y request
func (tm *TileManager) GetTile(z string, x string, y string) *tile {

	tile := NewTileStr(z, x, y)
	var tiledata []byte
	key := buildKey(z, x, y)
	tiledata, err := tm.cache.Get(key)

	//if tile is empty we can check the DBs unless we know everything is already loaded in the cache
	if tile == nil || err != nil {

		if !tm.fullyCached {
			//iterate over databases and look for a tile
			for _, stmt := range tm.prepStmts {
				row := stmt.QueryRow(z, x, y)
				row.Scan(&tiledata)
				if tiledata != nil {
					break
				}
			}
		}
	}
	tile.Data = tiledata
	return tile
}
//NewTileManager creates an instance of a tile manager based on a list of mbtile files
func NewTileManager(mbtilePath []string, useCache bool) *TileManager {

	tm := TileManager{}
	utils.GetLogging().Info("Initializing tile manager...")

	//initialize cache....100mb by default
	config := bigcache.Config{Shards: 1024, Verbose: false, HardMaxCacheSize: utils.GetSettings().GetCacheSizeMB() * 1000}
	cache, initErr := bigcache.NewBigCache(config)
	tm.cache = cache

	if initErr != nil {
		utils.GetLogging().Error("Error creating cache!")
	}
	utils.GetLogging().Debug("Cache initialized")

	var conns = make([]*sql.DB, 0)
	var preps = make([]*sql.Stmt, 0)

	for _, connStr := range mbtilePath {

		fi, err := os.Stat(connStr)
		if err != nil {
			utils.GetLogging().Error("Database %v does not exist...exiting", connStr)
			continue
		}

		//// Open database file
		db, err := sql.Open("sqlite3", connStr)
		conns = append(conns, db)
		if err != nil {
			utils.GetLogging().Error("Error opening database!")
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
		prepStmt, err := db.Prepare("SELECT tile_data as tile FROM tiles where zoom_level=? AND tile_column=? AND tile_row=?")
		preps = append(preps, prepStmt)

	}
	tm.prepStmts = preps
	tm.dbs = conns

	return &tm

}

func buildKey(z string, x string, y string) string {
	return y + "_" + x + "_" + z
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
		var key = buildKey(strconv.Itoa(tileRow), strconv.Itoa(tileColumn), strconv.Itoa(zoomLevel))
		//log.Println(key)
		cache.Set(key, tile)
	}
	elapsed := time.Since(start)

	utils.GetLogging().Info("Finished loading level %v (%v tiles) in %s", zoom, tilecount, elapsed)

}
