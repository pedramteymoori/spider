package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pedramteymoori/spider/internal/spider"
	"github.com/sirupsen/logrus"
)

func Execute() {
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.HandleFunc("/", rootHandler)
	logrus.Info("Starting server at port 8080\n")
	logrus.Fatal(http.ListenAndServe(":8080", nil))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	report, err := spider.Run(r.URL)
	if err != nil {
		logrus.WithError(err).Error("failed to respond the request")
		fmt.Fprintf(w, "Internal Error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(report)

	logrus.Info("processed the request successfully")
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {}
