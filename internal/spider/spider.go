package spider

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/pedramteymoori/spider/internal/pkg"
)

func Run(requestURL *url.URL) (*pkg.Report, error) {
	websiteURL := requestURL.Query().Get("url")
	if websiteURL == "" {
		return nil, errors.New("Please provide webiste url")
	}
	websiteURL = strings.Trim(websiteURL, "\"")

	parsedURL, err := url.Parse(websiteURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the url : %w", err)
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	body, err := pkg.GetBody(ctx, websiteURL)
	if err != nil {
		return nil, fmt.Errorf("error in fetch web page : %w", err)
	}

	reporter := pkg.NewReporter(parsedURL, body)

	report, err := reporter.GetReport(ctx)
	if err != nil {
		return nil, fmt.Errorf("error in parse web page : %w", err)
	}
	return report, nil
}
