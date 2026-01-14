package benchmark

func URLPattern() string {
	return `https?://(www\.)?[-a-zA-Z0-9@:%._\+~#=]+\.[a-zA-Z0-9()]+([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`
}

func EmailPattern() string {
	return `[\w\.+-]+@[\w\.-]+\.[\w\.-]+`
}

func IPv4Pattern() string {
	return `((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`
}
