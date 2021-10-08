package pkg

import (
	"context"
	"os"
	"testing"
)

func TestGetReport(t *testing.T) {
	url := "http://www.columbia.edu/~fdc/sample.html"

	data, err := os.ReadFile("./test/page.html")
	if err != nil {
		t.Errorf("reading sample file failed : %s", err)
	}

	r := getMockReporter(t)
	report, err := r.GetReport(context.Background(), string(data), url)
	if err != nil {
		t.Fatalf("failed to generate report : %s", err)
	}
	if report.HTMLVersion != "HTML 5" {
		t.Errorf("expected HTMLVersion to be HTML 5 but got %s", report.HTMLVersion)
	}
	if report.Title != "Sample Web Page" {
		t.Errorf("expected Title to be Sample Web Page but got %s", report.Title)
	}
	if report.Headings["h2"] != 1 {
		t.Errorf(`expected report.Headings["h2"] to be 1 but got %d`, report.Headings["h2"])
	}
	if report.HasLogin != false {
		t.Errorf(`expected report.HasLogin to be false but got %t`, report.HasLogin)
	}
}

func getMockReporter(t *testing.T) *Reporter {
	r := &Reporter{
		accessibilityChannel: make(chan *accessibilityRequest, 10),
		Communicator:         getFakeCommunicator(),
		report: &Report{
			Headings:      make(map[string]int32),
			InternalLinks: make([]string, 0),
			ExternalLinks: make([]string, 0),
		},
	}
	go r.checkAccessibility()
	return r
}

type fakeCommunicator struct{}

func (fakeCommunicator) GetBody(ctx context.Context, url string) (string, error) {
	return "", nil
}

func (fakeCommunicator) GetStatusCode(ctx context.Context, url string) (int, error) {
	return 200, nil
}

func getFakeCommunicator() ICommunicator {
	return &fakeCommunicator{}
}
