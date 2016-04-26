package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"randomImage/sources"
	"strings"
	"time"
)

var sourceSlice []sources.Source
var sourceMap map[string][]sources.Source
var blacklist []string
var hiddenTags map[string]struct{}

func listFiveRandoms(s sources.Source) {
	for i := 0; i < 5; i++ {
		image, source, err := s.GetRandomImage()
		if err != nil {
			log.Println("Error getting image:", err)
		}
		log.Println(source, ":", image)
	}
}

func loadSources(filename string) {
	sourceSlice = make([]sources.Source, 0, 5)
	sourceMap = make(map[string][]sources.Source)

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println("Failed to read file:", err)
		return
	}

	configs := make([]sources.TumblrConfig, 0, 5)
	err = json.Unmarshal(data, &configs)
	if err != nil {
		log.Println("Error in config:", err)
		return
	}
	for _, c := range configs {
		ts, err := sources.NewTumblrSource(c)
		if err != nil {
			log.Println("Failed to create Source from config:", err)
			continue
		}

		sourceSlice = append(sourceSlice, ts)
		for _, tag := range c.Tags {
			sourceMap[tag] = append(sourceMap[tag], ts)
		}
	}
}

func randomImage(tag string) (string, string, error) {
	var ok bool
	var err error
	var img string
	var source string
	var sourceIndex int
	var relevantSources []sources.Source

	if len(tag) == 0 {
		relevantSources = sourceSlice
	} else {
		if relevantSources, ok = sourceMap[tag]; !ok {
			return "", "", fmt.Errorf("No tag \"%s\" found", tag)
		}
	}

	sourceIndex = rand.Intn(len(relevantSources))
	for i := 0; i < 10; i++ {
		img, source, err = relevantSources[sourceIndex].GetRandomImage()
		if err == nil {
			if !blacklisted(img) {
				return img, source, nil
			} else {
				log.Println("randomImage: Got blacklisted. Trying again:", img)
				i-- //getting a blacklisted img doesn't count as a try
			}
		}
	}
	return "error", "", err
}

func blacklisted(img string) bool {
	for _, bi := range blacklist {
		if bi == img {
			return true
		}
	}

	//append img end roll out oldest blacklisted image
	blacklist = append(blacklist[1:], img)
	return false
}

func hidden(tag string) bool {
	_, hidden := hiddenTags[tag]
	return hidden

}

func parseTags(tagsString string) []string {
	if len(tagsString) == 0 {
		return nil
	}
	return strings.Split(tagsString, ",")
}

func jsonErrorMsg(msg interface{}) []byte {
	return []byte(fmt.Sprintf("{\"error\": \"%v\"}", msg))
}

func setupBlacklist(n int) {
	blacklist = make([]string, n)
}

func setupHiddenTags(secrets ...string) {
	hiddenTags = make(map[string]struct{})

	for _, secret := range secrets {
		hiddenTags[secret] = struct{}{}
	}
}

func main() {
	blacklistSize := 100
	servingAddress := ":12345"

	log.Println("Starting")
	rand.Seed(time.Now().UnixNano())
	setupBlacklist(blacklistSize)

	log.Println("Loading config")
	loadSources("tumblr.json")

	log.Println("Setting up handlers")
	http.HandleFunc("/sexy/", rootHandler)
	http.HandleFunc("/sexy/api/", apiHandler)
	http.HandleFunc("/sexy/tags/", tagHandler)
	http.HandleFunc("/sexy/config/", confHandler)
	http.HandleFunc("/sexy/config/add/", addConfHandler)

	log.Println("Serving")
	http.ListenAndServe(servingAddress, nil)
}
