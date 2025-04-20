package get_users_id

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/moevm/nosql1h25-writer/backend/internal/api"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/common/decorator"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/common/mw"
	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
	"github.com/moevm/nosql1h25-writer/backend/internal/service/users"
)

type Profile struct {
	Rating      float64   `json:"rating"`
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Request struct {
	ID       primitive.ObjectID `param:"id" validate:"required"`
	Profiles []string           `query:"profile" validate:"omitempty,dive,oneof=client freelancer"`
}

type Response struct {
	ID          primitive.ObjectID    `json:"id"`
	DisplayName string                `json:"displayName"`
	Email       string                `json:"email"`
	SystemRole  entity.SystemRoleType `json:"systemRole"`
	CreatedAt   time.Time             `json:"createdAt"`
	UpdatedAt   time.Time             `json:"updatedAt"`
	Balance     int                 `json:"balance"`
	Client      *Profile              `json:"client,omitempty"`
	Freelancer  *Profile              `json:"freelancer,omitempty"`
}

type handler struct {
	usersService users.Service
}

func New(usersService users.Service) api.Handler {
	return decorator.NewBindAndValidate(&handler{usersService: usersService})
}

// Handle - Get user by ID handler
//
//	@Summary		Get user by ID
//	@Description	Retrieves user details by their ObjectID. Requires authentication. Access restricted to the user themselves or administrators. Optionally filters profiles.
//	@Tags			Users
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string			true	"User ObjectID (Hex)"									example(507f1f77bcf86cd799439011)
//	@Param			profile	query		[]string		false	"Profile types to include ('client', 'freelancer')."	collectionFormat(multi)	Enums(client, freelancer)
//	@Success		200		{object}	Response		"Successfully retrieved user details"
//	@Failure		400		{object}	echo.HTTPError	"Invalid request format or ObjectID"
//	@Failure		401		{object}	echo.HTTPError	"Unauthorized (invalid or missing JWT)"
//	@Failure		403		{object}	echo.HTTPError	"Forbidden (access denied)"
//	@Failure		404		{object}	echo.HTTPError	"User not found"
//	@Failure		500		{object}	echo.HTTPError	"Internal server error"
//	@Router			/users/{id} [get]
func (h *handler) Handle(c echo.Context, inp Request) error {
	authUserID, okUserID := c.Get(mw.UserIDKey).(primitive.ObjectID)
	authUserRole, okUserRole := c.Get(mw.SystemRoleKey).(entity.SystemRoleType)

	if !okUserID || !okUserRole {
		log.Error("Handler: Failed to get user ID or role from context")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to process authentication context")
	}

	if authUserRole != entity.SystemRoleTypeAdmin && authUserID != inp.ID {
		log.Warnf("Handler: Access denied for user %s requesting user %s", authUserID.Hex(), inp.ID.Hex())
		return echo.NewHTTPError(http.StatusForbidden, "Access denied: You can only view your own profile or require admin privileges.")
	}

	ctx := c.Request().Context()
	user, err := h.usersService.FindUserByID(ctx, inp.ID)

	if err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			log.Infof("Handler: User not found: %s", inp.ID.Hex())
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		log.WithError(err).Errorf("Handler: Failed to find user by ID: %s", inp.ID.Hex())
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve user data")
	}

	response := Response{
        ID:          user.ID,
        DisplayName: user.DisplayName,
        Email:       user.Email,
        SystemRole:  entity.SystemRoleType(user.SystemRole),
        CreatedAt:   user.CreatedAt,
        UpdatedAt:   user.UpdatedAt,
        Balance:     user.Balance,
    }

    includeClient := lo.Contains(inp.Profiles, "client")
    includeFreelancer := lo.Contains(inp.Profiles, "freelancer")

    if includeClient {
        response.Client = &Profile{
            Rating:      user.Client.Rating,
            Description: user.Client.Description,
            UpdatedAt:   user.Client.UpdatedAt,
        }
    }

    if includeFreelancer {
        response.Freelancer = &Profile{
            Rating:      user.Freelancer.Rating,
            Description: user.Freelancer.Description,
            UpdatedAt:   user.Freelancer.UpdatedAt,
        }
    }

    return c.JSON(http.StatusOK, response)
}
