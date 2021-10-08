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

type Spider struct {
	communicator pkg.ICommunicator
	reporter     pkg.IReporter
}

func NewSpider() *Spider {
	return &Spider{
		communicator: &pkg.Communicator{},
		reporter:     pkg.NewReporter(),
	}
}

func (s *Spider) Run(requestURL *url.URL) (*pkg.Report, error) {
	websiteURL := requestURL.Query().Get("url")
	if websiteURL == "" {
		return nil, errors.New("Please provide webiste url")
	}
	websiteURL = strings.Trim(websiteURL, "\"")

	ctx, cancelFunc := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelFunc()

	body, err := s.communicator.GetBody(ctx, websiteURL)
	if err != nil {
		return nil, fmt.Errorf("error in fetch web page : %w", err)
	}

	report, err := s.reporter.GetReport(ctx, body, websiteURL)
	if err != nil {
		return nil, fmt.Errorf("error in parse web page : %w", err)
	}
	return report, nil
}
