package fastbelt

import (
	"net/url"
	"regexp"
	"strings"

	"typefox.dev/lsp"
)

const FileScheme = "file"

type URI interface {
	Scheme() string
	Authority() string
	Path() string
	Query() string
	Fragment() string
	// Returns the URI as a string with percent-encoding applied to the components as needed.
	// This is the standard way to serialize URIs and should be used when interoperability with other tools is required.
	String() string
	// Returns the URI as a string without percent-encoding.
	// This is useful for debugging and logging purposes.
	StringUnencoded() string
	// Converts the URI to a lsp.DocumentURI, which is the format used by the LSP library.
	DocumentURI() lsp.DocumentURI
	// Returns a new URI with the specified scheme, keeping other components unchanged.
	WithScheme(scheme string) URI
	// Returns a new URI with the specified authority, keeping other components unchanged.
	WithAuthority(authority string) URI
	// Returns a new URI with the specified path, keeping other components unchanged.
	WithPath(path string) URI
	// Returns a new URI with the specified query, keeping other components unchanged.
	WithQuery(query string) URI
	// Returns a new URI with the specified fragment, keeping other components unchanged.
	WithFragment(fragment string) URI
	// Returns a new URI with the specified components, keeping other components unchanged.
	// The components are pointers, so you can pass nil for components that should remain unchanged.
	//
	// Note that this method will not validate or normalize the components,
	// so it's the caller's responsibility to ensure that the resulting URI is valid.
	With(scheme, authority, path, query, fragment *string) URI
	// Checks if this [URI] is equal to another [URI], based on their components.
	// Note that this comparison is case-sensitive for all components.
	// Use [FileURI] and [ParseURI] to ensure consistent normalization for reliable comparisons.
	Equal(other URI) bool
}

type uri struct {
	scheme    string
	authority string
	path      string
	query     string
	fragment  string
	encoded   *string
	unencoded *string
}

func (u *uri) Scheme() string {
	return u.scheme
}

func (u *uri) WithScheme(scheme string) URI {
	return u.With(&scheme, nil, nil, nil, nil)
}

func (u *uri) Authority() string {
	return u.authority
}

func (u *uri) WithAuthority(authority string) URI {
	return u.With(nil, &authority, nil, nil, nil)
}

func (u *uri) Path() string {
	return u.path
}

func (u *uri) WithPath(path string) URI {
	return u.With(nil, nil, &path, nil, nil)
}

func (u *uri) Query() string {
	return u.query
}

func (u *uri) WithQuery(query string) URI {
	return u.With(nil, nil, nil, &query, nil)
}

func (u *uri) Fragment() string {
	return u.fragment
}

func (u *uri) WithFragment(fragment string) URI {
	return u.With(nil, nil, nil, nil, &fragment)
}

func (u *uri) Equal(other URI) bool {
	if u == nil && other == nil {
		return true
	}
	if u == nil || other == nil {
		return false
	}
	return u.Scheme() == other.Scheme() &&
		u.Authority() == other.Authority() &&
		u.Path() == other.Path() &&
		u.Query() == other.Query() &&
		u.Fragment() == other.Fragment()
}

func (u *uri) StringUnencoded() string {
	if u.unencoded != nil {
		return *u.unencoded
	}
	var result strings.Builder

	// Add scheme
	if u.scheme != "" {
		result.WriteString(u.scheme)
		result.WriteString(":")
	}

	if u.authority != "" || u.scheme == FileScheme {
		// Add authority prefix if authority is present or if it's a file URI (which can have an empty authority)
		result.WriteString("//")
	}

	// Add authority
	if u.authority != "" {
		result.WriteString(u.authority)
	}

	// Add path
	result.WriteString(u.path)

	// Add query
	if u.query != "" {
		result.WriteString("?")
		result.WriteString(u.query)
	}

	// Add fragment
	if u.fragment != "" {
		result.WriteString("#")
		result.WriteString(u.fragment)
	}

	value := result.String()
	u.unencoded = &value
	return value
}

