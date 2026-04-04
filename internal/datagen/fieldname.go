// internal/datagen/fieldname.go
package datagen

import (
	"strings"

	"github.com/brianvoe/gofakeit/v7"
)

// generateByFieldName infers data from the field name and an optional description hint.
// Returns (value, true) if a semantic rule matches, (nil, false) otherwise.
func generateByFieldName(fieldName, description string) (any, bool) {
	lower := strings.ToLower(fieldName)
	descLower := strings.ToLower(description)

	switch {
	case lower == "email" || strings.HasSuffix(lower, "_email") || strings.HasSuffix(lower, "email"):
		return gofakeit.Email(), true

	// "name" — disambiguated by description context.
	// file/attachment context → filename; product/item context → product-like noun phrase;
	// default → person name.
	case lower == "name" || lower == "full_name" || lower == "fullname":
		if descContainsAny(descLower, "file", "filename", "attachment") {
			return generateFilename(), true
		}
		if descContainsAny(descLower, "product", "item", "service", "category", "tag") {
			w1, w2 := gofakeit.Word(), gofakeit.Word()
			return strings.ToUpper(w1[:1]) + w1[1:] + " " + strings.ToUpper(w2[:1]) + w2[1:], true
		}
		return gofakeit.Name(), true // default: person name

	// Explicit filename fields always produce a filename with extension.
	case lower == "filename" || lower == "file_name":
		return generateFilename(), true

	case lower == "first_name" || lower == "firstname":
		return gofakeit.FirstName(), true
	case lower == "last_name" || lower == "lastname" || lower == "surname":
		return gofakeit.LastName(), true
	case lower == "username":
		return gofakeit.Username(), true
	case lower == "phone" || lower == "phone_number" || lower == "phonenumber":
		return gofakeit.Phone(), true
	case lower == "address" || lower == "street_address":
		return gofakeit.Street(), true
	case lower == "city":
		return gofakeit.City(), true
	case lower == "country":
		return gofakeit.Country(), true
	case lower == "zip" || lower == "zipcode" || lower == "postal_code":
		return gofakeit.Zip(), true
	case lower == "url" || lower == "website":
		return gofakeit.URL(), true
	case lower == "description" || lower == "bio" || lower == "summary":
		return gofakeit.Sentence(8), true
	case lower == "title":
		return gofakeit.JobTitle(), true
	case lower == "company" || lower == "company_name":
		return gofakeit.Company(), true
	case lower == "id" || strings.HasSuffix(lower, "_id") || strings.HasSuffix(lower, "id"):
		return gofakeit.UUID(), true
	case lower == "age":
		return int64(gofakeit.Number(18, 80)), true
	case lower == "price" || lower == "amount" || lower == "cost":
		return gofakeit.Price(1, 1000), true
	case lower == "quantity" || lower == "count":
		return int64(gofakeit.Number(1, 100)), true
	default:
		return nil, false
	}
}

// generateFilename returns a realistic-looking filename with an extension.
func generateFilename() string {
	extensions := []string{"txt", "pdf", "png", "jpg", "csv", "json", "xml", "zip"}
	ext := extensions[gofakeit.Number(0, len(extensions)-1)]
	return gofakeit.LetterN(8) + "." + ext
}

// descContainsAny returns true if descLower contains any of the given keywords.
func descContainsAny(descLower string, keywords ...string) bool {
	for _, kw := range keywords {
		if strings.Contains(descLower, kw) {
			return true
		}
	}
	return false
}
