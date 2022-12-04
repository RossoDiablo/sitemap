package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/RossoDiablo/sitemap/internal/encode"
	"github.com/RossoDiablo/sitemap/internal/sitemap"
)

const sitemapFile = "sitemap.xml"
var (
	URL := flag.String("url", "https://gophercises.com", "URL of a site")
	Depth := flag.Int("depth", 1, "maximum depth of links")
)

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func main() {
	fmt.Println("Hello!")

	flag.Parse()

	sitemap, err := sitemap.Create(*URL, *Depth)
	if err != nil {
		exit("Error creating sitemap!")
	}
	err = encode.XML(sitemapFile, sitemap)
	if err != nil {
		exit("Error encoding to XML!")
	}
	fmt.Println("Done successfully! Check sitemap.xml")
}
