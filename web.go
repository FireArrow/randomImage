package main

import (
	"encoding/json"
	"fmt"
	"github.com/FireArrow/randomImage/sources"
	"log"
	"net/http"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	tag := r.FormValue("tag")
	img, source, err := randomImage(tag)

	log.Printf("HTTP: Serving %v: %s", r.RemoteAddr, img)
	if err != nil {
		log.Println("HTTP: Error:", err)
		if err.Error() == "error" {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(rootPage(img, source))
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	tag := r.FormValue("tag")
	img, _, err := randomImage(tag)

	log.Printf("API:  Serving %v: %s", r.RemoteAddr, img)
	if err != nil {
		log.Println("API:  Error:", err)
		if err.Error() == "error" {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
		w.Write([]byte(err.Error()))
		return
	}
	http.Redirect(w, r, img, http.StatusTemporaryRedirect)
}

func tagHandler(w http.ResponseWriter, r *http.Request) {
	tags := make([]string, 0, len(sourceMap))
	for tag, _ := range sourceMap {
		if !hidden(tag) {
			tags = append(tags, tag)
		}
	}

	log.Printf("TAG:  Serving %v", r.RemoteAddr)
	jsonTags, err := json.Marshal(tags)
	if err != nil {
		log.Println("Error creating json tag-list:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonErrorMsg("Failed to compose tag list"))
		return
	}

	w.Write(jsonTags)
}

func confHandler(w http.ResponseWriter, r *http.Request) {
	configs := make([]sources.Config, 0, len(sourceSlice))
	for _, source := range sourceSlice {
		config := source.GetConfig()
		configs = append(configs, config)
	}

	jsonConfig, err := json.MarshalIndent(configs, "", "  ")
	if err != nil {
		log.Println("Error marshalling config list:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonErrorMsg("Failed to create config list"))
		return
	}

	w.Write(jsonConfig)
}

func addConfHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:

	case http.MethodGet:

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(jsonErrorMsg(fmt.Sprintf("Method %s not supported", r.Method)))
	}
}
