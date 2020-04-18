package unittests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/internal/controllers"
	"github.com/jz222/loggy/internal/mocks"
)

func TestRegisterError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	loggingControllers := controllers.LoggingControllers{
		LoggingService: &mocks.LoggingService{},
	}

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	r.POST("/", loggingControllers.RegisterError)

	c.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte("{}")))

	r.ServeHTTP(w, c.Request)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGetUserInformation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userControllers := controllers.UserControllers{
		UserService: &mocks.UserService{},
	}

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	r.Use(mocks.VerifyUserJWT)
	r.GET("/", userControllers.Get)

	c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)

	r.ServeHTTP(w, c.Request)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}
