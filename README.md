# A web crawler in go
This is a simple web crawler made in golang (following [boot.dev](https://www.boot.dev)'s make a [web scraper in go](https://www.boot.dev/courses/build-web-scraper-golang) course). The crawler concurrently fetches web pages and generates a report in specific format as below.

```json
[
	{
		"url": "https://learnwebscraping.dev/",
		"heading": "Web Scraping Practice Sandbox",
		"first_paragraph": "First paragraph",
		"outgoing_links": [
			"https://learnwebscraping.dev/",
			"https://learnwebscraping.dev/practice/",
		],
		"image_urls": []
	},
	{
		"url": "https://learnwebscraping.dev/practice/",
		"heading": "Scraping Practice",
		"first_paragraph": "First paragraph",
		"outgoing_links": [
			"https://learnwebscraping.dev/terms-of-service/",
			"https://www.boot.dev"
		],
		"image_urls": [
            "https://learnwebscraping.dev/practice/ecommerce/logo.png"
        ]
	}
]
```

## Build
```bash
# windows
go build -o build/crawler.exe
# linux
go build -o build/crawler
```

## Usage
```bash
# windows
.\build\crawler.exe -url=https://learnwebscraping.dev/practice/ecommerce/ -conn=5 -pages=100 -output="output/report.json"
#linux
./build/crawler -url=https://learnwebscraping.dev/practice/ecommerce/ -conn=5 -pages=100 -output="output/report.json"
```

## Flags
```text
  -conn int
        max number of concurrent workers (default 4)
  -output string
        the file name to write the report
  -pages int
        max pages to collect (default 10)
  -url string
        the url to crawl
```

## TODO
- [ ] The current version spawns one goroutine per url and controls concurrency using semaphores, another approach would be keeping a jobs queue and fixed number of workers to keep memory flat.
- [ ] Organize files into respective packages.