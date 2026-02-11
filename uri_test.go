package fastbelt

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseURI(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected URI
		hasError bool
	}{
		{
			name:  "complete HTTP URL",
			input: "http://example.com:8080/path/to/resource?query=value&foo=bar#section",
			expected: URI{
				Scheme:    "http",
				Authority: "example.com:8080",
				Path:      "/path/to/resource",
				Query:     "query=value&foo=bar",
				Fragment:  "section",
			},
			hasError: false,
		},
		{
			name:  "HTTPS URL without port",
			input: "https://www.example.com/api/v1/users",
			expected: URI{
				Scheme:    "https",
				Authority: "www.example.com",
				Path:      "/api/v1/users",
				Query:     "",
				Fragment:  "",
			},
			hasError: false,
		},
		{
			name:  "URL with query but no fragment",
			input: "https://search.example.com/search?q=golang&type=code",
			expected: URI{
				Scheme:    "https",
				Authority: "search.example.com",
				Path:      "/search",
				Query:     "q=golang&type=code",
				Fragment:  "",
			},
			hasError: false,
		},
		{
			name:  "URL with fragment but no query",
			input: "https://docs.example.com/guide#installation",
			expected: URI{
				Scheme:    "https",
				Authority: "docs.example.com",
				Path:      "/guide",
				Query:     "",
				Fragment:  "installation",
			},
			hasError: false,
		},
		{
			name:  "file URI with three slashes",
			input: "file:///home/user/document.txt",
			expected: URI{
				Scheme:    "file",
				Authority: "",
				Path:      "/home/user/document.txt",
				Query:     "",
				Fragment:  "",
			},
			hasError: false,
		},
		{
			name:  "relative path only",
			input: "/path/to/resource",
			expected: URI{
				Scheme:    "",
				Authority: "",
				Path:      "/path/to/resource",
				Query:     "",
				Fragment:  "",
			},
			hasError: false,
		},
		{
			name:  "path with query and fragment",
			input: "/api/users?active=true#list",
			expected: URI{
				Scheme:    "",
				Authority: "",
				Path:      "/api/users",
				Query:     "active=true",
				Fragment:  "list",
			},
			hasError: false,
		},
		{
			name:  "scheme only",
			input: "custom:",
			expected: URI{
				Scheme:    "custom",
				Authority: "",
				Path:      "",
				Query:     "",
				Fragment:  "",
			},
			hasError: false,
		},
		{
			name:  "authority with userinfo",
			input: "ftp://user:pass@ftp.example.com/files/",
			expected: URI{
				Scheme:    "ftp",
				Authority: "user:pass@ftp.example.com",
				Path:      "/files/",
				Query:     "",
				Fragment:  "",
			},
			hasError: false,
		},
		{
			name:  "empty string",
			input: "",
			expected: URI{
				Scheme:    "",
				Authority: "",
				Path:      "",
				Query:     "",
				Fragment:  "",
			},
			hasError: false,
		},
		{
			name:  "complex query string",
			input: "https://api.example.com/v2/search?q=golang+testing&lang=en&sort=date&page=1",
			expected: URI{
				Scheme:    "https",
				Authority: "api.example.com",
				Path:      "/v2/search",
				Query:     "q=golang+testing&lang=en&sort=date&page=1",
				Fragment:  "",
			},
			hasError: false,
		},
		{
			name:  "URI with port and special chars in fragment",
			input: "http://localhost:3000/app#section/subsection?param=value",
			expected: URI{
				Scheme:    "http",
				Authority: "localhost:3000",
				Path:      "/app",
				Query:     "",
				Fragment:  "section/subsection?param=value",
			},
			hasError: false,
		},
		{
			name:  "mailto scheme",
			input: "mailto:user@example.com?subject=Hello&body=Hi%20there",
			expected: URI{
				Scheme:    "mailto",
				Authority: "",
				Path:      "user@example.com",
				Query:     "subject=Hello&body=Hi there",
				Fragment:  "",
			},
			hasError: false,
		},
		{
			name:  "encoded characters in path (vscode URI)",
			input: "file:///c%3A/Users/Document/%23source.go",
			expected: URI{
				Scheme:    "file",
				Authority: "",
				Path:      "/C:/Users/Document/#source.go",
				Query:     "",
				Fragment:  "",
			},
			hasError: false,
		},
		{
			name:  "windows file URI",
			input: "file:///C:/Users/John%20Doe/Documents/my%20file.txt",
			expected: URI{
				Scheme:    "file",
				Authority: "",
				Path:      "/C:/Users/John Doe/Documents/my file.txt",
				Query:     "",
				Fragment:  "",
			},
			hasError: false,
		},
		{
			name:  "multiple question marks",
			input: "http://example.com/path?param1=value?param2=value2",
			expected: URI{
				Scheme:    "http",
				Authority: "example.com",
				Path:      "/path",
				Query:     "param1=value?param2=value2",
				Fragment:  "",
			},
			hasError: false,
		},
		{
			name:  "multiple hash symbols",
			input: "http://example.com/path#fragment#more",
			expected: URI{
				Scheme:    "http",
				Authority: "example.com",
				Path:      "/path",
				Query:     "",
				Fragment:  "fragment#more",
			},
			hasError: false,
		},
		{
			name:  "query without path",
			input: "http://example.com?param=value",
			expected: URI{
				Scheme:    "http",
				Authority: "example.com",
				Path:      "",
				Query:     "param=value",
				Fragment:  "",
			},
			hasError: false,
		},
		{
			name:  "fragment without path or query",
			input: "http://example.com#section",
			expected: URI{
				Scheme:    "http",
				Authority: "example.com",
				Path:      "",
				Query:     "",
				Fragment:  "section",
			},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseURI(tt.input)

			if tt.hasError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected.Scheme, result.Scheme)
			assert.Equal(t, tt.expected.Authority, result.Authority)
			assert.Equal(t, tt.expected.Path, result.Path)
			assert.Equal(t, tt.expected.Query, result.Query)
			assert.Equal(t, tt.expected.Fragment, result.Fragment)
		})
	}
}

