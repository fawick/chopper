package main

import (
	"log"
	"net/http"
	_ "github.com/mattn/go-sqlite3"
	"github.com/julienschmidt/httprouter"
	"strconv"
	"github.com/tingold/squirrelchopper/tiles"
	"flag"
	"github.com/tingold/squirrelchopper/handler"
	"github.com/tingold/squirrelchopper/utils"

)

var tm *tiles.TileManager
var th *handler.Tilehandler

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
	//initialize the tile manager with one or mone databases
	tm = tiles.NewTileManager(settings.GetDBs(), true)
	th = new(handler.Tilehandler)
	th.Manager = *tm

	//create the HTTP Router -- the tiles go to the tilehandler which uses the tile manager to access the DB
	//and potentially cache
	router := httprouter.New()
	router.GET("/tiles/:z/:x/:y", th.Handle)


	//any non tile request will default to serving files
	//files are NOT actually files but stored in GO code using
	// https://github.com/elazarl/go-bindata-assetfs
	fs := http.FileServer(assetFS())
	router.NotFound = fs

	//create server
	srv := &http.Server{
		Addr:    ":"+strconv.Itoa(settings.GetPort()),
		Handler: router,
	}

	logger.Info("Starting server on port %v",settings.GetPort())
	if(settings.GetSsl()) {
		logger.Info("Using certificate %v and key %v", settings.GetSslCert(), settings.GetSslKey())
		//this creates the SSL server which we really need to use server push
		error := srv.ListenAndServeTLS(settings.GetSslCert(),settings.GetSslKey())
		if(error != nil){
			log.Fatalf("Failed to start server: %v",error)
		}
	} else {
		error := srv.ListenAndServe()
		if(error != nil){
			log.Fatalf("Failed to start server: %v",error)
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

