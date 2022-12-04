package sitemap

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/RossoDiablo/html_link_parser/link"
)

func get(URL string) ([]string, error) {
	resp, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	reqURL := resp.Request.URL
	baseURL := &url.URL{
		Scheme: reqURL.Scheme,
		Host:   reqURL.Host,
	}
	base := baseURL.String()
	pages, err := hrefs(resp.Body, base)
	if err != nil {
		return nil, err
	}
	return filter(pages, withPrefix(base)), nil
}

func filter(links []string, filterOpts ...func(string) bool) []string {
	filtered := make([]string, 0)
	for _, link := range links {
		for _, opt := range filterOpts {
			if opt(link) {
				filtered = append(filtered, link)
			}
		}
	}
	return filtered
}

func withPrefix(prefix string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, prefix)
	}
}

func hrefs(r io.Reader, base string) ([]string, error) {
	links, err := link.Parse(r)
	if err != nil {
		return nil, err
	}
	hrefs := make([]string, 0)
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			hrefs = append(hrefs, base+l.Href)
		case strings.HasPrefix(l.Href, "http"):
			hrefs = append(hrefs, l.Href)
		}
	}
	return hrefs, nil
}

func Create(URL string, depth int) ([]string, error) {
	visited := make(map[string]struct{})
	q := make(map[string]struct{})
	nq := map[string]struct{}{URL: {}}
	countErrored := 0
	for i := 0; i <= depth; i++ {
		q, nq = nq, make(map[string]struct{})
		if len(q) == 0 {
			break
		}
		for link := range q {
			if _, ok := visited[link]; ok {
				continue
			}
			visited[link] = struct{}{}
			links, err := get(link)
			if err != nil {
				countErrored++
				continue
			}
			for _, l := range links {
				nq[l] = struct{}{}
			}
		}
	}
	sitemap := make([]string, 0, len(visited))
	for link := range visited {
		sitemap = append(sitemap, link)
	}
	if countErrored != 0 {
		return sitemap, errors.New("Some links were failed to get!")
	}
	return sitemap, nil
}
