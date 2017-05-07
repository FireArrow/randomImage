package sources

import (
	"fmt"
	"github.com/lestrrat/go-libxml2"
	"github.com/lestrrat/go-libxml2/xpath"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const urlTemplate = "%v/page/%d"

type TumblrSource struct {
	config      TumblrConfig
	imgPath     string
	sizePath    string
	sizePattern *regexp.Regexp
	url         string
	size        int64
}

type TumblrConfig struct {
	ImgPath     string   `json:"img_path"`
	SizePath    string   `json:"size_path"`
	SizePattern string   `json:"size_pattern"`
	Url         string   `json:"url"`
	Tags        []string `json:"tags"`
}

func validXpath(xp string) bool {
	e, err := xpath.NewExpression(xp)
	if err != nil {
		return false
	}
	e.Free()
	return true
}

func NewTumblrSource(config TumblrConfig) (*TumblrSource, error) {
	if !validXpath(config.ImgPath) {
		return nil, fmt.Errorf("Not valid xpath: %s", config.ImgPath)
	}
	if !validXpath(config.SizePath) {
		return nil, fmt.Errorf("Not valid xpath: %s", config.SizePath)
	}

	re, err := regexp.Compile(config.SizePattern)
	if err != nil {
		return nil, fmt.Errorf("Failed to compile regexp: %v", err)
	}

	ts := &TumblrSource{
		config:      config,
		imgPath:     config.ImgPath,
		sizePath:    config.SizePath,
		url:         strings.TrimSuffix(config.Url, "/"),
		sizePattern: re,
		size:        1,
	}

	err = ts.updateSize()
	if err != nil {
		return nil, err
	}
	return ts, nil
}

func (ts *TumblrSource) GetConfig() Config {
	return ts.config
}

func (ts *TumblrSource) GetTags() []string {
	return ts.config.Tags
}

func (ts *TumblrSource) updateSize() error {

	resp, err := http.Get(ts.url)
	if err != nil {
		return fmt.Errorf("Failed to update size: %v", err)
	}
	defer resp.Body.Close()

	doc, err := libxml2.ParseHTMLReader(resp.Body)
	if err != nil {
		return fmt.Errorf("Failed to parse HTML: %v", err)
	}
	defer doc.Free()

	nodes, err := doc.Find(ts.sizePath)
	if err != nil {
		return fmt.Errorf("Failed to apply xpath: %v", err)
	}
	defer nodes.Free()

	node := nodes.NodeList().First()
	if node == nil {
		return fmt.Errorf("Error when getting first node from nodes. Error in config?")
	}
	size, err := ts.parseSize(node.String())
	if err != nil {
		return err
	}
	ts.size = size

	return nil
}

func (ts *TumblrSource) Size() int64 {
	return ts.size
}

func (ts *TumblrSource) GetRandomImage() (string, string, error) {
	pageNumber := rand.Int63n(ts.Size()) + 1 //Page numbers are 1-indexed, and rand is [0-n)
	images, source, err := ts.ListPage(pageNumber)
	if err != nil {
		return "", "", fmt.Errorf("Failed to get image: %v", err)
	}

	imageNumber := rand.Intn(len(images))
	return images[imageNumber], source, nil
}

func (ts *TumblrSource) ListPage(pageNumber int64) ([]string, string, error) {
	retVal := make([]string, 0, 10)
	fullUrl := fmt.Sprintf(urlTemplate, ts.url, pageNumber)
	resp, err := http.Get(fullUrl)
	if err != nil {
		return retVal, fullUrl, err
	}
	defer resp.Body.Close()

	doc, err := libxml2.ParseHTMLReader(resp.Body)
	if err != nil {
		return retVal, fullUrl, fmt.Errorf("Failed to parse HTML: %v", err)
	}
	defer doc.Free()

	nodes, err := doc.Find(ts.imgPath)
	if err != nil {
		return retVal, fullUrl, fmt.Errorf("Failed to apply xpath: %v", err)
	}
	defer nodes.Free()

	it := nodes.NodeIter()
	for it.Next() {
		imgSrc := parseImgSrc(it.Node().String())
		retVal = append(retVal, imgSrc)
	}

	return retVal, fullUrl, nil
}

func parseImgSrc(nodeString string) string {
	parts := strings.Split(nodeString, `"`)
	if len(parts) < 3 {
		return ""
	}
	return parts[1]
}

func (ts *TumblrSource) parseSize(s string) (int64, error) {

	matches := ts.sizePattern.FindStringSubmatch(s)
	if matches == nil {
		return 1, fmt.Errorf("No match for size in \"%s\"", s)
	}
	if len(matches) < 2 {
		return 1, fmt.Errorf("Failed to match size from \"%s\"", s)
	}

	size, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return 1, fmt.Errorf("Failed to parse int from \"%s\"\n", matches[1])
	}

	return size, nil
}
