package communication

import (
	"encoding/json"
	"fmt"
	"github.com/21stio/go-playground/gifts/pkg/business"
	"github.com/21stio/go-playground/gifts/pkg/types"
	"net/http"
)

const (
	AssignPath = "/assign"
)

type Server struct {
	facade      *business.Facade
	addr        string
	listener    *http.Server
	logErrFunc  func(...interface{})
	logInfoFunc func(...interface{})
}

func NewServer(facade *business.Facade, addr string, logErrFunc func(...interface{}), logInfoFunc func(...interface{})) (*Server) {
	return &Server{
		facade:      facade,
		addr:        addr,
		logErrFunc:  logErrFunc,
		logInfoFunc: logInfoFunc,
	}
}

func (s *Server) Serve() (err error) {
	if s.listener != nil {
		panic("already listening")
	}

	mux := http.NewServeMux()

	mux.HandleFunc(AssignPath, func(w http.ResponseWriter, r *http.Request) {
		rsp := types.AssignEmployeeToGiftResponse{}

		err := func() (err error) {
			req := types.AssignEmployeeToGiftRequest{}
			err = json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				return
			}

			err = s.facade.AssignEmployeeToGift(req.EmployeeName)
			if err != nil {
				return
			}

			return
		}()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			rsp.Error = err.Error()
		}

		err = json.NewEncoder(w).Encode(rsp)
		if err != nil {
			s.logErrFunc(err)
		}
	})

	s.logInfoFunc(fmt.Sprintf("listening on addr: %v", s.addr))

	s.listener = &http.Server{Addr: s.addr, Handler: mux}
	err = s.listener.ListenAndServe()
	if err != nil {
		return
	}

	return
}

func (s *Server) Close() (err error) {
	if s.listener == nil {
		return
	}

	err = s.listener.Close()
	if err != nil {
	    return
	}

	s.listener = nil

	return
}
