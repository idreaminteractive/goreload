package hotreload

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Handler struct {
	sse    *SSEHandler
	writer *io.Writer
}

func newHandler() (h *Handler) {
	h = &Handler{
		sse:    NewSSEHandler(),
		writer: &io.Discard,
	}
	return h
}

func (p *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// w.Header().Add("Access-Control-Allow-Origin", "*")
	if r.URL.Path == "/hotreload" {
		switch r.Method {
		case http.MethodGet:
			// Provides a list of messages including a reload message.
			p.sse.ServeHTTP(w, r)
			return
		case http.MethodPost:
			fmt.Fprint(*p.writer, "Reload triggered\n")
			// Send a reload message to all connected clients.
			p.sse.Send("message", "reload")
			return
		}
		http.Error(w, "only GET or POST method allowed", http.StatusMethodNotAllowed)
		return
	}
	http.NotFound(w, r)

}

func (p *Handler) SendSSE(eventType string, data string) {
	p.sse.Send(eventType, data)
}

type HotReloadServer struct {
	server    *http.Server
	logwriter *io.Writer
}

// this is such a neat pattern.
func WithLogger(w *io.Writer) func(*HotReloadServer) {
	return func(hrs *HotReloadServer) {
		hrs.logwriter = w
		// also add logger to the handler
		if h, ok := hrs.server.Handler.(*Handler); ok {
			h.writer = w
		}

	}
}

func InitHotReloadServer(port int, options ...func(*HotReloadServer)) *HotReloadServer {

	// set a default logger
	h := enableCORS(newHandler())
	hrs := &HotReloadServer{
		logwriter: &io.Discard,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: h,
		},
	}
	for _, o := range options {
		o(hrs)
	}

	return hrs
}

func (s *HotReloadServer) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	return s.server.Shutdown(ctx)
}

func (s *HotReloadServer) Run(ctx context.Context) (err error) {
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()
	return err
}

func ValidateUrl(u string) (string, int, error) {
	parsed, err := url.Parse(u)
	if err != nil {
		return "", 0, err
	}
	portNum := 8082
	if parsed.Port() != "" {
		portNum, err = strconv.Atoi(parsed.Port())
		if err != nil {
			return "", 0, err
		}
	}

	fullPath := parsed.Scheme + "://" + parsed.Host

	return fullPath, portNum, nil
}
