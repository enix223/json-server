package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/enix223/json-server/serializer"
)

type Server struct {
	internal http.Server
	stop     chan struct{}
	stopped  chan struct{}
	// Handler mapping
	// /foo
	//   - GET: handler1
	//   - POST: handler2
	// /bar
	//   - GET: handler3
	muxMap      map[string]map[string]http.HandlerFunc
	mapping     []EndpointMapping
	serializers map[PayloadType]serializer.Serializer
}

type ServerOptions struct {
	ListenAddress string
	MappingJson   string
}

func NewServer(opts ServerOptions) *Server {
	Assert(len(opts.ListenAddress) > 0, "Server address is required")
	Assert(len(opts.MappingJson) > 0, "Handler mapping file path is required")
	content, err := os.ReadFile(opts.MappingJson)
	if err != nil {
		log.Panic("failed to read handler mapping file", err)
	}

	mapping := make([]EndpointMapping, 0)
	err = json.Unmarshal(content, &mapping)
	if err != nil {
		panic(err)
	}

	s := Server{
		stop:     make(chan struct{}),
		stopped:  make(chan struct{}),
		internal: http.Server{},
		mapping:  mapping,
		muxMap:   make(map[string]map[string]http.HandlerFunc),
		serializers: map[PayloadType]serializer.Serializer{
			"json": serializer.NewJsonSerializer(),
			"text": serializer.NewTextSerializer(),
		},
	}
	s.internal.Addr = opts.ListenAddress
	mux := http.NewServeMux()
	for _, entry := range mapping {
		s.registerHandler(mux, entry)
	}
	s.internal.Handler = mux
	return &s
}

func (s *Server) Run() {
	go s.internalRun()
	go func() {
		<-s.stop
		log.Println("JSON Server gracefully stopped")
		s.internal.Close()
		s.stopped <- struct{}{}
	}()
}

func (s *Server) Stop() {
	s.stop <- struct{}{}
	<-s.stopped
}

func (s *Server) internalRun() {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
		}
	}()
	log.Printf("JSON server listening %s", s.internal.Addr)
	s.internal.ListenAndServe()
}

func (s *Server) registerHandler(mux *http.ServeMux, entry EndpointMapping) {
	if entry.Method == nil {
		method := http.MethodGet
		entry.Method = &method
	}

	hm, ok := s.muxMap[entry.Path]
	if ok {
		// already register outter handler
		_, ok := hm[*entry.Method]
		if ok {
			log.Panicf("handler alread register: %s %s", *entry.Method, entry.Path)
		}
		hm[*entry.Method] = s.createHandler(entry)
		return
	}

	hm = make(map[string]http.HandlerFunc)
	hm[*entry.Method] = s.createHandler(entry)
	s.muxMap[entry.Path] = hm

	outterHandler := func(w http.ResponseWriter, r *http.Request) {
		hm := s.muxMap[entry.Path]
		h, ok := hm[r.Method]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		h(w, r)
	}
	mux.HandleFunc(entry.Path, outterHandler)
}

func (s *Server) createHandler(entry EndpointMapping) http.HandlerFunc {
	log.Printf("Register handler for path: %s %s", *entry.Method, entry.Path)
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("REQUEST: %s %s", r.Method, entry.Path)
		if entry.ResponseHeaders != nil {
			for key, value := range *entry.ResponseHeaders {
				w.Header().Set(key, value)
			}
		}

		if entry.StatusCode != nil {
			w.WriteHeader(*entry.StatusCode)
		}

		if entry.Payload != nil {
			ser, ok := s.serializers[entry.PayloadType]
			if !ok {
				w.Write([]byte(`unsupported payload type`))
				return
			}

			payload, err := ser.Serialize(*entry.Payload)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`failed to serialize payload`))
				return
			}
			w.Write(payload)
		}
	}
}