func (u *uri) String() string {
	if u.encoded != nil {
		return *u.encoded
	}
	var result strings.Builder

	// Add scheme
	if u.scheme != "" {
		// Scheme should not be escaped, because it cannot contain reserved characters
		result.WriteString(u.scheme)
		result.WriteString(":")
	}

	if u.authority != "" || u.scheme == FileScheme {
		// Add authority prefix if authority is present or if it's a file URI (which can have an empty authority)
		result.WriteString("//")
	}

	// Add authority
	if u.authority != "" {
		result.WriteString(encodeURIComponent(u.authority, false, true))
	}

	// Add path
	result.WriteString(encodeURIComponent(u.path, true, false))

	// Add query
	if u.query != "" {
		result.WriteString("?")
		result.WriteString(encodeURIComponent(u.query, false, false))
	}

	// Add fragment
	if u.fragment != "" {
		result.WriteString("#")
		result.WriteString(encodeURIComponent(u.fragment, false, false))
	}

	value := result.String()
	u.encoded = &value
	return value
}

func (u *uri) DocumentURI() lsp.DocumentURI {
	return lsp.DocumentURI(u.String())
}

// Constructs a new URI from the given components.
// The components are not validated or normalized, so it's the caller's responsibility to ensure that they are valid.
//
// Instead of creating a new URI from scratch, it is recommended to parse URIs from a string using [ParseURI].
func NewURI(scheme, authority, path, query, fragment string) URI {
	return &uri{
		scheme:    scheme,
		authority: authority,
		path:      path,
		query:     query,
		fragment:  fragment,
	}
}

func (u *uri) With(scheme, authority, path, query, fragment *string) URI {
	result := uri{}
	if u != nil {
		// Copy all existing components
		result = uri{
			scheme:    u.scheme,
			authority: u.authority,
			path:      u.path,
			query:     u.query,
			fragment:  u.fragment,
		}
	}
	if scheme != nil {
		result.scheme = *scheme
	}
	if authority != nil {
		result.authority = *authority
	}
	if path != nil {
		result.path = *path
	}
	if query != nil {
		result.query = *query
	}
	if fragment != nil {
		result.fragment = *fragment
	}
	return &result
}

var parseRegexp = regexp.MustCompile(`^(([^:/?#]+?):)?(\/\/([^/?#]*))?([^?#]*)(\?([^#]*))?(#(.*))?`)

// Parses a URI from a string. The parsing is lenient and will not fail for invalid URIs, but it will try to extract the components as best as possible.
func ParseURI(value string) URI {
	result := uri{}
	matches := parseRegexp.FindStringSubmatch(value)
	if matches == nil {
		// Generally only an empty string would fail to parse
		// But this is still a (mostly) valid URI
		return &result
	}
	var err error
	// Decode components that might contain percent-encoded characters
	result.scheme = strings.ToLower(matches[2]) // Schemes are not URL-encoded, but always lowercased

	if matches[4] != "" {
		result.authority, err = url.PathUnescape(matches[4])
		if err != nil {
			// If decoding fails, use original value
			result.authority = matches[4]
		}
		// Authorities are case-insensitive, so we normalize to lowercase
		result.authority = strings.ToLower(result.authority)
	}

	if matches[5] != "" {
		result.path, err = url.PathUnescape(matches[5])
		if err != nil {
			result.path = matches[5]
		}
		result.path = normalizeDriveLetter(result.path)
	}

	if matches[7] != "" {
		result.query, err = url.PathUnescape(matches[7])
		if err != nil {
			result.query = matches[7]
		}
	}

	if matches[9] != "" {
		result.fragment, err = url.PathUnescape(matches[9])
		if err != nil {
			result.fragment = matches[9]
		}
	}

	return &result
}

func FileURI(path string) URI {
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
