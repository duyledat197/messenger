package blevesearch

import (
	"context"

	"github.com/blevesearch/bleve/v2"
)

type Index struct {
	bleve.Index
	path string
}

func NewBleveSearch(path string) *Index {
	return &Index{
		path: path,
	}
}

func (s *Index) Start(_ context.Context) error {
	index, err := bleve.New(s.path, bleve.NewIndexMapping())
	if err != nil {
		return err
	}

	s.Index = index

	return nil
}

func (s *Index) Stop(_ context.Context) error {
	return nil
}
