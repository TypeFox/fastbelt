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
	}{
		{
			name:  "complete HTTP URL",
			input: "http://example.com:8080/path/to/resource?query=value&foo=bar#section",
			expected: &uri{
				scheme:    "http",
				authority: "example.com:8080",
				path:      "/path/to/resource",
				query:     "query=value&foo=bar",
				fragment:  "section",
			},
		},
		{
			name:  "HTTPS URL without port",
			input: "https://www.example.com/api/v1/users",
			expected: &uri{
				scheme:    "https",
				authority: "www.example.com",
				path:      "/api/v1/users",
				query:     "",
				fragment:  "",
			},
		},
		{
			name:  "URL with query but no fragment",
			input: "https://search.example.com/search?q=golang&type=code",
			expected: &uri{
				scheme:    "https",
				authority: "search.example.com",
				path:      "/search",
				query:     "q=golang&type=code",
				fragment:  "",
			},
		},
		{
			name:  "URL with fragment but no query",
			input: "https://docs.example.com/guide#installation",
			expected: &uri{
				scheme:    "https",
				authority: "docs.example.com",
				path:      "/guide",
				query:     "",
				fragment:  "installation",
			},
		},
		{
			name:  "file URI with three slashes",
			input: "file:///home/user/document.txt",
			expected: &uri{
				scheme:    "file",
				authority: "",
				path:      "/home/user/document.txt",
				query:     "",
				fragment:  "",
			},
		},
		{
			name:  "relative path only",
			input: "/path/to/resource",
			expected: &uri{
				scheme:    "",
				authority: "",
				path:      "/path/to/resource",
				query:     "",
				fragment:  "",
			},
		},
		{
			name:  "path with query and fragment",
			input: "/api/users?active=true#list",
			expected: &uri{
				scheme:    "",
				authority: "",
				path:      "/api/users",
				query:     "active=true",
				fragment:  "list",
			},
		},
		{
			name:  "scheme only",
			input: "custom:",
			expected: &uri{
				scheme:    "custom",
				authority: "",
				path:      "",
				query:     "",
				fragment:  "",
			},
		},
		{
			name:  "authority with userinfo",
			input: "ftp://user:pass@ftp.example.com/files/",
			expected: &uri{
				scheme:    "ftp",
				authority: "user:pass@ftp.example.com",
				path:      "/files/",
				query:     "",
				fragment:  "",
			},
		},
		{
			name:  "empty string",
			input: "",
			expected: &uri{
				scheme:    "",
				authority: "",
				path:      "",
				query:     "",
				fragment:  "",
			},
		},
		{
			name:  "complex query string",
			input: "https://api.example.com/v2/search?q=golang+testing&lang=en&sort=date&page=1",
			expected: &uri{
				scheme:    "https",
				authority: "api.example.com",
				path:      "/v2/search",
				query:     "q=golang+testing&lang=en&sort=date&page=1",
				fragment:  "",
			},
		},
		{
			name:  "URI with port and special chars in fragment",
			input: "http://localhost:3000/app#section/subsection?param=value",
			expected: &uri{
				scheme:    "http",
				authority: "localhost:3000",
				path:      "/app",
				query:     "",
				fragment:  "section/subsection?param=value",
			},
		},
		{
			name:  "mailto scheme",
			input: "mailto:user@example.com?subject=Hello&body=Hi%20there",
			expected: &uri{
				scheme:    "mailto",
				authority: "",
				path:      "user@example.com",
				query:     "subject=Hello&body=Hi there",
				fragment:  "",
			},
		},
		{
			name:  "encoded characters in path (vscode URI)",
			input: "file:///c%3A/Users/Document/%23source.go",
			expected: &uri{
				scheme:    "file",
				authority: "",
				path:      "/C:/Users/Document/#source.go",
				query:     "",
				fragment:  "",
			},
		},
		{
			name:  "windows file URI",
			input: "file:///C:/Users/John%20Doe/Documents/my%20file.txt",
			expected: &uri{
				scheme:    "file",
				authority: "",
				path:      "/C:/Users/John Doe/Documents/my file.txt",
				query:     "",
				fragment:  "",
			},
		},
		{
			name:  "multiple question marks",
			input: "http://example.com/path?param1=value?param2=value2",
			expected: &uri{
				scheme:    "http",
				authority: "example.com",
				path:      "/path",
				query:     "param1=value?param2=value2",
				fragment:  "",
			},
		},
		{
			name:  "multiple hash symbols",
			input: "http://example.com/path#fragment#more",
			expected: &uri{
				scheme:    "http",
				authority: "example.com",
				path:      "/path",
				query:     "",
				fragment:  "fragment#more",
			},
		},
		{
			name:  "query without path",
			input: "http://example.com?param=value",
			expected: &uri{
				scheme:    "http",
				authority: "example.com",
				path:      "",
				query:     "param=value",
				fragment:  "",
			},
		},
		{
			name:  "fragment without path or query",
			input: "http://example.com#section",
			expected: &uri{
				scheme:    "http",
				authority: "example.com",
				path:      "",
				query:     "",
				fragment:  "section",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseURI(tt.input)

			assert.Equal(t, tt.expected.Scheme(), result.Scheme())
			assert.Equal(t, tt.expected.Authority(), result.Authority())
			assert.Equal(t, tt.expected.Path(), result.Path())
			assert.Equal(t, tt.expected.Query(), result.Query())
			assert.Equal(t, tt.expected.Fragment(), result.Fragment())
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
			uri: &uri{
				scheme:    "https",
				authority: "example.com:443",
				path:      "/path/to/resource",
				query:     "param=value",
				fragment:  "section",
			},
			expected: "https://example.com:443/path/to/resource?param=value#section",
		},
		{
			name: "URI without authority",
			uri: &uri{
				scheme:   "file",
				path:     "/home/user/file.txt",
				query:    "",
				fragment: "",
			},
			expected: "file:///home/user/file.txt",
		},
		{
			name: "URI with empty authority",
			uri: &uri{
				scheme:    "file",
				authority: "",
				path:      "/home/user/file.txt",
				query:     "",
				fragment:  "",
			},
			expected: "file:///home/user/file.txt",
		},
		{
			name: "path only",
			uri: &uri{
				path: "/api/data",
			},
			expected: "/api/data",
		},
		{
			name: "query and fragment only",
			uri: &uri{
				query:    "search=golang",
				fragment: "results",
			},
			expected: "?search=golang#results",
		},
		{
			name:     "empty URI",
			uri:      &uri{},
			expected: "",
		},
		{
			name: "scheme and path only",
			uri: &uri{
				scheme: "custom",
				path:   "/resource",
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
			parsed := ParseURI(original)

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
		{"file:///absolute/path/to/file", "file:///absolute/path/to/file"},
		{"file:///C:/Windows/file.txt", "file:///C:/Windows/file.txt"},
		{"file://server/share/file.txt", "file://server/share/file.txt"}, // non-empty authority preserved
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			parsed := ParseURI(tt.input)

			reconstructed := parsed.StringUnencoded()
			assert.Equal(t, tt.expected, reconstructed)
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
			uri: &uri{
				scheme:    "https",
				authority: "example.com",
				path:      "/api/users",
				query:     "active=true",
				fragment:  "section",
			},
			expected: "https://example.com/api/users?active%3Dtrue#section",
		},
		{
			name: "spaces in different components",
			uri: &uri{
				scheme:    "https",
				authority: "example server.com",
				path:      "/path with spaces/file",
				query:     "param=value with spaces",
				fragment:  "section with spaces",
			},
			expected: "https://example%20server.com/path%20with%20spaces/file?param%3Dvalue%20with%20spaces#section%20with%20spaces",
		},
		{
			name: "special characters in authority",
			uri: &uri{
				scheme:    "http",
				authority: "user:pass@host[::1]:8080",
				path:      "/path",
				query:     "",
				fragment:  "",
			},
			expected: "http://user:pass%40host[::1]:8080/path",
		},
		{
			name: "IPv6 address in authority (preserved)",
			uri: &uri{
				scheme:    "http",
				authority: "[::1]:8080",
				path:      "/path",
				query:     "",
				fragment:  "",
			},
			expected: "http://[::1]:8080/path",
		},
		{
			name: "path with forward slashes preserved",
			uri: &uri{
				scheme:    "https",
				authority: "example.com",
				path:      "/api/v1/users/123/profile?param=value",
				query:     "",
				fragment:  "",
			},
			expected: "https://example.com/api/v1/users/123/profile%3Fparam%3Dvalue",
		},
		{
			name: "special characters in query",
			uri: &uri{
				scheme:    "https",
				authority: "example.com",
				path:      "/search",
				query:     "q=hello world&filter=type:user&sort=date",
				fragment:  "",
			},
			expected: "https://example.com/search?q%3Dhello%20world%26filter%3Dtype%3Auser%26sort%3Ddate",
		},
		{
			name: "special characters in fragment",
			uri: &uri{
				scheme:    "https",
				authority: "example.com",
				path:      "/docs",
				query:     "",
				fragment:  "section/subsection?param=value",
			},
			expected: "https://example.com/docs#section%2Fsubsection%3Fparam%3Dvalue",
		},
		{
			name: "file URI with Windows path",
			uri: &uri{
				scheme:    "file",
				authority: "",
				path:      "/C:/Users/John Doe/Documents/my file.txt",
				query:     "",
				fragment:  "",
			},
			expected: "file:///C%3A/Users/John%20Doe/Documents/my%20file.txt",
		},
		{
			name: "URI with percent characters that need escaping",
			uri: &uri{
				scheme:    "https",
				authority: "example.com",
				path:      "/discount/50%off",
				query:     "code=SAVE%NOW",
				fragment:  "terms%conditions",
			},
			expected: "https://example.com/discount/50%25off?code%3DSAVE%25NOW#terms%25conditions",
		},
		{
			name: "URI with hash and question mark in wrong places",
			uri: &uri{
				scheme:    "https",
				authority: "example.com",
				path:      "/path#hash/more?query",
				query:     "param=value#fragment",
				fragment:  "section?param=value",
			},
			expected: "https://example.com/path%23hash/more%3Fquery?param%3Dvalue%23fragment#section%3Fparam%3Dvalue",
		},
		{
			name: "unreserved characters not escaped",
			uri: &uri{
				scheme:    "https",
				authority: "sub-domain.example-site.com",
				path:      "/api_v1/users-list/~john.doe_123",
				query:     "filter=active-users&sort=name_asc",
				fragment:  "user-profile.details",
			},
			expected: "https://sub-domain.example-site.com/api_v1/users-list/~john.doe_123?filter%3Dactive-users%26sort%3Dname_asc#user-profile.details",
		},
		{
			name: "empty components",
			uri: &uri{
				scheme:   "custom",
				path:     "/resource",
				query:    "",
				fragment: "",
			},
			expected: "custom:/resource",
		},
		{
			name: "only path with special characters",
			uri: &uri{
				path: "/path with spaces/file#name?query",
			},
			expected: "/path%20with%20spaces/file%23name%3Fquery",
		},
		{
			name: "authority with port (colon preserved)",
			uri: &uri{
				scheme:    "https",
				authority: "localhost:3000",
				path:      "/app",
			},
			expected: "https://localhost:3000/app",
		},
		{
			name: "mailto scheme with encoded characters",
			uri: &uri{
				scheme: "mailto",
				path:   "user name@example.com",
				query:  "subject=Hello World&body=Test message",
			},
			expected: "mailto:user%20name%40example.com?subject%3DHello%20World%26body%3DTest%20message",
		},
		{
			name: "authority with at sign (@ escaped)",
			uri: &uri{
				scheme:    "ftp",
				authority: "user:password@ftp.example.com",
				path:      "/files",
			},
			expected: "ftp://user:password%40ftp.example.com/files",
		},
		{
			name: "complex authority with multiple special chars",
			uri: &uri{
				scheme:    "http",
				authority: "user@example.com:8080",
				path:      "/path",
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
	testCases := []*uri{
		{
			scheme:    "https",
			authority: "user name@example server.com:8080",
			path:      "/path with spaces/file name.txt",
			query:     "param=value with spaces&other=test",
			fragment:  "section with spaces",
		},
		{
			scheme:    "file",
			authority: "",
			path:      "/C:/Users/John Doe/Documents/my file.txt",
			query:     "",
			fragment:  "",
		},
		{
			scheme:    "custom",
			authority: "host[::1]:3000",
			path:      "/api/v1/test?query#fragment",
			query:     "sort=name&filter=active",
			fragment:  "results/page1",
		},
	}

	for i, originalURI := range testCases {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			// Convert to escaped string
			escapedString := originalURI.String()

			// Parse the escaped string back
			parsedURI := ParseURI(escapedString)

			// The parsed URI should be equivalent to the original
			assert.Equal(t, originalURI.scheme, parsedURI.Scheme())
			assert.Equal(t, originalURI.authority, parsedURI.Authority())
			assert.Equal(t, originalURI.path, parsedURI.Path())
			assert.Equal(t, originalURI.query, parsedURI.Query())
			assert.Equal(t, originalURI.fragment, parsedURI.Fragment())
		})
	}
}

func TestParseUri_WindowsDriveHandling(t *testing.T) {
	input := "file:///c%3A/Users/Document/%23source.go"
	expected := "file:///C:/Users/Document/#source.go"
	parsed := ParseURI(input)
	assert.Equal(t, expected, parsed.StringUnencoded())
}

func TestFileUri_WindowsHandling(t *testing.T) {
	input := "C:\\Users\\John Doe\\Documents\\my file.txt"
	expected := "file:///C:/Users/John Doe/Documents/my file.txt"
	uri := FileURI(input)
	assert.Equal(t, expected, uri.StringUnencoded())
}

func TestFileUri_UnixHandling(t *testing.T) {
	input := "/home/user/document.txt"
	expected := "file:///home/user/document.txt"
	uri := FileURI(input)
	assert.Equal(t, expected, uri.StringUnencoded())
}

func TestParseUri_DropsEmptyAuthority(t *testing.T) {
	// By the RFC for URIs, empty authorities don't really exist
	// The input "file:///path" becomes "file:/path" after parsing and reconstructing
	input := "file:///absolute/path/to/file"
	expected := "file:///absolute/path/to/file"
	parsed := ParseURI(input)
	assert.Equal(t, expected, parsed.StringUnencoded())
}

func TestParseUri_LowerCaseSchemeAndAuthority(t *testing.T) {
	// URI RFC specifies that the scheme and authority are case-insensitive and should be normalized to lowercase.
	input := "HTTP://EXAMPLE.COM/PATH"
	expected := "http://example.com/PATH"
	parsed := ParseURI(input)
	assert.Equal(t, expected, parsed.StringUnencoded())
}
