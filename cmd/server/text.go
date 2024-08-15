package main

import (
	"html/template"
	"log/slog"
	"os"
	"sort"
)

type TextItem struct {
	Name    string
	Path    string
	Summary string
	Data    template.HTML
}

func (ti *TextItem) SafeName() string {
	if ti == nil {
		return ""
	}
	return ti.Name
}

type TextCache struct {
	cache map[string]*TextItem
	items []*TextItem
	names []string
}

func NewTextCache(items ...*TextItem) *TextCache {
	tc := &TextCache{
		cache: make(map[string]*TextItem),
		items: make([]*TextItem, 0, len(items)),
		names: make([]string, 0, len(items)),
	}
	for _, item := range items {
		if data, err := os.ReadFile(item.Path); err != nil {
			slog.Warn("Unable to load text file",
				"Path", item.Path, "error", err.Error())
		} else if len(data) > 0 {
			item.Data = template.HTML(data)
			tc.cache[item.Name] = item
			tc.names = append(tc.names, item.Name)
		} else {
			slog.Warn("Loaded empty text file", "Path", item.Path)
		}
	}
	sort.Strings(tc.names)
	for _, name := range tc.names {
		tc.items = append(tc.items, tc.cache[name])
	}
	return tc
}

func (tc *TextCache) HasText() bool {
	return len(tc.cache) > 0
}

func (tc *TextCache) TextItem(path string) *TextItem {
	var found bool
	var item *TextItem
	if item, found = tc.cache[path]; found && item != nil {
		slog.Error("Text not found", "Path", path)
	} else {
		item = &TextItem{
			Path:    path,
			Summary: "Text Not Found",
			Data:    template.HTML("Path: " + path),
		}
	}
	return item
}

func (tc *TextCache) TextItems() []*TextItem {
	return tc.items
}
