package router

import (
	"cspr-fetcher/data"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type ErrorResp struct {
	Message string
	Code    int
}

type Router struct {
	router *mux.Router
	db     data.DBConnectorI
}

func NewRouter(db data.DBConnectorI) (*Router, error) {
	router := mux.NewRouter()

	r := &Router{
		router: router,
		db:     db,
	}

	router.HandleFunc("/blocks", r.blocks).Methods(http.MethodGet)
	router.HandleFunc("/blocks/{height}", r.block).Methods(http.MethodGet)
	router.Use(panicRecovery)
	return r, nil
}

func panicRecovery(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				errorResponse(w, "unexpected internal error", http.StatusInternalServerError)
				return
			}
		}()

		h.ServeHTTP(w, r)
	})
}

func (rt Router) blocks(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	limit := params["limit"]
	offset := params["offset"]

	if limit == nil || offset == nil {
		errorResponse(w, "missing limit or offset info", http.StatusBadRequest)
		return
	}

	l, err := strconv.ParseInt(limit[0], 10, 64)
	if err != nil {
		errorResponse(w, "limit should be integer", http.StatusBadRequest)
		return
	}
	if l > 100 {
		errorResponse(w, "max limit is 100", http.StatusBadRequest)
		return
	}

	o, err := strconv.ParseInt(offset[0], 10, 64)
	if err != nil {
		errorResponse(w, "offset should be integer", http.StatusBadRequest)
		return
	}

	blocks, err := rt.db.GetAllBlocks(uint64(l), uint64(o))
	if err != nil {
		errorResponse(w, "unexpected internal error on getting blocks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(blocks)
}

func (rt Router) block(w http.ResponseWriter, r *http.Request) {
	height := mux.Vars(r)["height"]
	h, err := strconv.ParseInt(height, 10, 64)
	if err != nil {
		errorResponse(w, "block height is not integer", http.StatusBadRequest)
		return
	}

	block, err := rt.db.GetBlockByHeight(uint64(h))
	if err != nil {
		errorResponse(w, "unexpected internal error on getting blocks", http.StatusInternalServerError)
		return
	}

	if block == nil {
		errorResponse(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(block)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}

func errorResponse(w http.ResponseWriter, message string, code int) {
	errObj := ErrorResp{
		Message: message,
		Code:    code,
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(errObj)
}
