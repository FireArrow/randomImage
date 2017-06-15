package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/FireArrow/randomImage/sources"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"strings"
	"time"
)

var sourceSlice []sources.Source
var sourceMap map[string][]sources.Source
var blacklist []string
var hiddenTags map[string]struct{}
var servingAddress = flag.String("addr", ":12345", "Listening address")

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

func getMainTags() []string {
	tagMap := make(map[string]struct{})
	for _, source := range sourceSlice {
		tags := source.GetTags()
		if tags != nil {
			sort.Strings(tags)
			tagMap[tags[0]] = struct{}{}
		}
	}

	mainTags := make([]string, len(tagMap))
	i := 0
	for tag, _ := range tagMap {
		mainTags[i] = tag
		i++
	}
	return mainTags
}

func randomImage(tag string) (*image, error) {
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
			return nil, fmt.Errorf("No tag \"%s\" found", tag)
		}
	}

	sourceIndex = rand.Intn(len(relevantSources))
	for i := 0; i < 10; i++ {
		img, source, err = relevantSources[sourceIndex].GetRandomImage()
		if err == nil {
			if !blacklisted(img) {
				return &image{
					Url:    img,
					Source: source,
				}, nil
			} else {
				log.Println("randomImage: Got blacklisted. Trying again:", img)
				i-- //getting a blacklisted img doesn't count as a try
			}
		}
	}
	return nil, err
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
	blacklistSize := flag.Int("blsize", 800, "Size of the backlist")
	rawBasePath := flag.String("basepath", "", "Base path, e.g. \"/test\" would cause server to listen to host:post/test/")
	flag.Parse()

	log.Println("Starting")
	rand.Seed(time.Now().UnixNano())
	setupBlacklist(*blacklistSize)

	basePath := fmt.Sprintf("/%v", strings.Trim(*rawBasePath, "/ "))
	if len(basePath) == 1 {
		basePath = ""
	}

	log.Println("Loading config")
	loadSources("tumblr.json")

	log.Println("Setting up handlers, base", basePath+"/")
	http.Handle(basePath+"/", http.StripPrefix(basePath+"/", http.FileServer(http.Dir("static"))))
	http.HandleFunc(basePath+"/api/", apiHandler)
	http.HandleFunc(basePath+"/tags/", tagHandler)
	http.HandleFunc(basePath+"/config/", confHandler)
	http.HandleFunc(basePath+"/config/add/", addConfHandler)

	log.Println("Serving on", *servingAddress)
	err := http.ListenAndServe(*servingAddress, nil)
	if err != nil {
		log.Println("Error while trying to serve:", err)
	}
	log.Println("Exiting")
}
