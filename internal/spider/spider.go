package spider

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/pedramteymoori/spider/internal/pkg"
)

func Run(requestURL *url.URL) (*pkg.Report, error) {
	url := requestURL.Query().Get("url")
	if url == "" {
		return nil, errors.New("Please provide webiste url")
	}
	url = strings.Trim(url, "\"")

	body, err := pkg.GetBody(url)
	if err != nil {
		return nil, fmt.Errorf("error in fetch web page : %w", err)
	}

	reporter := pkg.NewReporter(body)
	report, err := reporter.GetReport()
	if err != nil {
		return nil, fmt.Errorf("error in parse web page : %w", err)
	}
	return report, nil
}
