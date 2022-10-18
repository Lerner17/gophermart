package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Lerner17/gophermart/internal/models"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var (
	mockDB = map[string]models.RegisterUser{
		"user1": models.RegisterUser{"user1", "qwe1Wdasd"},
	}
	userJSON = `{"name":"Jon Snow","email":"jon@labstack.com"}`
	// userJson = json.MustMarshal(models.RegisterUser{ name: "John Snow", email: "john@labstack.com" })
)

func TestRegistration(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := Registration(db) // hi there

	if assert.NoError(t, h.getUser(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, userJSON, rec.Body.String())
	}
}
