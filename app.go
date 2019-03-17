package main

import (
	. "./config"
	. "./dao"
	. "./models"
	"encoding/json"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)

var config = Config{}
var dao = TracksDAO{}

func AllTracksEndPoint(w http.ResponseWriter, r *http.Request) {
	tracks, err := dao.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, tracks)
}

func FindTrackEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	movie, err := dao.FindById(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Track ID")
		return
	}
	respondWithJson(w, http.StatusOK, movie)
}

func CreateTrackEndPoint(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	var track Track
	if err := json.NewDecoder(r.Body).Decode(&track); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload in: CreateTrackEndPoint")
		return
	}
	track.ID = bson.NewObjectId()
	if err := dao.Insert(track); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, track)
}

func UpdateTrackEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var track Track
	if err := json.NewDecoder(r.Body).Decode(&track); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := dao.Update(track); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func DeleteTrackEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var track Track
	if err := json.NewDecoder(r.Body).Decode(&track); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := dao.Delete(track); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func init() {
	config.Read()

	dao.Server = config.Server
	dao.Database = config.Database
	dao.Connect()
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/tracks", AllTracksEndPoint).Methods("GET")
	r.HandleFunc("/tracks", CreateTrackEndPoint).Methods("POST")
	r.HandleFunc("/tracks", UpdateTrackEndPoint).Methods("PUT")
	r.HandleFunc("/tracks", DeleteTrackEndPoint).Methods("DELETE")
	r.HandleFunc("/tracks/{id}", FindTrackEndpoint).Methods("GET")

	/*c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:[]string{"X-Requested-With"},
		AllowedMethods: []string{"GET", "HEAD", "POST", "PUT", "OPTIONS"},
	})*/

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	log.Fatal(http.ListenAndServe(":8081", handlers.CORS(originsOk, headersOk, methodsOk)(r)))


}
