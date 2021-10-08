package pkg

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type ICommunicator interface {
	GetBody(ctx context.Context, url string) (string, error)
	GetStatusCode(ctx context.Context, url string) (int, error)
}

type Communicator struct{}

const (
	httpTimeout = 5 * time.Second
)

func (*Communicator) GetBody(ctx context.Context, url string) (string, error) {
	logrus.Info("trying to fetch : ", url)

	reqCtx, cancel := context.WithTimeout(ctx, httpTimeout)
	defer cancel()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Fatalf("%v", err)
	}

	req = req.WithContext(reqCtx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	logrus.Info("response code is : ", resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (*Communicator) GetStatusCode(ctx context.Context, url string) (int, error) {
	logrus.Info("trying to fetch : ", url)

	reqCtx, cancel := context.WithTimeout(ctx, httpTimeout)
	defer cancel()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Fatalf("%v", err)
	}

	req = req.WithContext(reqCtx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}
