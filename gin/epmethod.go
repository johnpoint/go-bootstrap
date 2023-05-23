package gin

type GETEndpoint struct{}

func (e GETEndpoint) Method() string {
	return "GET"
}

type POSTEndpoint struct{}

func (e POSTEndpoint) Method() string {
	return "POST"
}

type PUTEndpoint struct{}

func (e PUTEndpoint) Method() string {
	return "PUT"
}

type DELETEEndpoint struct{}

func (e DELETEEndpoint) Method() string {
	return "DELETE"
}

type PATCHEndpoint struct{}

func (e PATCHEndpoint) Method() string {
	return "PATCH"
}

type HEADEndpoint struct{}

func (e HEADEndpoint) Method() string {
	return "HEAD"
}

type OPTIONSEndpoint struct{}

func (e OPTIONSEndpoint) Method() string {
	return "OPTIONS"
}

type AnyEndpoint struct{}

func (e AnyEndpoint) Method() string {
	return "Any"
}
