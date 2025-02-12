package kp

import (
	"net/http"
	"testing"
)

func TestRemoveBraces(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"{hello}", "hello"},
		{"{hello}{world}", "helloworld"},
		{"no_braces", "no_braces"},
		{"{start}middle{end}", "startmiddleend"},
		{"{}", ""},
		{"{a}{b}{c}", "abc"},
	}

	for _, test := range tests {
		result := removeBraces(test.input)
		if result != test.expected {
			t.Errorf("removeBraces(%q) = %q; expected %q", test.input, result, test.expected)
		}
	}
}
func TestSetParam(t *testing.T) {
	tests := []struct {
		path     string
		urlPath  string
		expected map[string]string
	}{
		{"/user/{id}", "/user/123", map[string]string{"id": "123"}},
		{"/user/{id}/profile", "/user/123/profile", map[string]string{"id": "123"}},
		{"/{entity}/{id}", "/product/456", map[string]string{"entity": "product", "id": "456"}},
		{"/{entity}/{id}/details", "/order/789/details", map[string]string{"entity": "order", "id": "789"}},
		{"/static/path", "/static/path", map[string]string{}},
		{"/{a}/{b}/{c}", "/x/y/z", map[string]string{"a": "x", "b": "y", "c": "z"}},
	}

	for _, test := range tests {
		req, _ := http.NewRequest("GET", test.urlPath, nil)
		req = setParam(test.path, req)

		for key, expectedValue := range test.expected {
			value := req.Context().Value(ContextKey(key))
			if value != expectedValue {
				t.Errorf("setParam(%q, %q) = %q; expected %q", test.path, test.urlPath, value, expectedValue)
			}
		}
	}
}
