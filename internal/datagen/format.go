// internal/datagen/format.go
package datagen

import (
	"github.com/brianvoe/gofakeit/v7"
)

// generateByFormat returns a fake value for a known OpenAPI format string.
// Returns (value, true) if format is recognized, ("", false) otherwise.
func generateByFormat(format string) (any, bool) {
	switch format {
	case "email":
		return gofakeit.Email(), true
	case "uuid":
		return gofakeit.UUID(), true
	case "date":
		return gofakeit.Date().Format("2006-01-02"), true
	case "date-time":
		return gofakeit.Date().Format("2006-01-02T15:04:05Z"), true
	case "uri", "url":
		return gofakeit.URL(), true
	case "hostname":
		return gofakeit.DomainName(), true
	case "ipv4":
		return gofakeit.IPv4Address(), true
	case "ipv6":
		return gofakeit.IPv6Address(), true
	case "password":
		return gofakeit.Password(true, true, true, true, false, 12), true
	case "phone":
		return gofakeit.Phone(), true
	default:
		return nil, false
	}
}
