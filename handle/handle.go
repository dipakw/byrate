package handle

import (
	"embed"
	"fmt"
	"net"
	"strings"
)

//go:embed ui/index.html ui/app.js ui/style.css
var ui embed.FS

func Handle(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 64)

	n, err := conn.Read(buf)

	if err != nil {
		return
	}

	req := (&Req{}).Parse(buf[:n])

	if req == nil {
		return
	}

	res := &Res{
		Code: 404,
		Data: []byte("Not Found"),

		Headers: map[string]string{
			"Content-Type": "text/plain",
		},
	}

	if req.Path == "download" {
		chunkSize := 4 * 1024                // 4 KB
		totalSize := 10 * 1024 * 1024 * 1024 // 10 GB
		buffer := make([]byte, chunkSize)

		res = &Res{
			Code: 200,

			Headers: map[string]string{
				"Content-Disposition": "attachment; filename=data.bin",
				"Content-Type":        "application/octet-stream",
				"Content-Length":      fmt.Sprintf("%d", totalSize),
				"Cache-Control":       "no-cache, no-store, must-revalidate",
				"Pragma":              "no-cache",
				"Expires":             "0",
			},
		}

		if _, err := conn.Write(res.Bytes()); err != nil {
			return
		}

		for i := 0; i < totalSize; i += chunkSize {
			buffer = buffer[:chunkSize]

			if _, err := conn.Write(buffer); err != nil {
				return
			}
		}

		return
	}

	filename := req.Path

	if filename == "" {
		filename = "index.html"
	}

	file, err := ui.ReadFile("ui/" + filename)

	if err == nil {
		mimeType := "text/plain"

		for ext, mime := range mimes {
			if strings.HasSuffix(filename, ext) {
				mimeType = mime
				break
			}
		}

		res = &Res{
			Code: 200,
			Data: file,

			Headers: map[string]string{
				"Content-Type": mimeType,
			},
		}
	}

	conn.Write(res.Bytes())
}
