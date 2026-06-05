package main

import (
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestGetHeadingFromHTML(t *testing.T) {
	t.Run("it extracts H1", func(t *testing.T) {
		input := `
<html>
  <body>
    <h1>Welcome to Boot.dev</h1>
    <main>
      <p>Learn to code by building real projects.</p>
      <p>This is the second paragraph.</p>
    </main>
  </body>
</html>
		`
		expected := "Welcome to Boot.dev"

		out := getHeadingFromHTML(input)
		if out != expected {
			t.Fatalf("expected %q got %q", expected, out)
		}
	})

	t.Run("extract H2 as fallback", func(t *testing.T) {
		input := `
<html>
  <body>
    <h2>H2 is awesome</h2>
    <main>
      <p>Learn to code by building real projects.</p>
      <p>This is the second paragraph.</p>
    </main>
  </body>
</html>
		`
		expected := "H2 is awesome"

		out := getHeadingFromHTML(input)
		if out != expected {
			t.Fatalf("expected %q got %q", expected, out)
		}
	})

	t.Run("return empty if no heading", func(t *testing.T) {
		input := `
<html>
  <body>
    <main>
    </main>
  </body>
</html>
		`
		expected := ""

		out := getHeadingFromHTML(input)
		if out != expected {
			t.Fatalf("expected %q got %q", expected, out)
		}
	})
}

func TestGetParagraphFromHTML(t *testing.T) {
	t.Run("it extracts main", func(t *testing.T) {
		input := `
<html>
  <body>
    <h1>Welcome to Boot.dev</h1>
    <main>
		Learn to code by building real projects.
    </main>
  </body>
</html>
		`
		expected := "Learn to code by building real projects."
		got := getParagraphFromHTML(input)

		if !strings.Contains(got, expected) {
			t.Fatalf("expected %q got %q", expected, got)
		}
	})
	t.Run("it extracts p as fallback", func(t *testing.T) {
		input := `
<html>
  <body>
    <h1>Welcome to Boot.dev</h1>
    <p>
		Learn to code by building real projects.
    </p>
  </body>
</html>
		`
		expected := "Learn to code by building real projects."
		got := getParagraphFromHTML(input)

		if !strings.Contains(got, expected) {
			t.Fatalf("expected %q got %q", expected, got)
		}
	})
}

func TestExtractLinksFromHTML(t *testing.T) {
	baseURL, err := url.Parse("https://crawler-test.com")
	if err != nil {
		t.Fatalf("invalid base url")
	}

	t.Run("extract links from anchors", func(t *testing.T) {
		input := `
<html>
  <body>
    <a href="https://crawler-test.com">Go to Boot.dev</a>
    <a href="/about">About</a>
    <a href="">Nothing</a>
  </body>
</html>
	`
		expected := []string{
			"https://crawler-test.com",
			"https://crawler-test.com/about",
		}

		got, err := getLinksFromHTML(input, baseURL)
		if err != nil {
			t.Fatalf("unexpected error %v", err)
		}

		if !reflect.DeepEqual(got, expected) {
			t.Fatalf("expected %v got %v", expected, got)
		}
	})

	t.Run("extract links from images", func(t *testing.T) {
		input := `
<html>
  <body>
	<img src="/logo.png" alt="Logo">
	<img src="" alt="Logo">
  </body>
</html>
	`
		expected := []string{
			"https://crawler-test.com/logo.png",
		}

		got, err := getLinksFromImages(input, baseURL)
		if err != nil {
			t.Fatalf("unexpected error %v", err)
		}

		if !reflect.DeepEqual(got, expected) {
			t.Fatalf("expected %v got %v", expected, got)
		}
	})
}

func TestExtractPageData(t *testing.T) {
	baseURL, err := url.Parse("https://crawler-test.com")
	if err != nil {
		t.Fatalf("invalid base url")
	}

	input := `
<html>
	<body>
        <h1>Test Title</h1>
        <p>This is the first paragraph.</p>
        <a href="/link1">Link 1</a>
        <img src="/image1.jpg" alt="Image 1">
    </body>
</html>`

	got, errs := extractPageData(input, baseURL)
	if len(errs) != 0 {
		t.Fatalf("expected 0 errors got %d errors %v", len(errs), errs)
	}

	expected := PageData{
		URL:            "https://crawler-test.com",
		Heading:        "Test Title",
		FirstParagraph: "This is the first paragraph.",
		OutgoingLinks:  []string{"https://crawler-test.com/link1"},
		ImageURLs:      []string{"https://crawler-test.com/image1.jpg"},
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("expected %v got %v", expected, got)
	}
}
