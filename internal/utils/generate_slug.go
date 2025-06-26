package utils

import (
	"regexp"
	"strings"
)

// Helper function to generate slug
func GenerateSlug(name string) string {
	// Replace spaces with hyphens and convert to lowercase
	slug := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	// Remove any non-alphanumeric characters (except hyphens)
	reg := regexp.MustCompile("[^a-z0-9-]+")
	slug = reg.ReplaceAllString(slug, "")
	return slug
}
