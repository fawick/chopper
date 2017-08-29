package handler

import (
	"net/http"
	"github.com/ruptivespatial/chopper/tiles"
	"github.com/julienschmidt/httprouter"
	"encoding/json"
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