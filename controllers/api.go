package controllers

import (
	"context"
	"strings"

	"github.com/berto/kerbal/services"
	"github.com/pkg/errors"
)

// Item is a bucket asset
type Item string

// Items includes all assets organized by categories
type Items map[string][]Item

// GetItems returns a list of s3 item assets
func GetItems(ctx context.Context) (Items, error) {
	awsService := services.New(ctx)
	if err := awsService.AWSConnect(); err != nil {
		return nil, errors.Wrap(err, "Failed to connect to aws: %s")
	}
	folder := "images"
	items, err := awsService.List(&folder)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to list items: %s")
	}
	return NewItems(items), nil
}

// NewItems converts s3 objects into items
func NewItems(objects []*services.S3Object) Items {
	items := map[string][]Item{}
	for _, obj := range objects {
		splitName := strings.Split(obj.Name, "/")
		if len(splitName) < 3 {
			continue
		}
		folder := splitName[1]
		if obj.Size == 0 {
			if _, ok := items[obj.Name]; !ok {
				items[folder] = []Item{}
			}
			continue
		}
		items[folder] = append(items[folder], Item(splitName[2]))
	}
	return items
}
