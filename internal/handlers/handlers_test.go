package handlers

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"github.com/EvgeniyBudaev/yandex-go-musthave-diploma-tpl/internal/auth"
	"github.com/EvgeniyBudaev/yandex-go-musthave-diploma-tpl/internal/config"
	"github.com/EvgeniyBudaev/yandex-go-musthave-diploma-tpl/internal/mocks"
	"github.com/EvgeniyBudaev/yandex-go-musthave-diploma-tpl/internal/storage"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type wantResponse struct {
	code            int
	headerContent   string
	responseContent string
}

func TestValidateOrder(t *testing.T) {
	tests := []struct {
		name        string
		order       string
		resultOrder uint
		errCode     int
	}{
		{
			name:        "not integer order",
			order:       "SomeText",
			resultOrder: 0,
			errCode:     http.StatusBadRequest,
		},
		{
			name:        "negative order",
			order:       "-5",
			resultOrder: 0,
			errCode:     http.StatusBadRequest,
		},
		{
			name:        "correct order with odd digit quantity",
			order:       "101",
			resultOrder: 101,
			errCode:     http.StatusOK,
		},
		{
			name:        "incorrect order with odd digit quantity",
			order:       "124",
			resultOrder: 0,
			errCode:     http.StatusUnprocessableEntity,
		},
		{
			name:        "correct order with even digit quantity",
			order:       "4953",
			resultOrder: 4953,
			errCode:     http.StatusOK,
		},
		{
			name:        "incorrect order with even digit quantity",
			order:       "3743",
			resultOrder: 0,
			errCode:     http.StatusUnprocessableEntity,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, errCode, _ := ValidateOrder(tt.order)
			assert.Equal(t, tt.resultOrder, result)
			assert.Equal(t, tt.errCode, errCode)
		})
	}
}

func TestRegisterHandler(t *testing.T) {
	tt := []struct {
		name                string
		want                wantResponse
		registerData        storage.Auth
		mockResponseID      string
		mockResponseErrCode int
	}{
		{
			"success_register",
			wantResponse{
				http.StatusOK,
				"",
				``,
			},
			storage.Auth{Login: "TestLogin", Password: "TestPassword"},
			"5f0319ee-bd23-40a2-9b56-f6d21726c425",
			http.StatusOK,
		},
		{
			"fail_register_same_login",
			wantResponse{
				http.StatusFailedDependency,
				"text/plain; charset=utf-8",
				"Could not register user\n",
			},
			storage.Auth{Login: "TestLogin", Password: "TestPassword"},
			"",
			http.StatusFailedDependency,
		},
		{
			"fail_register_internal_error",
			wantResponse{
				http.StatusInternalServerError,
				"text/plain; charset=utf-8",
				"Could not register user\n",
			},
			storage.Auth{Login: "TestLogin", Password: "TestPassword"},
			"",
			http.StatusInternalServerError,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			marshalledData, _ := json.Marshal(tc.registerData)
			request := httptest.NewRequest(http.MethodGet, "/api/user/register", bytes.NewBuffer(marshalledData))
			w := httptest.NewRecorder()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			storage := mocks.NewMockStorage(ctrl)
			storage.EXPECT().Register(ctx, tc.registerData).Return(tc.mockResponseID, tc.mockResponseErrCode)
			handler := http.HandlerFunc(GetHandlerWithStorage(storage).Register)
			handler.ServeHTTP(w, request)
			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, tc.want.code, result.StatusCode)
			if result.StatusCode == http.StatusOK {
				cookies := result.Cookies()
				for _, cookie := range cookies {
					if cookie.Name == config.GetUserCookie() {
						h := auth.GenerateCookie()
						h.Write([]byte(tc.mockResponseID))
						sign := h.Sum(nil)
						assert.Equal(t, hex.EncodeToString(append([]byte(tc.mockResponseID)[:], sign[:]...)), cookie.Value)
						break
					}
					assert.Fail(t, "get no cookies for UserID")
				}
			}
			responseBody, err := io.ReadAll(result.Body)
			assert.Nil(t, err)
			assert.Equal(t, tc.want.headerContent, result.Header.Get("Content-Type"))
			assert.Equal(t, tc.want.responseContent, string(responseBody))
		})
	}
}
