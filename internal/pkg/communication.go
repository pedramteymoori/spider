package pkg

import (
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

func GetBody(url string) (string, error) {
	logrus.Info("trying to fetch : ", url)
	resp, err := http.Get(url)
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
