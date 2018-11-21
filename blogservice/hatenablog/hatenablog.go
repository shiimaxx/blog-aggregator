package hatenablog

import (
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/shiimaxx/blog-aggregator/blogservice"
	"github.com/shiimaxx/blog-aggregator/structs"
	"golang.org/x/tools/blog/atom"
)

const baseURL = "https://blog.hatena.ne.jp"

// Result for hatenablog correction uri
type Result struct {
	Entries []atom.Entry `xml:"entry"`
}

// FetchEntries fetch entry list of hatena blog
func FetchEntries(ctx context.Context, userID, blogID, apiKey string) ([]structs.Entry, error) {
	endpoint := fmt.Sprintf("%s/%s/%s/atom/entry", baseURL, userID, blogID)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch hatenablog entries: %s", err.Error())
	}

	req.SetBasicAuth(userID, apiKey)

	req = req.WithContext(ctx)

	var body []byte
	errCh := make(chan error)
	doneCh := make(chan struct{})
	go func() {
		res, err := blogservice.HTTPClient.Do(req)
		if err != nil {
			errCh <- err
			return
		}

		if res.StatusCode != http.StatusOK {
			errCh <- errors.New("request failed")
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
			return nil, fmt.Errorf("failed to parse xml: %s", err.Error())
		}
	}

	var entries []structs.Entry

	for _, e := range r.Entries {
		title := e.Title

		var url string
		for _, l := range e.Link {
			if l.Rel == "alternate" {
				url = l.Href
			}
		}

		layout := "2006-01-02T15:04:05-07:00"
		createdAt, err := time.Parse(layout, string(e.Published))
		if err != nil {
			return nil, err
		}

		entries = append(entries, structs.Entry{
			Title:     title,
			URL:       url,
			CreatedAt: createdAt,
		})
	}

	return entries, nil
}
