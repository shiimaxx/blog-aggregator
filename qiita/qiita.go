package qiita

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/shiimaxx/blog-aggregator/structs"
)

const baseURL = "https://qiita.com/api/v2"

// FetchEntries fetch qiita entries of specified user id
func FetchEntries(userID string) ([]structs.Entry, error) {
	endpoint := fmt.Sprintf("%s/users/%s/items", baseURL, userID)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch qiita entries: %s", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	var body []byte
	errCh := make(chan error)
	doneCh := make(chan struct{})
	go func() {
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			errCh <- err
			return
		}
		body, err = ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		if err != nil {
			errCh <- err
			return
		}
		doneCh <- struct{}{}
	}()

	var e []structs.Entry

	select {
	case err := <-errCh:
		return nil, fmt.Errorf("failed to fetch qiita entries: %s", err.Error())
	case <-doneCh:
		if err := json.Unmarshal(body, &e); err != nil {
			return nil, fmt.Errorf("failed to fetch qiita entries: %s", err.Error())
		}
	}

	return e, nil
}
