package authorizer

type Request struct {
	Method  string
	Path    string
	Headers map[string][]string
}

type Authorizer interface {
	IsAuthorized(*Request) bool
}
