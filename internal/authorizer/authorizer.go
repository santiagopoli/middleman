package authorizer

type Request struct {
	Host    string
	Method  string
	Path    string
	Headers map[string][]string
}

type Authorizer interface {
	IsAuthorized(*Request) bool
}
