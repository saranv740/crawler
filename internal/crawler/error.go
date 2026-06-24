package crawler

import "errors"

var (
	ErrInvalidURL   = errors.New("invalid url")
	ExternalHostErr = errors.New("external host")
	URLSanityErr    = errors.New("url sanitizing failed")
	PageFetchErr    = errors.New("error in fetching page")
)
