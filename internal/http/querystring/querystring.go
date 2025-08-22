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
func Marshal(v interface{}) ([]byte, error) {
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
