package fastbelt

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/TypeFox/go-lsp/protocol"
)

type URI struct {
	Scheme    string
	Authority string
	Path      string
	Query     string
	Fragment  string
}

// StringUnencoded returns the URI as a string without percent-encoding. This is useful for debugging and logging purposes.
func (u URI) StringUnencoded() string {
	var result strings.Builder

	// Add scheme
	if u.Scheme != "" {
		result.WriteString(u.Scheme)
		result.WriteString(":")
	}

	// Add authority
	if u.Authority != "" {
		result.WriteString("//")
		result.WriteString(u.Authority)
	}

	// Add path
	result.WriteString(u.Path)

	// Add query
	if u.Query != "" {
		result.WriteString("?")
		result.WriteString(u.Query)
	}

	// Add fragment
	if u.Fragment != "" {
		result.WriteString("#")
		result.WriteString(u.Fragment)
	}

	return result.String()
}

// String returns the URI as a string with percent-encoding applied to the components as needed.
func (u URI) String() string {
	var result strings.Builder

	// Add scheme
	if u.Scheme != "" {
		// Scheme should not be escaped, because it cannot contain reserved characters
		result.WriteString(u.Scheme)
		result.WriteString(":")
	}

	// Add authority
	if u.Authority != "" {
		result.WriteString("//")
		result.WriteString(encodeURIComponent(u.Authority, false, true))
	}

	// Add path
	result.WriteString(encodeURIComponent(u.Path, true, false))

	// Add query
	if u.Query != "" {
		result.WriteString("?")
		result.WriteString(encodeURIComponent(u.Query, false, false))
	}

	// Add fragment
	if u.Fragment != "" {
		result.WriteString("#")
		result.WriteString(encodeURIComponent(u.Fragment, false, false))
	}

	return result.String()
}

func (u URI) DocumentURI() protocol.DocumentURI {
	return protocol.DocumentURI(u.String())
}

func NewURI(scheme, authority, path, query, fragment string) (URI, error) {
	if !isValidScheme(scheme) {
		return URI{}, fmt.Errorf("invalid scheme: %s", scheme)
	}
	return URI{
		Scheme:    strings.ToLower(scheme),
		Authority: authority,
		Path:      path,
		Query:     query,
		Fragment:  fragment,
	}, nil
}

func (u URI) With(scheme, authority, path, query, fragment *string) URI {
	if scheme != nil {
		u.Scheme = strings.ToLower(*scheme)
	}
	if authority != nil {
		u.Authority = *authority
	}
	if path != nil {
		u.Path = *path
	}
	if query != nil {
		u.Query = *query
	}
	if fragment != nil {
		u.Fragment = *fragment
	}
	return u
}

var parseRegexp = regexp.MustCompile(`^(([^:/?#]+?):)?(\/\/([^/?#]*))?([^?#]*)(\?([^#]*))?(#(.*))?`)

func ParseURI(uri string) (URI, error) {
	var result URI
	matches := parseRegexp.FindStringSubmatch(uri)
	if matches == nil {
		return result, fmt.Errorf("invalid URI: %s", uri)
	}

	var err error

	// Decode components that might contain percent-encoded characters
	result.Scheme = strings.ToLower(matches[2]) // Schemes are not URL-encoded, but always lowercased

	if matches[4] != "" {
		result.Authority, err = url.PathUnescape(matches[4])
		if err != nil {
			// If decoding fails, use original value
			result.Authority = matches[4]
		}
	}

	if matches[5] != "" {
		result.Path, err = url.PathUnescape(matches[5])
		if err != nil {
			// If decoding fails, use original value
			result.Path = matches[5]
		}
		result.Path = normalizeDriveLetter(result.Path)
	}

	if matches[7] != "" {
		result.Query, err = url.PathUnescape(matches[7])
		if err != nil {
			// If decoding fails, use original value
			result.Query = matches[7]
		}
	}

	if matches[9] != "" {
		result.Fragment, err = url.PathUnescape(matches[9])
		if err != nil {
			// If decoding fails, use original value
			result.Fragment = matches[9]
		}
	}

	return result, nil
}

func MustParseURI(uri string) URI {
	parsed, err := ParseURI(uri)
	if err != nil {
		panic(fmt.Sprintf("failed to parse URI '%s': %s", uri, err))
	}
	return parsed
}

func FileURI(path string) (URI, error) {
	// Normalize Windows paths by replacing backslashes with forward slashes
	normalized := strings.ReplaceAll(path, "\\", "/")
	// Ensure the path starts with a slash
	if !strings.HasPrefix(normalized, "/") {
		normalized = "/" + normalized
	}
	// Further normalize the drive letter
	normalized = normalizeDriveLetter(normalized)
	return NewURI("file", "", normalized, "", "")
}

// Uppercases the drive letter in Windows file paths to ensure consistent URIs across platforms
func normalizeDriveLetter(path string) string {
	offset := 0
	if strings.HasPrefix(path, "/") {
		offset = 1
	}
	if len(path) >= 3+offset && path[2+offset] == '/' && path[1+offset] == ':' && (path[0+offset] >= 'a' && path[0+offset] <= 'z') {
		return path[0:offset] + strings.ToUpper(string(path[offset:1+offset])) + path[1+offset:]
	}
	return path
}

// isValidScheme checks if the scheme contains only valid characters
func isValidScheme(scheme string) bool {
	if len(scheme) == 0 {
		return false
	}

	// First character must be a letter
	if !isLetter(scheme[0]) {
		return false
	}

	// Rest can be letters, digits, '+', '-', '.'
	for i := 1; i < len(scheme); i++ {
		c := scheme[i]
		if !isLetter(c) && !isDigit(c) && c != '+' && c != '-' && c != '.' {
			return false
		}
	}

	return true
}

const upperhex = "0123456789ABCDEF"

// Modified version of the one in net/url to preserve unescaped slashes in paths
// Also escapes ' ' (whitespace) to '%20' instead of '+'
func encodeURIComponent(component string, isPath, isAuthority bool) string {
	hexCount := 0
	for i := 0; i < len(component); i++ {
		c := component[i]
		if shouldEscape(c, isPath, isAuthority) {
			hexCount++
		}
	}

	// Nothing to escape
	if hexCount == 0 {
		return component
	}

	var buf [64]byte
	var t []byte

	required := len(component) + 2*hexCount
	if required <= len(buf) {
		t = buf[:required]
	} else {
		t = make([]byte, required)
	}

	j := 0
	for i := 0; i < len(component); i++ {
		switch c := component[i]; {
		case shouldEscape(c, isPath, isAuthority):
			t[j] = '%'
			t[j+1] = upperhex[c>>4]
			t[j+2] = upperhex[c&15]
			j += 3
		default:
			t[j] = component[i]
			j++
		}
	}
	return string(t)
}

func shouldEscape(c byte, isPath, isAuthority bool) bool {
	// unreserved characters: https://tools.ietf.org/html/rfc3986#section-2.3
	shouldKeep := isLetter(c) ||
		isDigit(c) ||
		c == '-' ||
		c == '_' ||
		c == '.' ||
		c == '~' ||
		(isPath && c == '/') ||
		(isAuthority && c == '[') ||
		(isAuthority && c == ']') ||
		(isAuthority && c == ':')
	return !shouldKeep
}

// isLetter checks if a byte is an ASCII letter
func isLetter(c byte) bool {
	return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')
}

// isDigit checks if a byte is an ASCII digit
func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}
