package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/nobia/simplebank/token"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func addAuthorization(
	t *testing.T,
	request *http.Request,
	tokenManager token.TokenManager,
	authorizationType string,
	username string,
	duration time.Duration,
) {
	token, err := tokenManager.CreateToken(username, duration)
	require.NoError(t, err)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	request.Header.Set(authorizationHeaderKey, authorizationHeader)
}

func TestAuthMiddleware(t *testing.T) {
	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenManager token.TokenManager)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenManager token.TokenManager) {
				addAuthorization(t, request, tokenManager, authorizationTypeBearer, "u1", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "No authorization header",
			setupAuth: func(t *testing.T, request *http.Request, tokenManager token.TokenManager) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Unsupported authorization type",
			setupAuth: func(t *testing.T, request *http.Request, tokenManager token.TokenManager) {
				addAuthorization(t, request, tokenManager, "unsupported", "u1", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Invalid authorization format",
			setupAuth: func(t *testing.T, request *http.Request, tokenManager token.TokenManager) {
				addAuthorization(t, request, tokenManager, "", "u1", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Expired token",
			setupAuth: func(t *testing.T, request *http.Request, tokenManager token.TokenManager) {
				addAuthorization(t, request, tokenManager, authorizationTypeBearer, "u1", -time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := newTestServer(t, nil)

			authPath := "/auth"
			server.router.GET(
				authPath,
				authMiddleware(server.tokenManager),
				func(ctx *gin.Context) {
					ctx.JSON(http.StatusOK, gin.H{})
				},
			)

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenManager)
			server.router.ServeHTTP(recorder, request)

			res := recorder.Result()
			defer res.Body.Close()
			data, err := ioutil.ReadAll(res.Body)
			require.NoError(t, err)
			t.Log(string(data))

			tc.checkResponse(t, recorder)
		})
	}
}
