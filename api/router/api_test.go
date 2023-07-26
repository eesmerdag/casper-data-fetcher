package router

import (
	"cspr-fetcher/data"
	"cspr-fetcher/models"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

func TestBlocksEndpoint(t *testing.T) {
	tests := []struct {
		name            string
		req             *http.Request
		errFromDatabase bool
		validationError bool
		code            int
	}{
		{
			name: "missing limit and offset",
			req: &http.Request{
				URL: &url.URL{},
			},
			validationError: true,
			code:            http.StatusBadRequest,
		},
		{
			name: "invalid limit type",
			req: &http.Request{
				URL: &url.URL{
					RawQuery: "limit=XXX&offset=1",
				},
			},
			validationError: true,
			code:            http.StatusBadRequest,
		},
		{
			name: "invalid offset type",
			req: &http.Request{
				URL: &url.URL{
					RawQuery: "limit=1&offset=XXX",
				},
			},
			validationError: true,
			code:            http.StatusBadRequest,
		},
		{
			name: "exceed limit",
			req: &http.Request{
				URL: &url.URL{
					RawQuery: "limit=101&offset=1",
				},
			},
			validationError: true,
			code:            http.StatusBadRequest,
		},
		{
			name: "db error",
			req: &http.Request{
				URL: &url.URL{
					RawQuery: "limit=10&offset=10",
				},
			},
			errFromDatabase: true,
			code:            http.StatusInternalServerError,
		},
		{
			name: "ok",
			req: &http.Request{
				URL: &url.URL{
					RawQuery: "limit=10&offset=10",
				},
			},
			code: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockDbConnector := data.NewMockDBConnectorI(mockCtrl)

			if !tt.validationError && tt.errFromDatabase {
				mockDbConnector.EXPECT().GetAllBlocks(uint64(10), uint64(10)).Return(nil, errors.New("some db problem"))
			} else if !tt.validationError {
				mockDbConnector.EXPECT().GetAllBlocks(uint64(10), uint64(10)).Return([]*models.Block{}, nil)
			}

			router, _ := NewRouter(mockDbConnector)
			res := http.HandlerFunc(router.blocks)
			rr := httptest.NewRecorder()
			res.ServeHTTP(rr, tt.req)
			assert.True(t, rr.Code == tt.code)
		})
	}
}

func TestBlockEndpoint(t *testing.T) {
	tests := []struct {
		name            string
		req             *http.Request
		height          string
		errFromDatabase bool
		validationError bool
		notFoundBlock   bool
		ok              bool
		code            int
	}{
		{
			name: "invalid height type",
			req: &http.Request{
				URL: &url.URL{},
			},
			height:          "XXX",
			validationError: true,
			code:            http.StatusBadRequest,
		},
		{
			name: "db error",
			req: &http.Request{
				URL: &url.URL{},
			},
			height:          "100",
			errFromDatabase: true,
			code:            http.StatusInternalServerError,
		},
		{
			name: "not found block",
			req: &http.Request{
				URL: &url.URL{},
			},
			height:        "100",
			notFoundBlock: true,
			code:          http.StatusNotFound,
		},
		{
			name: "ok",
			req: &http.Request{
				URL: &url.URL{},
			},
			height: "100",
			ok:     true,
			code:   http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockDbConnector := data.NewMockDBConnectorI(mockCtrl)

			tt.req = mux.SetURLVars(tt.req, map[string]string{
				"height": tt.height,
			})

			if tt.errFromDatabase {
				h, _ := strconv.ParseInt(tt.height, 10, 64)

				mockDbConnector.EXPECT().GetBlockByHeight(uint64(h)).Return(nil, errors.New("some db problem"))
			} else if tt.notFoundBlock {
				h, _ := strconv.ParseInt(tt.height, 10, 64)

				mockDbConnector.EXPECT().GetBlockByHeight(uint64(h)).Return(nil, nil)
			} else if tt.ok {
				h, _ := strconv.ParseInt(tt.height, 10, 64)

				mockDbConnector.EXPECT().GetBlockByHeight(uint64(h)).Return(&models.Block{}, nil)
			}

			router, _ := NewRouter(mockDbConnector)
			res := http.HandlerFunc(router.block)
			rr := httptest.NewRecorder()
			res.ServeHTTP(rr, tt.req)
			assert.True(t, rr.Code == tt.code)
		})
	}
}
