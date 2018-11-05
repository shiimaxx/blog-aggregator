package hatenablog

import (
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/shiimaxx/blog-aggregator/structs"
)

const baseURL = "https://blog.hatena.ne.jp"

// Result for hatenablog correction uri
type Result struct {
	Entry []structs.Entry `xml:"entry"`
}

// FetchEntries fetch entry list of hatena blog
func FetchEntries(userID, blogID, apiKey string) ([]structs.Entry, error) {
	endpoint := fmt.Sprintf("%s/%s/%s/atom/entry", baseURL, userID, blogID)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch qiita entries: %s", err.Error())
	}

	req.SetBasicAuth(userID, apiKey)

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

	var r Result

	select {
	case err := <-errCh:
		return nil, fmt.Errorf("failed to fetch hatenablog entries: %s", err.Error())
	case <-doneCh:
		if err := xml.Unmarshal(body, &r); err != nil {
			return nil, fmt.Errorf("failed to fetch hatenablog entries: %s", err.Error())
		}
	}

	return r.Entry, nil
}
