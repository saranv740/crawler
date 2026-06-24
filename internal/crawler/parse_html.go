package crawler

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type PageData struct {
	URL            string   `json:"url"`
	Heading        string   `json:"heading"`
	FirstParagraph string   `json:"first_paragraph"`
	OutgoingLinks  []string `json:"outgoing_links"`
	ImageURLs      []string `json:"image_urls"`
}

func getHeadingFromHTML(doc *goquery.Document) string {
	node := doc.Find("h1")
	if node.Length() == 0 {
		node = doc.Find("h2")
	}

	return node.Text()
}

func getParagraphFromHTML(doc *goquery.Document) string {
	node := doc.Find("main")
	if node.Length() == 0 {
		node = doc.Find("p")
	}

	return node.Text()
}

func getLinksFromHTML(doc *goquery.Document, baseURL *url.URL) ([]string, error) {
	result := make([]string, 0)
	doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		href := s.AttrOr("href", "")

		if strings.HasPrefix(href, "/") {
			href = fmt.Sprintf("%s://%s%s", baseURL.Scheme, baseURL.Hostname(), href)
		}

		if href != "" {
			result = append(result, href)
		}
	})

	return result, nil
}

func getLinksFromImages(doc *goquery.Document, baseURL *url.URL) ([]string, error) {
	result := make([]string, 0)
	doc.Find("img[src]").Each(func(_ int, s *goquery.Selection) {
		src := s.AttrOr("src", "")

		if strings.HasPrefix(src, "/") {
			src = fmt.Sprintf("%s%s", baseURL, src)
		}

		if src != "" {
			result = append(result, src)
		}
	})

	return result, nil
}

func extractPageData(rawHTML string, currentURL string, baseURL *url.URL) (PageData, []error) {
	reader := strings.NewReader(rawHTML)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return PageData{}, []error{err}
	}

	errs := make([]error, 0)
	hrefs, err := getLinksFromHTML(doc, baseURL)

	if err != nil {
		errs = append(errs, err)
	}

	imageURLs, err := getLinksFromImages(doc, baseURL)

	if err != nil {
		errs = append(errs, err)
	}

	return PageData{
		URL:            currentURL,
		Heading:        getHeadingFromHTML(doc),
		FirstParagraph: getParagraphFromHTML(doc),
		OutgoingLinks:  hrefs,
		ImageURLs:      imageURLs,
	}, errs
}

func getHTML(rawURL string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	req.Header.Set("User-Agent", "BootCrawler/1.0")
	if err != nil {
		return "", err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error when requesting ", rawURL, err)
		return "", err
	}
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	if !strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		return "", fmt.Errorf("response is not html")
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	return string(body), nil
}
