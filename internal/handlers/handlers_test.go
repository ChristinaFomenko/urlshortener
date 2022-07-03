package handlers

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mock "github.com/ChristinaFomenko/shortener/internal/handlers/mocks"
)

func TestShortenHandler(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		shortcut    string
	}
	tests := []struct {
		name     string
		request  string
		url      string
		shortcut string
		want     want
		err      error
	}{
		{
			name:     "success",
			url:      "https://yandex.ru",
			shortcut: "http://localhost:8080/abcde",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  201,
				shortcut:    "http://localhost:8080/abcde",
			},
			request: "/",
			err:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serviceMock := mock.NewMockservice(ctrl)
			serviceMock.EXPECT().Shorten(tt.url).Return(tt.shortcut, nil)

			httpHandler := New(serviceMock, nil)

			buffer := new(bytes.Buffer)
			buffer.WriteString(tt.url)
			request := httptest.NewRequest(http.MethodPost, tt.request, buffer)

			writer := httptest.NewRecorder()
			HandlerFunc := http.HandlerFunc(httpHandler.Shorten)
			HandlerFunc.ServeHTTP(writer, request)
			result := writer.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			bodyResult, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			require.NoError(t, err)

			assert.Equal(t, tt.want.shortcut, string(bodyResult))
		})
	}
}

func TestAPIJSONShorten_Success(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		response    string
	}
	tests := []struct {
		name     string
		request  string
		url      string
		body     string
		shortcut string
		want     want
		err      error
	}{
		{
			name:     "success",
			url:      "https://yandex.ru",
			body:     "{\"url\":\"https://yandex.ru\"}",
			shortcut: "http://localhost:8080/abcde",
			want: want{
				contentType: "application/json",
				statusCode:  201,
				response:    "{\"result\":\"http://localhost:8080/abcde\"}",
			},
			request: "/api/shorten",
			err:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serviceMock := mock.NewMockservice(ctrl)
			serviceMock.EXPECT().Shorten(tt.url).Return(tt.shortcut, tt.err)

			httpHandler := New(serviceMock, nil)

			buffer := new(bytes.Buffer)
			buffer.WriteString(tt.body)
			request := httptest.NewRequest(http.MethodPost, tt.request, buffer)

			writer := httptest.NewRecorder()
			handlerFunc := http.HandlerFunc(httpHandler.APIJSONShorten)
			handlerFunc.ServeHTTP(writer, request)
			result := writer.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			bodyResult, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			require.NoError(t, err)

			assert.Equal(t, tt.want.response, string(bodyResult))
		})
	}
}

func TestAPIJSONShorten_BadRequest(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		response    string
	}
	tests := []struct {
		name     string
		request  string
		url      string
		body     string
		shortcut string
		want     want
	}{
		{
			name:     "bad-request",
			url:      "https://yandex.ru",
			body:     "{\"url\":\"\"}",
			shortcut: "http://localhost:8080/abcde",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
				response:    "request in not valid\n",
			},
			request: "/api/shorten",
		},
		{
			name:     "bad-request",
			url:      "https://yandex.ru",
			body:     "{\"url\":\"qwerty\"}",
			shortcut: "http://localhost:8080/abcde",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
				response:    "request in not valid\n",
			},
			request: "/api/shorten",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			httpHandler := New(nil, nil)

			buffer := new(bytes.Buffer)
			buffer.WriteString(tt.body)
			request := httptest.NewRequest(http.MethodPost, tt.request, buffer)

			writer := httptest.NewRecorder()
			handlerFunc := http.HandlerFunc(httpHandler.APIJSONShorten)
			handlerFunc.ServeHTTP(writer, request)
			result := writer.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			bodyResult, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			require.NoError(t, err)

			assert.Equal(t, tt.want.response, string(bodyResult))
		})
	}
}

func Test_handler_GetByUserID(t *testing.T) {
	type fields struct {
		service service
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	var tests []struct {
		name   string
		fields fields
		args   args
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &handler{
				service: tt.fields.service,
			}
			h.GetByUserID(tt.args.w, tt.args.r)
		})
	}
}