func TestURI_StringUnencoded(t *testing.T) {
	tests := []struct {
		name     string
		uri      URI
		expected string
	}{
		{
			name: "complete URI",
			uri: URI{
				Scheme:    "https",
				Authority: "example.com:443",
				Path:      "/path/to/resource",
				Query:     "param=value",
				Fragment:  "section",
			},
			expected: "https://example.com:443/path/to/resource?param=value#section",
		},
		{
			name: "URI without authority",
			uri: URI{
				Scheme:   "file",
				Path:     "/home/user/file.txt",
				Query:    "",
				Fragment: "",
			},
			expected: "file:/home/user/file.txt",
		},
		{
			name: "URI with empty authority",
			uri: URI{
				Scheme:    "file",
				Authority: "",
				Path:      "/home/user/file.txt",
				Query:     "",
				Fragment:  "",
			},
			expected: "file:/home/user/file.txt",
		},
		{
			name: "path only",
			uri: URI{
				Path: "/api/data",
			},
			expected: "/api/data",
		},
		{
			name: "query and fragment only",
			uri: URI{
				Query:    "search=golang",
				Fragment: "results",
			},
			expected: "?search=golang#results",
		},
		{
			name:     "empty URI",
			uri:      URI{},
			expected: "",
		},
		{
			name: "scheme and path only",
			uri: URI{
				Scheme: "custom",
				Path:   "/resource",
			},
			expected: "custom:/resource",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.uri.StringUnencoded()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewURI(t *testing.T) {
	tests := []struct {
		name      string
		scheme    string
		authority string
		path      string
		query     string
		fragment  string
		expected  URI
		hasError  bool
	}{
		{
			name:      "valid HTTP URI",
			scheme:    "http",
			authority: "example.com",
			path:      "/path",
			query:     "q=value",
			fragment:  "section",
			expected: URI{
				Scheme:    "http",
				Authority: "example.com",
				Path:      "/path",
				Query:     "q=value",
				Fragment:  "section",
			},
			hasError: false,
		},
		{
			name:      "valid scheme with special characters",
			scheme:    "custom-scheme.v1+ext",
			authority: "",
			path:      "/resource",
			query:     "",
			fragment:  "",
			expected: URI{
				Scheme:    "custom-scheme.v1+ext",
				Authority: "",
				Path:      "/resource",
				Query:     "",
				Fragment:  "",
			},
			hasError: false,
		},
		{
			name:      "invalid scheme starting with digit",
			scheme:    "1invalid",
			authority: "example.com",
			path:      "/path",
			query:     "",
			fragment:  "",
			expected:  URI{},
			hasError:  true,
		},
		{
			name:      "invalid scheme with invalid character",
			scheme:    "test@scheme",
			authority: "example.com",
			path:      "/path",
			query:     "",
			fragment:  "",
			expected:  URI{},
			hasError:  true,
		},
		{
			name:      "empty scheme",
			scheme:    "",
			authority: "example.com",
			path:      "/path",
			query:     "",
			fragment:  "",
			expected:  URI{},
			hasError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewURI(tt.scheme, tt.authority, tt.path, tt.query, tt.fragment)

			if tt.hasError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected.Scheme, result.Scheme)
			assert.Equal(t, tt.expected.Authority, result.Authority)
			assert.Equal(t, tt.expected.Path, result.Path)
			assert.Equal(t, tt.expected.Query, result.Query)
			assert.Equal(t, tt.expected.Fragment, result.Fragment)
		})
	}
}

func TestParseAndStringUnencoded_RoundTrip(t *testing.T) {
	testURIs := []string{
		"https://example.com:8080/path?query=value#fragment",
		"http://user:pass@host.com/path",
		"custom://authority/path?q=v",
		"/relative/path?param=value",
		"mailto:test@example.com",
		"ftp://ftp.example.com/files/",
		"https://api.github.com/repos/owner/repo/issues?state=open&labels=bug",
	}

	for _, original := range testURIs {
		t.Run(original, func(t *testing.T) {
			parsed, err := ParseURI(original)
			assert.NoError(t, err, "failed to parse URI %q", original)

			reconstructed := parsed.StringUnencoded()
			assert.Equal(t, original, reconstructed, "round trip failed")
		})
	}
}

func TestFileURI_ExpectedBehavior(t *testing.T) {
	// This test documents the current behavior with file URIs containing empty authority
	// The input "file:///path" becomes "file:/path" after parsing and reconstructing
	// because the String() method only adds "//" when authority is non-empty
	tests := []struct {
		input    string
		expected string
	}{
		{"file:///absolute/path/to/file", "file:/absolute/path/to/file"},
		{"file:///C:/Windows/file.txt", "file:/C:/Windows/file.txt"},
		{"file://server/share/file.txt", "file://server/share/file.txt"}, // non-empty authority preserved
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			parsed, err := ParseURI(tt.input)
			assert.NoError(t, err, "failed to parse URI %q", tt.input)

			reconstructed := parsed.StringUnencoded()
			assert.Equal(t, tt.expected, reconstructed)
		})
	}
}

