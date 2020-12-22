package postfile

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

// 12/22/2020, gonyi, wrote it for slack file upload.
//    unlike i thought, it seems like i cannot upload multiple fils to slack at once..?
//    maybe a bug on my end..

type postMultipartForm struct {
	buf *bytes.Buffer
	w   *multipart.Writer
	h   map[string][]string
}

func New() *postMultipartForm {
	m := postMultipartForm{
		buf: &bytes.Buffer{},
		h:   make(map[string][]string),
	}
	m.w = multipart.NewWriter(m.buf)
	return &m
}

func (m *postMultipartForm) AddField(key, value string) {
	m.w.WriteField(key, value)
}

func (m *postMultipartForm) AddFile(filename string) (n int64, err error) {
	if f, err := m.w.CreateFormFile("file", filename); err != nil {
		n = 0
	} else {
		file, err := os.Open(filename)
		if err != nil {
			return 0, err
		}
		n, err = io.Copy(f, file)
		file.Close()
	}
	return
}

func (m *postMultipartForm) AddFileReader(filename string, ior io.Reader) (n int64, err error) {
	if f, err := m.w.CreateFormFile("file", filename); err != nil {
		n = 0
	} else {
		n, err = io.Copy(f, ior)
	}
	return
}

func (m *postMultipartForm) Reader() io.Reader {
	return m.buf
}

func (m *postMultipartForm) ContentType() string {
	return m.w.FormDataContentType()
}

func (m *postMultipartForm) SetHTTPHeader(key, value string) {
	m.h[key] = []string{value}
}

func (m *postMultipartForm) AddHTTPHeader(key, value string) {
	if v, ok := m.h[key]; ok {
		m.h[key] = append(v, value)
	} else {
		m.SetHTTPHeader(key, value)
	}
}

func (m *postMultipartForm) GetRequest(method, url string) (*http.Request, error) {
	if err := m.w.Close(); err != nil {
		return nil, err
	}
	return http.NewRequest(method, url, m.buf)
}

func (m *postMultipartForm) Send(method, url string, resp *bytes.Buffer) (received int64, err error) {
	req, err := m.GetRequest(method, url)
	if err != nil {
		return 0, err
	}
	// content type will be automatically prepared by multipart writer
	// and add authroization keyee if any

	req.Header.Set("Content-Type", m.ContentType())
	for k, v := range m.h {
		if len(v) > 1 {
			for _, v2 := range v {
				req.Header.Add(k, v2)
			}
		} else {
			req.Header.Set(k, v[0])
		}
	}

	// send the request
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}

	// out, err := ioutil.ReadAll(r.Body)
	if resp != nil {
		received, err = io.Copy(resp, r.Body)
	}

	if r != nil && r.Body != nil {
		r.Body.Close()
	}

	return received, err
}
