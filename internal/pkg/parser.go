package pkg

import (
	"strings"

	"golang.org/x/net/html"
)

type Reporter struct {
	data   string
	report *Report
}

type Report struct {
	HTMLVersion string
	Title       string
	Headings    map[string]int32
	Links       []string
}

func NewReporter(data string) *Reporter {
	r := &Reporter{
		data: data,
		report: &Report{
			Headings: make(map[string]int32),
			Links:    make([]string, 0),
		},
	}
	return r
}

func (r *Reporter) GetReport() (*Report, error) {
	err := r.parseHTML()
	if err != nil {
		return nil, err
	}
	return r.report, nil
}

func (r *Reporter) parseHTML() error {
	reader := strings.NewReader(r.data)
	doc, err := html.Parse(reader)
	if err != nil {
		return err
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.DoctypeNode {
			r.setHTMLVersion(n)
		}
		if n.Type == html.ElementNode {
			switch n.Data {
			case "title":
				r.setTitle(n)
			case "a":
				r.appendLink(n)
			case "h1", "h2", "h3", "h4", "h5", "h6":
				r.appendHeading(n)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return nil
}

func (r *Reporter) setHTMLVersion(n *html.Node) {
	if n.Data == "html" {
		r.report.HTMLVersion = "HTML 5"
	} else if strings.HasPrefix(n.Data, "HTML PUBLIC \"-//W3C//DTD HTML 4.01") {
		r.report.HTMLVersion = "HTML 4.01"
	} else if strings.HasPrefix(n.Data, "html PUBLIC \"-//W3C//DTD XHTML 1.0") {
		r.report.HTMLVersion = "XHTML 1.0"
	} else if strings.HasPrefix(n.Data, "html PUBLIC \"-//W3C//DTD XHTML 1.1") {
		r.report.HTMLVersion = "XHTML 1.1"
	}
}

func (r *Reporter) setTitle(n *html.Node) {
	r.report.Title = n.FirstChild.Data
}

func (r *Reporter) appendHeading(n *html.Node) {
	if v, ok := r.report.Headings[n.Data]; ok {
		r.report.Headings[n.Data] = v + 1
	} else {
		r.report.Headings[n.Data] = 1
	}
}

func (r *Reporter) appendLink(n *html.Node) {
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			r.report.Links = append(r.report.Links, attr.Val)
			break
		}
	}
}
