package spider

import (
	"context"
	"net/url"
	"testing"

	"github.com/pedramteymoori/spider/internal/pkg"
)

func TestRun(t *testing.T) {
	requestURL := "http://localhost:8080/?url=http://www.columbia.edu/~fdc/sample.html"
	websiteURL, err := url.Parse(requestURL)
	if err != nil {
		t.Fatalf("failed to parse url : %s", err)
	}

	s := getMockSpider(t)

	report, err := s.Run(websiteURL)
	if err != nil {
		t.Fatalf("failed to generate report : %s", err)
	}
	if report.HTMLVersion != "HTML 5" {
		t.Errorf("expected HTMLVersion to be HTML 5 but got %s", report.HTMLVersion)
	}
}

func getMockSpider(t *testing.T) *Spider {
	r := &Spider{
		reporter:     getFakeReporter(),
		communicator: getFakeCommunicator(),
	}
	return r
}

type fakeCommunicator struct{}

func (fakeCommunicator) GetBody(ctx context.Context, url string) (string, error) {
	return "", nil
}

func (fakeCommunicator) GetStatusCode(ctx context.Context, url string) (int, error) {
	return 200, nil
}

func getFakeCommunicator() pkg.ICommunicator {
	return &fakeCommunicator{}
}

type fakeReporter struct{}

func (fakeReporter) GetReport(ctx context.Context, body string, websiteURL string) (*pkg.Report, error) {
	r := &pkg.Report{
		HTMLVersion: "HTML 5",
	}
	return r, nil
}

func getFakeReporter() pkg.IReporter {
	return &fakeReporter{}
}
