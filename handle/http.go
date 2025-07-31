package handle

import (
	"bytes"
	"fmt"
	"strings"
)

func (r *Req) Parse(b []byte) *Req {
	// Find index separating headers and body
	headerEnd := bytes.Index(b, []byte("\r\n\r\n"))

	if headerEnd == -1 {
		return nil // Malformed request, missing header-body separator
	}

	r.Index = headerEnd

	// Split headers from body
	headerSection := b[:headerEnd]
	lines := bytes.Split(headerSection, []byte("\r\n"))

	if len(lines) < 1 {
		return nil // Malformed request
	}

	// Parse request line: METHOD PATH VERSION
	reqLine := string(lines[0])
	parts := strings.SplitN(reqLine, " ", 3)
	if len(parts) != 3 {
		return nil // Malformed request line
	}

	method := strings.ToUpper(parts[0])
	path := strings.Trim(parts[1], "/")
	version := parts[2]

	// Parse query parameters
	params := map[string]string{}
	pathParts := strings.SplitN(path, "?", 2)
	if len(pathParts) > 1 {
		path = pathParts[0]
		query := pathParts[1]
		for _, pair := range strings.Split(query, "&") {
			if kv := strings.SplitN(pair, "=", 2); len(kv) == 2 {
				params[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
			}
		}
	}

	// Parse headers
	headers := map[string]string{}
	for _, line := range lines[1:] {
		if kv := bytes.SplitN(line, []byte(":"), 2); len(kv) == 2 {
			key := strings.TrimSpace(string(kv[0]))
			val := strings.TrimSpace(string(kv[1]))
			headers[strings.ToLower(key)] = val
		}
	}

	return &Req{
		Method:  method,
		Path:    path,
		Ver:     version,
		Params:  params,
		Headers: headers,
		Index:   headerEnd + 4,
		Body:    b[headerEnd+4:],
	}
}

func (r *Res) Bytes() []byte {
	headers := []string{
		fmt.Sprintf("HTTP/1.1 %d %s", r.Code, statuses[r.Code]),
	}

	if r.Headers != nil {
		if _, ok := r.Headers["Content-Length"]; !ok {
			r.Headers["Content-Length"] = fmt.Sprintf("%d", len(r.Data))
		}
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
