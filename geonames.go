package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

var (
	language string
	output   *os.File
	wg       sync.WaitGroup
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getJSON(code string, outputDir string, wg *sync.WaitGroup) {

	defer wg.Done()

	fp := path.Join(outputDir, code+".json")
	output, err := os.Create(fp)
	check(err)
	// url := "http://api.geonames.org/countryInfoJSON?continentCode=" + code + "&lang=" + language + "&username=raconteur"
	url := "https://jsonplaceholder.typicode.com/posts/2"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = io.Copy(output, resp.Body)

	defer resp.Body.Close()

}

func main() {
	start := time.Now()

	regions := [6]string{"EU", "NA", "SA", "OC", "AS", "AF"}

	flag.StringVar(&language, "lang", "en", "ISO 639 language code")
	flag.Parse()
	arg := flag.Args()
	outputDir := arg[0] // ex: pt-BR

	if err := os.Mkdir(outputDir, 0744); err != nil && !os.IsExist(err) {
		panic(err)
	}

	for _, code := range regions {
		wg.Add(1)
		go getJSON(code, outputDir, &wg)

	}

	wg.Wait()

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("%v", elapsed)

}