func TestIsValidScheme(t *testing.T) {
	tests := []struct {
		name     string
		scheme   string
		expected bool
	}{
		{"valid lowercase", "http", true},
		{"valid uppercase", "HTTP", true},
		{"valid mixed case", "Http", true},
		{"valid with digits", "h2", true},
		{"valid with hyphen", "custom-scheme", true},
		{"valid with plus", "scheme+ext", true},
		{"valid with dot", "scheme.v1", true},
		{"complex valid", "My-Custom.Scheme+Extension2", true},
		{"empty string", "", false},
		{"starts with digit", "2scheme", false},
		{"starts with hyphen", "-scheme", false},
		{"starts with plus", "+scheme", false},
		{"starts with dot", ".scheme", false},
		{"contains invalid char", "scheme@test", false},
		{"contains space", "my scheme", false},
		{"contains underscore", "my_scheme", false},
		{"contains colon", "scheme:test", false},
		{"contains slash", "scheme/test", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidScheme(tt.scheme)
			assert.Equal(t, tt.expected, result, "isValidScheme(%q)", tt.scheme)
		})
	}
}

func TestURI_String(t *testing.T) {
	tests := []struct {
		name     string
		uri      URI
		expected string
	}{
		{
			name: "basic URI without special characters",
			uri: URI{
				Scheme:    "https",
				Authority: "example.com",
				Path:      "/api/users",
				Query:     "active=true",
				Fragment:  "section",
			},
			expected: "https://example.com/api/users?active%3Dtrue#section",
		},
		{
			name: "spaces in different components",
			uri: URI{
				Scheme:    "https",
				Authority: "example server.com",
				Path:      "/path with spaces/file",
				Query:     "param=value with spaces",
				Fragment:  "section with spaces",
			},
			expected: "https://example%20server.com/path%20with%20spaces/file?param%3Dvalue%20with%20spaces#section%20with%20spaces",
		},
		{
			name: "special characters in authority",
			uri: URI{
				Scheme:    "http",
				Authority: "user:pass@host[::1]:8080",
				Path:      "/path",
				Query:     "",
				Fragment:  "",
			},
			expected: "http://user:pass%40host[::1]:8080/path",
		},
		{
			name: "IPv6 address in authority (preserved)",
			uri: URI{
				Scheme:    "http",
				Authority: "[::1]:8080",
				Path:      "/path",
				Query:     "",
				Fragment:  "",
			},
			expected: "http://[::1]:8080/path",
		},
		{
			name: "path with forward slashes preserved",
			uri: URI{
				Scheme:    "https",
				Authority: "example.com",
				Path:      "/api/v1/users/123/profile?param=value",
				Query:     "",
				Fragment:  "",
			},
			expected: "https://example.com/api/v1/users/123/profile%3Fparam%3Dvalue",
		},
		{
			name: "special characters in query",
			uri: URI{
				Scheme:    "https",
				Authority: "example.com",
				Path:      "/search",
				Query:     "q=hello world&filter=type:user&sort=date",
				Fragment:  "",
			},
			expected: "https://example.com/search?q%3Dhello%20world%26filter%3Dtype%3Auser%26sort%3Ddate",
		},
		{
			name: "special characters in fragment",
			uri: URI{
				Scheme:    "https",
				Authority: "example.com",
				Path:      "/docs",
				Query:     "",
				Fragment:  "section/subsection?param=value",
			},
			expected: "https://example.com/docs#section%2Fsubsection%3Fparam%3Dvalue",
		},
		{
			name: "file URI with Windows path",
			uri: URI{
				Scheme:    "file",
				Authority: "",
				Path:      "/C:/Users/John Doe/Documents/my file.txt",
				Query:     "",
				Fragment:  "",
			},
			expected: "file:/C%3A/Users/John%20Doe/Documents/my%20file.txt",
		},
		{
			name: "URI with percent characters that need escaping",
			uri: URI{
				Scheme:    "https",
				Authority: "example.com",
				Path:      "/discount/50%off",
				Query:     "code=SAVE%NOW",
				Fragment:  "terms%conditions",
			},
			expected: "https://example.com/discount/50%25off?code%3DSAVE%25NOW#terms%25conditions",
		},
		{
			name: "URI with hash and question mark in wrong places",
			uri: URI{
				Scheme:    "https",
				Authority: "example.com",
				Path:      "/path#hash/more?query",
				Query:     "param=value#fragment",
				Fragment:  "section?param=value",
			},
			expected: "https://example.com/path%23hash/more%3Fquery?param%3Dvalue%23fragment#section%3Fparam%3Dvalue",
		},
		{
			name: "unreserved characters not escaped",
			uri: URI{
				Scheme:    "https",
				Authority: "sub-domain.example-site.com",
				Path:      "/api_v1/users-list/~john.doe_123",
				Query:     "filter=active-users&sort=name_asc",
				Fragment:  "user-profile.details",
			},
			expected: "https://sub-domain.example-site.com/api_v1/users-list/~john.doe_123?filter%3Dactive-users%26sort%3Dname_asc#user-profile.details",
		},
		{
			name: "empty components",
			uri: URI{
				Scheme:   "custom",
				Path:     "/resource",
				Query:    "",
				Fragment: "",
			},
			expected: "custom:/resource",
		},
		{
			name: "only path with special characters",
			uri: URI{
				Path: "/path with spaces/file#name?query",
			},
			expected: "/path%20with%20spaces/file%23name%3Fquery",
		},
		{
			name: "authority with port (colon preserved)",
			uri: URI{
				Scheme:    "https",
				Authority: "localhost:3000",
				Path:      "/app",
			},
			expected: "https://localhost:3000/app",
		},
		{
			name: "mailto scheme with encoded characters",
			uri: URI{
				Scheme: "mailto",
				Path:   "user name@example.com",
				Query:  "subject=Hello World&body=Test message",
			},
			expected: "mailto:user%20name%40example.com?subject%3DHello%20World%26body%3DTest%20message",
		},
		{
			name: "authority with at sign (@ escaped)",
			uri: URI{
				Scheme:    "ftp",
				Authority: "user:password@ftp.example.com",
				Path:      "/files",
			},
			expected: "ftp://user:password%40ftp.example.com/files",
		},
		{
			name: "complex authority with multiple special chars",
			uri: URI{
				Scheme:    "http",
				Authority: "user@example.com:8080",
				Path:      "/path",
			},
			expected: "http://user%40example.com:8080/path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.uri.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestString_RoundTrip(t *testing.T) {
	// Test that URI -> StringEscaped -> ParseURI results in equivalent URI
	testCases := []URI{
		{
			Scheme:    "https",
			Authority: "user name@example server.com:8080",
			Path:      "/path with spaces/file name.txt",
			Query:     "param=value with spaces&other=test",
			Fragment:  "section with spaces",
		},
		{
			Scheme:    "file",
			Authority: "",
			Path:      "/C:/Users/John Doe/Documents/my file.txt",
			Query:     "",
			Fragment:  "",
		},
		{
			Scheme:    "custom",
			Authority: "host[::1]:3000",
			Path:      "/api/v1/test?query#fragment",
			Query:     "sort=name&filter=active",
			Fragment:  "results/page1",
		},
	}

	for i, originalURI := range testCases {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			// Convert to escaped string
			escapedString := originalURI.String()

			// Parse the escaped string back
			parsedURI, err := ParseURI(escapedString)
			assert.NoError(t, err, "Failed to parse escaped URI: %s", escapedString)

			// The parsed URI should be equivalent to the original
			assert.Equal(t, originalURI.Scheme, parsedURI.Scheme)
			assert.Equal(t, originalURI.Authority, parsedURI.Authority)
			assert.Equal(t, originalURI.Path, parsedURI.Path)
			assert.Equal(t, originalURI.Query, parsedURI.Query)
			assert.Equal(t, originalURI.Fragment, parsedURI.Fragment)
		})
	}
}

