package blocks

import (
	"testing"
)

func TestValidatePath(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		want      string
		wantError bool
	}{
		{
			name:      "empty path",
			path:      "",
			want:      "",
			wantError: false,
		},
		{
			name:      "simple path",
			path:      "a/b/c",
			want:      "a/b/c",
			wantError: false,
		},
		{
			name:      "path with leading ./",
			path:      "./a/b/c",
			want:      "a/b/c",
			wantError: false,
		},
		{
			name:      "path traversal with ..",
			path:      "a/../../../etc/passwd",
			want:      "",
			wantError: true,
		},
		{
			name:      "absolute path",
			path:      "/etc/passwd",
			want:      "",
			wantError: true,
		},
		{
			name:      "path with null byte",
			path:      "a/b\x00c",
			want:      "",
			wantError: true,
		},
		{
			name:      "path too deep",
			path:      "a/b/c/d/e/f/g/h/i/j/k",
			want:      "",
			wantError: true,
		},
		{
			name:      "path at max depth",
			path:      "a/b/c/d/e/f/g/h/i/j",
			want:      "a/b/c/d/e/f/g/h/i/j",
			wantError: false,
		},
		{
			name:      "path with dots cleaned",
			path:      "a/./b/./c",
			want:      "a/b/c",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidatePath(tt.path)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidatePath() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if got != tt.want {
				t.Errorf("ValidatePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateContentType(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		wantError   bool
	}{
		{
			name:        "valid text/plain",
			contentType: "text/plain",
			wantError:   false,
		},
		{
			name:        "valid application/json",
			contentType: "application/json",
			wantError:   false,
		},
		{
			name:        "valid with charset",
			contentType: "text/html; charset=utf-8",
			wantError:   false,
		},
		{
			name:        "valid image type",
			contentType: "image/png",
			wantError:   false,
		},
		{
			name:        "empty content type",
			contentType: "",
			wantError:   true,
		},
		{
			name:        "missing subtype",
			contentType: "text",
			wantError:   true,
		},
		{
			name:        "missing type",
			contentType: "/plain",
			wantError:   true,
		},
		{
			name:        "too many slashes",
			contentType: "text/plain/extra",
			wantError:   true,
		},
		{
			name:        "with null byte",
			contentType: "text/plain\x00",
			wantError:   true,
		},
		{
			name:        "with newline",
			contentType: "text/plain\n",
			wantError:   true,
		},
		{
			name:        "with carriage return",
			contentType: "text/plain\r",
			wantError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateContentType(tt.contentType)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateContentType() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}
