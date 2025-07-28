// Package querystring provides a marshaller that encodes structs into URL query strings.
// It uses struct tags similar to encoding/json for flexible and readable control.
//
// Struct fields must be tagged with `query:"key"` to be included in the output.
// Fields tagged with `query:"key,omitempty"` are omitted if they hold the zero value for their type.
// Boolean values are always encoded as "1" (true) or "0" (false).
//
// Example:
//
//	type SearchParams struct {
//	    Offset     int    `query:"offset"`
//	    Limit      int    `query:"limit,omitempty"`
//	    Query      string `query:"query,omitempty"`
//	    TitlesOnly bool   `query:"titles_only"`
//	}
//
//	params := SearchParams{
//	    Offset:     10,
//	    TitlesOnly: true,
//	}
//
//	q, _ := querystring.MarshalQuery(params)
//	// q: []byte("offset=10&titles_only=1")
package querystring

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// Marshal converts a struct with `query:"key,omitempty"` tags into a query string.
// Booleans are encoded as "1"/"0". Omits zero values if `omitempty` is set.
func Marshal(v any) ([]byte, error) {
	rv := reflect.ValueOf(v)
	rt := reflect.TypeOf(v)

	if rv.Kind() != reflect.Struct {
		return nil, fmt.Errorf("MarshalQuery: expected struct, got %s", rv.Kind())
	}

	values := url.Values{}

	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)
		fieldType := rt.Field(i)

		tag := fieldType.Tag.Get("query")
		if tag == "-" {
			continue
		}

		tagParts := strings.Split(tag, ",")
		if len(tagParts) == 0 || tagParts[0] == "" {
			continue
		}

		key := tagParts[0]
		omitempty := false
		for _, opt := range tagParts[1:] {
			if opt == "omitempty" {
				omitempty = true
			}
		}

		if omitempty && field.IsZero() {
			continue
		}

		var str string
		switch field.Kind() {
		case reflect.String:
			str = field.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			str = strconv.FormatInt(field.Int(), 10)
		case reflect.Bool:
			if field.Bool() {
				str = "1"
			} else {
				str = "0"
			}
		default:
			continue // unsupported type
		}

		values.Add(key, str)
	}

	return []byte(values.Encode()), nil
}
