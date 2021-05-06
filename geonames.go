package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

var (
	language string
	output   *os.File
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	regions := [6]string{"EU", "NA", "SA", "OC", "AS", "AF"}

	flag.StringVar(&language, "lang", "en", "ISO 639 language code")
	flag.Parse()
	arg := flag.Args()
	outputDir := arg[0] // ex: pt-BR

	// does dir exist?
	_, err := os.ReadDir(outputDir)
	if err != nil {
		fmt.Println("error: does not exist")
		err = os.Mkdir(outputDir, 0744)
		check(err)
	}

	for _, code := range regions {

		fp := path.Join(outputDir, code+".json")
		output, err = os.Create(fp)
		check(err)
		url := "http://api.geonames.org/countryInfoJSON?continentCode=" + code + "&lang=" + language + "&username=raconteur"
		resp, err := http.Get(url)
		if err != nil {
			log.Fatalln(err)
		}

		//We Read the response body on the line below.
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		//Convert the body to type string
		sb := string(body)

		_, err = output.WriteString(sb)
		check(err)
		fmt.Printf("Output success: %s\n", fp)
	}

}