func TestParseUri_WindowsDriveHandling(t *testing.T) {
	input := "file:///c%3A/Users/Document/%23source.go"
	expected := "file:/C:/Users/Document/#source.go"
	parsed, err := ParseURI(input)
	assert.NoError(t, err)
	assert.Equal(t, expected, parsed.StringUnencoded())
}

func TestFileUri_WindowsHandling(t *testing.T) {
	input := "C:\\Users\\John Doe\\Documents\\my file.txt"
	expected := "file:/C:/Users/John Doe/Documents/my file.txt"
	uri, err := FileURI(input)
	assert.NoError(t, err)
	assert.Equal(t, expected, uri.StringUnencoded())
}

func TestFileUri_UnixHandling(t *testing.T) {
	input := "/home/user/document.txt"
	expected := "file:/home/user/document.txt"
	uri, err := FileURI(input)
	assert.NoError(t, err)
	assert.Equal(t, expected, uri.StringUnencoded())
}

func TestParseUri_DropsEmptyAuthority(t *testing.T) {
	// By the RFC for URIs, empty authorities don't really exist
	// The input "file:///path" becomes "file:/path" after parsing and reconstructing
	input := "file:///absolute/path/to/file"
	expected := "file:/absolute/path/to/file"
	parsed, err := ParseURI(input)
	assert.NoError(t, err)
	assert.Equal(t, expected, parsed.StringUnencoded())
}
