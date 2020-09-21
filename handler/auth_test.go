package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/qiniu/x/xlog"
	"github.com/qrtc/qlive/protocol"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticate(t *testing.T) {
	handler := &AuthHandler{
		Auth: &MockAuth{},
	}
	testCases := []struct {
		authHeader         string
		expectedStatusCode int
		expectedUserID     string
	}{
		{
			authHeader:         "",
			expectedStatusCode: 401,
		},
		{
			authHeader:         "user-1#login-token",
			expectedStatusCode: 401,
		},
		{
			authHeader:         "Bearer user-1login-token",
			expectedStatusCode: 401,
		},
		{
			authHeader:         "Bearer user-1#login-token",
			expectedStatusCode: 0,
			expectedUserID:     "user-1",
		},
	}

	for i, testCase := range testCases {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(protocol.XLogKey, xlog.New(fmt.Sprintf("test-authenticate-%d", i)))

		req, err := http.NewRequest("POST", "/v1/profile", nil)
		assert.Nilf(t, err, "failed to build request for case %d", i)
		req.Header.Set("Authorization", testCase.authHeader)
		c.Request = req

		handler.Authenticate(c)
		if testCase.expectedStatusCode != 0 {
			assert.Equalf(t, testCase.expectedStatusCode, w.Code, "code is not the same as expected for test case %d", i)
		}
		if testCase.expectedUserID != "" {
			assert.Equalf(t, testCase.expectedUserID, c.GetString(protocol.UserIDContextKey), "user ID is not the same as expected in case %d", i)
		}
	}
}
