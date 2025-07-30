package handle

import (
	"bytes"
	"fmt"
	"strings"
)

func (r *Req) Parse(b []byte) *Req {
	n := len(b)

	reqLineEnd := bytes.Index(b, []byte("\r\n"))
	if reqLineEnd == -1 {
		reqLineEnd = n // no newline, treat entire chunk as request line
	}

	reqLine := string(b[:reqLineEnd])
	parts := strings.SplitN(reqLine, " ", 3)

	if len(parts) != 3 {
		return nil
	}

	theme := ""

	if strings.Contains(string(b), "theme=light") {
		theme = "light"
	}

	path := strings.SplitN(parts[1], "?", 2)[0]

	return &Req{
		Path:  strings.Trim(path, "/"),
		Ver:   parts[2],
		Theme: theme,
	}
}

func (r *Res) Bytes() []byte {
	headers := []string{
		fmt.Sprintf("HTTP/1.1 %d %s", r.Code, statuses[r.Code]),
	}

	if _, ok := r.Headers["Content-Length"]; !ok {
		r.Headers["Content-Length"] = fmt.Sprintf("%d", len(r.Data))
	}

	for k, v := range r.Headers {
		headers = append(headers, fmt.Sprintf("%s: %s", k, v))
	}

	buf := bytes.NewBufferString(strings.Join(headers, "\r\n"))
	buf.WriteString("\r\n\r\n")

	if len(r.Data) > 0 {
		buf.Write(r.Data)
	}

	return buf.Bytes()
}
