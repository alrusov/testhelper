package testhelper

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"

	"github.com/alrusov/stdhttp"
)

//----------------------------------------------------------------------------------------------------------------------------//

type (
	liveTest struct {
		proc LiveTestProcessor
	}

	LiveTestProcessor func(id uint64, prefix string, path string, w http.ResponseWriter, r *http.Request)
)

var (
	lt *liveTest
)

//----------------------------------------------------------------------------------------------------------------------------//

func SetLiveTest(h *stdhttp.HTTP, proc LiveTestProcessor) (err error) {
	if proc == nil {
		return
	}

	if lt != nil {
		err = fmt.Errorf("live test already set to %s", runtime.FuncForPC(reflect.ValueOf(lt.proc).Pointer()).Name())
		return
	}

	lt = &liveTest{proc: proc}
	h.AddHandler(lt, false)
	return
}

//----------------------------------------------------------------------------------------------------------------------------//

func (lt *liveTest) Handler(id uint64, prefix string, path string, w http.ResponseWriter, r *http.Request) (processed bool) {
	switch path {
	default:

	case "/live-test":
		processed = true
		lt.proc(id, prefix, path, w, r)
	}

	return
}

//----------------------------------------------------------------------------------------------------------------------------//
