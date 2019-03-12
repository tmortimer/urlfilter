package server

import (
	"github.com/tmortimer/urlfilter/handlers"
	"testing"
)

type TestHandler struct {
	called int
}

func (h *TestHandler) Init() {
	h.called++
}

type TestServer struct {
	called int
}

func (s *TestServer) ListenAndServe() error {
	s.called++
	return nil
}

func TestRunCallsHandlersStartsServer(t *testing.T) {
	h := &TestHandler{}
	h2 := &TestHandler{}
	handlers := []handlers.Handler{
		h,
		h2,
	}

	s := &TestServer{}

	Run(handlers, s)

	if h.called != 1 {
		t.Errorf("The TestHandler Init function was called %d time(s).", h.called)
	}

	if h2.called != 1 {
		t.Errorf("The TestHandler2 Init function was called %d time(s).", h2.called)
	}

	if s.called != 1 {
		t.Errorf("The TestServer ListenAndServe function was called %d time(s).", s.called)
	}
}
