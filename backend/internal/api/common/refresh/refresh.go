package refresh

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const RefreshToken = "refreshToken"

func ExtractTokenFromCookie(c echo.Context) *uuid.UUID {
	cookie, err := c.Cookie(RefreshToken)
	if err != nil {
		return nil
	}

	token, err := uuid.Parse(cookie.Value)
	if err != nil {
		return nil
	}

	return &token
}
