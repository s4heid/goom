package cmd

import (
	gobrowser "github.com/pkg/browser"
)

// Browser opens url in a new window of the default web browser.
//go:generate counterfeiter -o ./fakes/browser.go . Browser
type Browser interface {
	OpenURL(string) error
}

// DefaultBrowser is the default web browser
type DefaultBrowser struct {
}

// OpenURL opens url in the default web browser.
func (e DefaultBrowser) OpenURL(url string) error {
	err := gobrowser.OpenURL(url)
	return err
}

func openURL(url string, b Browser) error {
	err := b.OpenURL(url)
	return err
}
