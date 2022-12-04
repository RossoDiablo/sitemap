package encode

import (
	"encoding/xml"
	"os"
)

type loc struct {
	URL string `xml:"loc"`
}

type sitemapXML struct {
	Xmlns string `xml:"xmlns,attr"`
	URLs  []loc  `xml:"url"`
}

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

func defaultSitemapXML() sitemapXML {
	return sitemapXML{Xmlns: xmlns, URLs: nil}
}

func toSitemapXML(data []string) sitemapXML {
	toXML := defaultSitemapXML()
	for _, page := range data {
		toXML.URLs = append(toXML.URLs, loc{URL: page})
	}
	return toXML
}

func XML(filename string, data []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	encoder := xml.NewEncoder(file)
	encoder.Indent("", "   ")

	toXML := toSitemapXML(data)
	err = encoder.Encode(toXML)
	if err != nil {
		return err
	}
	return nil
}
