package handle

import (
	"bufio"
	"bytes"
	"embed"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

//go:embed ui/index.html ui/app.js ui/style.css ui/favicon.svg
var ui embed.FS

func Handle(conn net.Conn) {
	defer conn.Close()

	writer := bufio.NewWriter(conn)

	buf := make([]byte, 4096)

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
		opts := map[string]int{}

		for key, opt := range uploadDownloadOptions {
			if val, ok := req.Params[key]; ok {
				if opt.Opts[val] {
					opts[key], _ = strconv.Atoi(val)
				}
			}
		}

		for key, opt := range uploadDownloadOptions {
			if _, ok := opts[key]; !ok {
				opts[key], _ = strconv.Atoi(opt.Default)
			}
		}

		totalSize := opts["size"] * 1024 * 1024

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

		conn.SetWriteDeadline(time.Now().Add(time.Duration(opts["duration"]+2) * time.Second))

		if _, err := writer.Write(res.Bytes()); err != nil {
			return
		}

		if err := writer.Flush(); err != nil {
			return
		}

		if sendDummyBytes(writer, totalSize, opts["chunk"]) != nil {
			return
		}

		return
	}

	if req.Path == "upload" && req.Method == "POST" {
		opts := map[string]int{}

		for key, opt := range uploadDownloadOptions {
			if val, ok := req.Params[key]; ok {
				if opt.Opts[val] {
					opts[key], _ = strconv.Atoi(val)
				}
			}
		}

		for key, opt := range uploadDownloadOptions {
			if _, ok := opts[key]; !ok {
				opts[key], _ = strconv.Atoi(opt.Default)
			}
		}

		res := &Res{
			Code: 200,
			Data: []byte("done"),
		}

		if _, ok := req.Headers["content-length"]; !ok {
			return
		}

		contentLength, err := strconv.Atoi(req.Headers["content-length"])

		if err != nil {
			return
		}

		if contentLength > MAX_UPLOAD_SIZE {
			return
		}

		if contentLength != opts["size"]*1024*1024 {
			return
		}

		remainSizeToRead := contentLength - len(req.Body)

		if remainSizeToRead > 0 {
			conn.SetReadDeadline(time.Now().Add(time.Duration(opts["duration"]+2) * time.Second))

			buf := make([]byte, opts["chunk"]*1024)
			read := 0

			for read < remainSizeToRead {
				n, err := conn.Read(buf)

				if err != nil {
					return
				}

				read += n

				if n == 0 {
					break
				}
			}
		}

		writer.Write(res.Bytes())
		writer.Flush()

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

		if filename == "index.html" {
			setTheme := ""

			if req.Headers["cookie"] != "" && strings.Contains(req.Headers["cookie"], "theme=light") {
				setTheme = " light"
			}

			file = bytes.Replace(file, []byte(" curentTheme"), []byte(setTheme), -1)
		}

		res = &Res{
			Code: 200,
			Data: file,

			Headers: map[string]string{
				"Content-Type": mimeType,
			},
		}
	}

	writer.Write(res.Bytes())
	writer.Flush()
}

func sendDummyBytes(writer *bufio.Writer, dataSize int, chunkKb int) error {
	chunkSize := chunkKb * 1024
	buffer := make([]byte, chunkSize)
	copy(buffer, []byte("start"))
	remain := dataSize

	for remain > 0 {
		size := chunkSize

		if remain < chunkSize {
			size = remain
		}

		buffer = buffer[:size]

		if _, err := writer.Write(buffer); err != nil {
			return err
		}

		remain -= size
	}

	return writer.Flush()
}
