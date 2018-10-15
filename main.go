package main

import (
	"bytes"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	url  = flag.String("url", "", "url")
	path = flag.String("path", "", "file base path")
)

func main() {
	flag.Parse()
	if *url == "" || *path == "" {
		log.Fatal("need url and path")
	}
	resp, err := http.Get(*url)
	if err != nil {
		log.Fatal(err.Error())
	}
	rawBody, _ := ioutil.ReadAll(resp.Body)
	reg, _ := regexp.Compile(`(http://|https://).+?\.jpg`)
	imgs := reg.FindAllString(string(rawBody), -1)
	for _, i := range imgs {
		now := time.Now().UnixNano()
		path := filepath.Join(*path, strconv.FormatInt(now, 10)+".jpg")
		saveImg(i, path)
		log.Println("%s finished===", i)
	}
	log.Println("=======all finished=========")
}

func saveImg(imgUrl string, path string) (err error) {
	resp, err := http.Get(imgUrl)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	out, err := os.Create(path)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	_, err = io.Copy(out, bytes.NewReader(body))
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	return
}

func trimHttpPrefix(url string) string {
	return strings.TrimSuffix(strings.TrimPrefix(url, "http://"), "https://")
}
