package main

import (
	"github.com/pkg/browser"
)

func OpenURL(url string) error {
	return browser.OpenURL(url)
}
