package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Time struct {
	Unix    interface{} `json:"unix"`
	Natural interface{} `json:"natural"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", http.StripPrefix("/", fs))

	http.HandleFunc("/api/", handler)
	log.Fatal(http.ListenAndServe("localhost:"+port, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	var timeNow time.Time
	var timeJSON Time
	const layout = "2006-Jan-02"

	urlParam := strings.TrimPrefix(r.URL.Path, "/api/")
	timeNow, _ = time.Parse(layout, urlParam)

	if timeNow.IsZero() {
		i, err := strconv.ParseInt(urlParam, 10, 64)
		if err == nil {
			timeNow = time.Unix(i, 0)
			timeJSON = Time{
				Unix:    timeNow.Unix(),
				Natural: timeNow.Format(time.UnixDate),
			}
		}
	} else {
		timeJSON = Time{
			Unix:    timeNow.Unix(),
			Natural: timeNow.Format(time.UnixDate),
		}
	}

	js, err := json.Marshal(timeJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
