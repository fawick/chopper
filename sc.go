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
package main

import (
	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
	"github.com/namsral/flag"
	"github.com/ruptivespatial/chopper/handler"
	"github.com/ruptivespatial/chopper/tiles"
	"github.com/ruptivespatial/chopper/utils"
	"log"
	"net/http"
	"strconv"
)

var tm *tiles.TileManager
var th *handler.Tilehandler
var infoHandler *handler.InfoHandler

func main() {
	//parse out settings
	settings := utils.GetSettings()

	//print help message if needed
	if settings.GetHelp() {
		flag.PrintDefaults()
		return
	}
	logger := utils.GetLogging()

	logger.Debug("Setting up tile manager with data from %v", settings.GetDBs())
	//initialize the tile manager with one or more databases
	th = new(handler.Tilehandler)
	tm = tiles.NewTileManager(settings.GetDBs(), true)
	th.Manager = *tm
	infoHandler = &handler.InfoHandler{ tm}

	//create the HTTP Router -- the tiles go to the tilehandler which uses the tile manager to access the DB
	//and potentially cache
	router := httprouter.New()
	router.GET("/tiles/:z/:x/:y", th.Handle)

	//this is a work around to rewrite values in the json file as the Mapbox style format doesn't support relative paths
	ph := handler.NewProxyHostHandler(assetFS())
	router.GET("/style/osm-liberty/osm-liberty.json", ph.Handle)
	router.GET("/style/osm-liberty.json", ph.Handle)
	router.GET("/style/osm-liberty@2x.json", ph.Handle)
	router.GET("/style/osm-liberty/osm-liberty@2x.json", ph.Handle)

	router.GET("/info", infoHandler.Handle)

	//any non tile request will default to serving files
	//files are NOT actually files but stored in GO code using
	// https://github.com/elazarl/go-bindata-assetfs
	fs, ok := http.FileServer(assetFS()).(http.HandlerFunc)
	if ok{
		router.NotFound = fs
	}
	

	//create server
	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(settings.GetPort()),
		Handler: router,
	}

	logger.Info("Starting server on port %v", settings.GetPort())
	if settings.GetSsl() {
		logger.Info("Using certificate %v and key %v", settings.GetSslCert(), settings.GetSslKey())
		//this creates the SSL server which we really need to use server push
		error := srv.ListenAndServeTLS(settings.GetSslCert(), settings.GetSslKey())
		if error != nil {
			log.Fatalf("Failed to start server: %v", error)
		}
	} else {
		error := srv.ListenAndServe()
		if error != nil {
			logger.Fatal("Failed to start server: %v", error)
		}
	}

	logger.Info("Exiting")

}

func corsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
