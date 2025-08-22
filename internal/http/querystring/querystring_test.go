package querystring

import (
	"net/url"
	"reflect"
	"testing"
)

// testStruct is a struct used for testing Marshal function.
type testStruct struct {
	Name   string `query:"name"`
	Age    int    `query:"age,omitempty"`
	Active bool   `query:"active"`
	Skip   string `query:"-"`
}

// TestMarshal_SimpleStruct tests Marshal with a simple struct and all fields set.
func TestMarshal_SimpleStruct(t *testing.T) {
	input := testStruct{
		Name:   "Alice",
		Age:    30,
		Active: true,
		Skip:   "should be skipped",
	}
	got, err := Marshal(input)
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}
	values, err := url.ParseQuery(string(got))
	if err != nil {
		t.Fatalf("Failed to parse query string: %v", err)
	}
	want := map[string]string{
		"name":   "Alice",
		"age":    "30",
		"active": "1",
	}
	for k, v := range want {
		if values.Get(k) != v {
			t.Errorf("Expected %s=%s, got %s", k, v, values.Get(k))
		}
	}
	if values.Get("Skip") != "" {
		t.Errorf("Expected Skip field to be omitted, got %s", values.Get("Skip"))
	}
}

// TestMarshal_OmitemptyZeroValue tests Marshal with omitempty and zero value.
func TestMarshal_OmitemptyZeroValue(t *testing.T) {
	input := testStruct{
		Name:   "Bob",
		Age:    0, // omitempty, should be omitted
		Active: false,
	}
	got, err := Marshal(input)
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}
	values, err := url.ParseQuery(string(got))
	if err != nil {
		t.Fatalf("Failed to parse query string: %v", err)
	}
	if values.Get("age") != "" {
		t.Errorf("Expected age to be omitted, got %s", values.Get("age"))
	}
	if values.Get("active") != "0" {
		t.Errorf("Expected active=0, got %s", values.Get("active"))
	}
	if values.Get("name") != "Bob" {
		t.Errorf("Expected name=Bob, got %s", values.Get("name"))
	}
}

// TestMarshal_UnsupportedType tests Marshal with a struct containing unsupported types.
func TestMarshal_UnsupportedType(t *testing.T) {
	type unsupported struct {
		Data []string `query:"data"`
	}
	input := unsupported{
		Data: []string{"a", "b"},
	}
	got, err := Marshal(input)
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}
	if string(got) != "" {
		t.Errorf("Expected empty query string for unsupported type, got %s", string(got))
	}
}

// TestMarshal_NoTags tests Marshal with a struct with no query tags.
func TestMarshal_NoTags(t *testing.T) {
	type noTags struct {
		Field1 string
		Field2 int
	}
	input := noTags{
		Field1: "foo",
		Field2: 42,
	}
	got, err := Marshal(input)
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}
	if string(got) != "" {
		t.Errorf("Expected empty query string for struct with no tags, got %s", string(got))
	}
}

// TestMarshal_NonStruct tests Marshal with a non-struct input.
func TestMarshal_NonStruct(t *testing.T) {
	input := 123
	_, err := Marshal(input)
	if err == nil {
		t.Errorf("Expected error for non-struct input, got nil")
	}
}

// TestMarshal_EmptyStruct tests Marshal with an empty struct.
func TestMarshal_EmptyStruct(t *testing.T) {
	type empty struct{}
	input := empty{}
	got, err := Marshal(input)
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}
	if string(got) != "" {
		t.Errorf("Expected empty query string for empty struct, got %s", string(got))
	}
}

// TestMarshal_BoolFalse tests Marshal with a boolean field set to false.
func TestMarshal_BoolFalse(t *testing.T) {
	type boolStruct struct {
		Flag bool `query:"flag"`
	}
	input := boolStruct{Flag: false}
	got, err := Marshal(input)
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}
	values, err := url.ParseQuery(string(got))
	if err != nil {
		t.Fatalf("Failed to parse query string: %v", err)
	}
	if values.Get("flag") != "0" {
		t.Errorf("Expected flag=0, got %s", values.Get("flag"))
	}
}

// TestMarshal_IntTypes tests Marshal with various int types.
func TestMarshal_IntTypes(t *testing.T) {
	type intTypes struct {
		I8  int8  `query:"i8"`
		I16 int16 `query:"i16"`
		I32 int32 `query:"i32"`
		I64 int64 `query:"i64"`
	}
	input := intTypes{
		I8:  1,
		I16: 2,
		I32: 3,
		I64: 4,
	}
	got, err := Marshal(input)
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}
	values, err := url.ParseQuery(string(got))
	if err != nil {
		t.Fatalf("Failed to parse query string: %v", err)
	}
	want := map[string]string{
		"i8":  "1",
		"i16": "2",
		"i32": "3",
		"i64": "4",
	}
	for k, v := range want {
		if values.Get(k) != v {
			t.Errorf("Expected %s=%s, got %s", k, v, values.Get(k))
		}
	}
}

// TestMarshal_TagOmitted tests Marshal with a field tagged as "-".
func TestMarshal_TagOmitted(t *testing.T) {
	type omit struct {
		Field1 string `query:"-"`
		Field2 string `query:"field2"`
	}
	input := omit{
		Field1: "should be omitted",
		Field2: "included",
	}
	got, err := Marshal(input)
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}
	values, err := url.ParseQuery(string(got))
	if err != nil {
		t.Fatalf("Failed to parse query string: %v", err)
	}
	if values.Get("field2") != "included" {
		t.Errorf("Expected field2=included, got %s", values.Get("field2"))
	}
	if _, ok := values["Field1"]; ok {
		t.Errorf("Expected Field1 to be omitted, but found in query string")
	}
}

// TestMarshal_EmptyTag tests Marshal with a field with an empty tag.
func TestMarshal_EmptyTag(t *testing.T) {
	type emptyTag struct {
		Field1 string `query:""`
		Field2 string `query:"field2"`
	}
	input := emptyTag{
		Field1: "should be omitted",
		Field2: "included",
	}
	got, err := Marshal(input)
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}
	values, err := url.ParseQuery(string(got))
	if err != nil {
		t.Fatalf("Failed to parse query string: %v", err)
	}
	if values.Get("field2") != "included" {
		t.Errorf("Expected field2=included, got %s", values.Get("field2"))
	}
	if _, ok := values["Field1"]; ok {
		t.Errorf("Expected Field1 to be omitted, but found in query string")
	}
}

// TestMarshal_MultipleFieldsSameKey tests Marshal with multiple fields using the same key.
func TestMarshal_MultipleFieldsSameKey(t *testing.T) {
	type sameKey struct {
		A string `query:"dup"`
		B string `query:"dup"`
	}
	input := sameKey{
		A: "foo",
		B: "bar",
	}
	got, err := Marshal(input)
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}
	values, err := url.ParseQuery(string(got))
	if err != nil {
		t.Fatalf("Failed to parse query string: %v", err)
	}
	// Both values should be present for the same key.
	gotVals := values["dup"]
	wantVals := []string{"foo", "bar"}
	if !reflect.DeepEqual(gotVals, wantVals) && !reflect.DeepEqual(gotVals, []string{"bar", "foo"}) {
		t.Errorf("Expected dup values %v, got %v", wantVals, gotVals)
	}
}
