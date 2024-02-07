package testhelper

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/alrusov/config"
	"github.com/alrusov/db"
	"github.com/alrusov/misc"
	"github.com/alrusov/panic"
	"github.com/alrusov/rest/v4"
	"github.com/alrusov/stdhttp"
)

//----------------------------------------------------------------------------------------------------------------------------//

type (
	Data struct {
		Cfg     config.App
		Init    Init
		CfgPath string
		Calls   []*Call
	}

	Call struct {
		// In
		Method        string
		RequestHeader http.Header
		URI           string
		Body          []byte

		// Out
		ID             uint64
		ResponseHeader http.Header
		Status         int
		Answer         []byte
	}

	Init func(any) (hh *stdhttp.HTTP, err error)
)

//----------------------------------------------------------------------------------------------------------------------------//

func Rest(t *testing.T, data *Data) {
	defer Stop(0)

	panicID := panic.ID()
	defer panic.SaveStackToLogEx(panicID)

	t.Helper()

	env := misc.StringMap{}
	err := Start(t, "", data.CfgPath, env, data.Cfg)
	if err != nil {
		t.Fatal(err)
	}

	db.EnableMock()

	h, err := data.Init(data.Cfg)
	if err != nil {
		t.Fatal(err)
	}

	for i, call := range data.Calls {
		call.ID = uint64(i) + 1

		requestURL, err := url.Parse(call.URI)
		if err != nil {
			t.Fatalf("[%d] %s", call.ID, err)
		}

		_, processed := rest.Handler(h, call.ID, "", requestURL.Path, call,
			&http.Request{
				Method:           call.Method,
				URL:              requestURL,
				Proto:            "HTTP/1.1",
				ProtoMajor:       1,
				ProtoMinor:       1,
				Header:           call.RequestHeader,
				Body:             io.NopCloser(bytes.NewReader(call.Body)),
				GetBody:          nil,
				ContentLength:    int64(len(call.Body)),
				TransferEncoding: nil,
				Close:            false,
				Host:             "localhost:0000",
				Form:             nil,
				PostForm:         nil,
				MultipartForm:    nil,
				Trailer:          nil,
				RemoteAddr:       "localhost",
				RequestURI:       call.URI,
				TLS:              nil,
				Cancel:           nil,
				Response:         nil,
			},
		)

		if !processed {
			t.Fatalf("[%d] not processed", call.ID)
			return
		}

		if call.Status/100 != 2 {
			t.Fatalf("[%d] status = %d", call.ID, call.Status)
		}
	}
}

//----------------------------------------------------------------------------------------------------------------------------//

func (call *Call) Header() (h http.Header) {
	if call.ResponseHeader == nil {
		call.ResponseHeader = http.Header{}
	}

	h = call.ResponseHeader
	return
}

func (call *Call) WriteHeader(statusCode int) {
	call.Status = statusCode
}

func (call *Call) Write(data []byte) (n int, err error) {
	call.Answer = data
	return
}

//----------------------------------------------------------------------------------------------------------------------------//
