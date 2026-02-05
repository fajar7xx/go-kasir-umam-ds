package handlers

import "time"

const (
	requestTimeout = 5 * time.Second // Centralized timeout
	maxNameLength  = 255             // Maximum length of product name
)
