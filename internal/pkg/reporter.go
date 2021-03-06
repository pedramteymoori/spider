package pkg

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

const (
	numberOfWorkers = 100
	channelCapacity = 100000
)

type IReporter interface {
	GetReport(ctx context.Context, body string, websiteURL string) (*Report, error)
}

type Reporter struct {
	report               *Report
	accessibilityChannel chan (*accessibilityRequest)
	Communicator         ICommunicator
}

type Report struct {
	HTMLVersion       string
	Title             string
	Headings          map[string]int32
	InternalLinks     []string
	ExternalLinks     []string
	InAccessibleLinks int32
	HasLogin          bool
}

type loginFormAttrs struct {
	textFound     bool
	passwordFound bool
	submitFound   bool
}

type accessibilityRequest struct {
	link     string
	ctx      context.Context
	doneChan chan (struct{})
}

func NewReporter() *Reporter {
	r := &Reporter{
		accessibilityChannel: make(chan *accessibilityRequest, channelCapacity),
		Communicator:         &Communicator{},
		report: &Report{
			Headings:      make(map[string]int32),
			InternalLinks: make([]string, 0),
			ExternalLinks: make([]string, 0),
		},
	}

	for i := 0; i < numberOfWorkers; i++ {
		go r.checkAccessibility()
	}
	return r
}

func (r *Reporter) GetReport(ctx context.Context, body string, websiteURL string) (*Report, error) {
	parsedURL, err := url.Parse(websiteURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the url : %w", err)
	}

	err = r.traverseHTML(ctx, body, parsedURL)
	if err != nil {
		return nil, err
	}
	return r.report, nil
}

func (r *Reporter) traverseHTML(ctx context.Context, body string, parsedURL *url.URL) error {
	reader := strings.NewReader(body)
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
			case "form":
				r.checkLoginForm(n)
			case "a":
				r.appendLink(n, parsedURL)
			case "h1", "h2", "h3", "h4", "h5", "h6":
				r.appendHeading(n)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	r.checkLinksAccessibility(ctx)
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

func (r *Reporter) appendLink(n *html.Node, parsedURL *url.URL) {
	for _, attr := range n.Attr {
		if attr.Key != "href" {
			continue
		}
		link := attr.Val

		parsed, err := url.Parse(link)
		if err != nil {
			logrus.Errorf("failed to parse link : %s", err)
			continue
		}

		if parsed.Scheme == "" {
			parsed.Scheme = parsedURL.Scheme
		}
		if parsed.Host == "" {
			parsed.Host = parsedURL.Host
		}

		if parsed.Host == parsedURL.Host {
			r.report.InternalLinks = append(r.report.InternalLinks, parsed.String())
		} else {
			r.report.ExternalLinks = append(r.report.ExternalLinks, parsed.String())
		}
		break
	}
}

func (r *Reporter) checkLinksAccessibility(ctx context.Context) {
	var wg sync.WaitGroup

	for _, link := range r.report.ExternalLinks {
		wg.Add(1)
		doneChan := make(chan struct{})
		req := &accessibilityRequest{
			ctx:      ctx,
			doneChan: doneChan,
			link:     link,
		}

		r.accessibilityChannel <- req
		go func() {
			select {
			case <-ctx.Done():
				wg.Done()
			case <-doneChan:
				wg.Done()
			}
		}()
	}
	wg.Wait()
}

func (r *Reporter) checkAccessibility() {
	for {
		request := <-r.accessibilityChannel
		func(req *accessibilityRequest) {
			statusCode, err := r.Communicator.GetStatusCode(req.ctx, req.link)
			if err != nil {
				logrus.Errorf("failed to check accessiblity for %s : %s", req.link, err)
				r.report.InAccessibleLinks++
				req.doneChan <- struct{}{}
				return
			}
			if statusCode >= 400 {
				r.report.InAccessibleLinks++
			}
			req.doneChan <- struct{}{}
		}(request)
	}
}

func (r *Reporter) checkLoginForm(n *html.Node) {
	lfa := &loginFormAttrs{}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		r.checkForInput(c, lfa)
		if lfa.passwordFound { // We can also check for submit and text inputs here
			r.report.HasLogin = true
			return
		}
	}
}

func (r *Reporter) checkForInput(n *html.Node, lfa *loginFormAttrs) {
	if n.Data == "input" {
		for _, attr := range n.Attr {
			checkForAttributes(attr, lfa)
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		r.checkForInput(c, lfa)
	}
}

func checkForAttributes(attr html.Attribute, lfa *loginFormAttrs) {
	if attr.Key != "type" {
		return
	}

	switch attr.Val {
	case "text":
		lfa.textFound = true
	case "password":
		lfa.passwordFound = true
	case "submit":
		lfa.submitFound = true
	}
}
