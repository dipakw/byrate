package handle

var statuses = map[int]string{
	200: "OK",
	404: "Not Found",
}

var mimes = map[string]string{
	".html": "text/html",
	".css":  "text/css",
	".js":   "application/javascript",
	".svg":  "image/svg+xml",
}

type Req struct {
	Path string
	Ver  string
}

type Res struct {
	Code    int
	Data    []byte
	Headers map[string]string
	NoSize  bool
}
