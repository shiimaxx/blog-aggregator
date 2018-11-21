package blogservice

import (
	"net/http"
	"sync"

	"github.com/shiimaxx/blog-aggregator/structs"
	"golang.org/x/sync/errgroup"
)

var HTTPClient = http.DefaultClient

type BlogService struct {
	FetchFunc []func() ([]structs.Entry, error)
}

func (b *BlogService) Add(fn func() ([]structs.Entry, error)) {
	b.FetchFunc = append(b.FetchFunc, fn)
}

func (b *BlogService) Fetch() ([]structs.Entry, error) {
	eg := errgroup.Group{}
	var entries []structs.Entry
	var mu sync.Mutex
	for _, fn := range b.FetchFunc {
		fn := fn
		eg.Go(func() error {
			e, err := fn()
			if err != nil {
				return err
			}
			mu.Lock()
			entries = append(entries, e...)
			mu.Unlock()
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return entries, nil
}
