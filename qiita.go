package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type qiitaEntriy struct {
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
}

const baseURL = "https://qiita.com/api/v2"

func fetchEntries(userID string) ([]qiitaEntriy, error) {
	endpoint := fmt.Sprintf("%s/users/%s/items", baseURL, userID)

	res, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var q []qiitaEntriy
	if err := json.Unmarshal(body, &q); err != nil {
		return nil, err
	}

	return q, nil
}
