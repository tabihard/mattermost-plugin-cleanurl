package main

import "testing"

func TestCleanMessage(t *testing.T) {
	tests := []struct {
		name    string
		message string
		extra   map[string]bool
		want    string
	}{
		{
			name:    "strips single utm param",
			message: "https://example.com/article?utm_source=twitter",
			want:    "https://example.com/article",
		},
		{
			name:    "strips multiple tracking params and keeps real ones",
			message: "https://example.com/item?id=42&utm_source=x&fbclid=abc&color=red",
			want:    "https://example.com/item?id=42&color=red",
		},
		{
			name:    "leaves url without tracking params untouched",
			message: "https://example.com/item?id=42&color=red",
			want:    "https://example.com/item?id=42&color=red",
		},
		{
			name:    "leaves url without query untouched",
			message: "check this out https://example.com/path",
			want:    "check this out https://example.com/path",
		},
		{
			name:    "handles markdown link syntax without eating closing paren",
			message: "[link](https://example.com?utm_source=x)",
			want:    "[link](https://example.com)",
		},
		{
			name:    "cleans multiple urls in one message",
			message: "https://a.com?utm_source=x and https://b.com?gclid=y",
			want:    "https://a.com and https://b.com",
		},
		{
			name:    "respects admin configured extra params",
			message: "https://example.com?ref=homepage&id=1",
			extra:   map[string]bool{"ref": true},
			want:    "https://example.com?id=1",
		},
		{
			name:    "removes all params leaving no trailing question mark",
			message: "https://example.com?utm_source=x&utm_medium=y",
			want:    "https://example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := cleanMessage(tt.message, tt.extra)
			if got != tt.want {
				t.Errorf("cleanMessage(%q) = %q, want %q", tt.message, got, tt.want)
			}
		})
	}
}
