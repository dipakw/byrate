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

var uploadDownloadOptions = map[string]Opts{
	"size": {
		Default: "50",

		Opts: map[string]bool{
			"5":   true,
			"10":  true,
			"20":  true,
			"50":  true,
			"100": true,
			"200": true,
			"300": true,
			"400": true,
			"500": true,
		},
	},

	"chunk": {
		Default: "16",

		Opts: map[string]bool{
			"2":  true,
			"4":  true,
			"8":  true,
			"16": true,
			"32": true,
			"64": true,
		},
	},

	"duration": {
		Default: "10",

		Opts: map[string]bool{
			"5":  true,
			"10": true,
			"15": true,
			"20": true,
			"25": true,
			"30": true,
		},
	},
}

type Opts struct {
	Default string
	Opts    map[string]bool
}

type Req struct {
	Method  string
	Path    string
	Ver     string
	Params  map[string]string
	Headers map[string]string
	Length  int
	Index   int
	Body    []byte
}

type Res struct {
	Code    int
	Data    []byte
	Headers map[string]string
	NoSize  bool
}
