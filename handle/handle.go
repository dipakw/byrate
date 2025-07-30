package handle

import (
	"embed"
	"fmt"
	"net"
	"strings"
	"time"
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
		defer time.Sleep(500 * time.Millisecond)

		totalSize := 50 * 1024 * 1024 // 50 MB

		res = &Res{
			Code: 200,

			Headers: map[string]string{
				"Content-Disposition": "attachment; filename=data.bin",
				"Content-Type":        "application/octet-stream",
				"Content-Length":      fmt.Sprintf("%d", totalSize),
				"Cache-Control":       "no-cache, no-store, must-revalidate",
				"Pragma":              "no-cache",
				"Expires":             "0",
				"Connection":          "close",
			},
		}

		if _, err := conn.Write(res.Bytes()); err != nil {
			return
		}

		if sendDummyBytes(conn, totalSize) != nil {
			return
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

func sendDummyBytes(conn net.Conn, dataSize int) error {
	chunkSize := 16 * 1024 // 16 KB
	buffer := make([]byte, chunkSize)
	copy(buffer, []byte("start"))
	remain := dataSize

	for remain > 0 {
		size := chunkSize

		if remain < chunkSize {
			size = remain
		}

		buffer = buffer[:size]

		if _, err := conn.Write(buffer); err != nil {
			return err
		}

		remain -= size
	}

	return nil
}
