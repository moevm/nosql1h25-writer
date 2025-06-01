package get_admin_users

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/moevm/nosql1h25-writer/backend/internal/api"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/common/decorator"
	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
	"github.com/moevm/nosql1h25-writer/backend/internal/service/users"
)

type handler struct {
	usersService users.Service
}

func New(usersService users.Service) api.Handler {
	return decorator.NewBindAndValidate(&handler{usersService: usersService})
}

type Request struct {
	Offset              *int       `query:"offset" validate:"omitempty,gte=0" example:"0"`
	Limit               *int       `query:"limit" validate:"omitempty,gte=1,lte=200" example:"10"`
	NameSearch          *string    `query:"nameSearch" validate:"omitempty" example:"Joh"`
	EmailSearch         *string    `query:"emailSearch" validate:"omitempty" example:"mail.ru"`
	Roles               []string   `query:"role" validate:"omitempty,dive,oneof=user admin"`
	MinFreelancerRating *float64   `query:"minFreelancerRating" validate:"omitempty,gte=0,lte=5"`
	MaxFreelancerRating *float64   `query:"maxFreelancerRating" validate:"omitempty,gte=0,lte=5"`
	MinClientRating     *float64   `query:"minClientRating" validate:"omitempty,gte=0,lte=5"`
	MaxClientRating     *float64   `query:"maxClientRating" validate:"omitempty,gte=0,lte=5"`
	MinCreatedAt        *time.Time `query:"minCreatedAt" validate:"omitempty"`
	MaxCreatedAt        *time.Time `query:"maxCreatedAt" validate:"omitempty"`
	MaxBalance          *int       `query:"maxBalance" validate:"omitempty,gte=0"`
	MinBalance          *int       `query:"minBalance" validate:"omitempty,gte=0"`
	SortBy              *string    `query:"sortBy" validate:"omitempty,oneof=newest oldest rich poor name_asc name_desc freelancer_rating_asc freelancer_rating_desc client_rating_asc client_rating_desc" example:"newest"`
}

type User struct {
	ID               primitive.ObjectID    `json:"id" validate:"required" example:"582ebf010936ac3ba5cd00e4"`
	DisplayName      string                `json:"displayName" validate:"required" example:"John Doe"`
	Email            string                `json:"email" validate:"required" example:"goida@mail.ru"`
	SystemRole       entity.SystemRoleType `json:"systemRole" validate:"required" example:"admin"`
	Balance          int                   `json:"balance" validate:"required" example:"500"`
	CreatedAt        time.Time             `json:"createdAt" validate:"required" example:"2020-01-01T00:00:00Z"`
	UpdatedAt        time.Time             `json:"updatedAt" validate:"required" example:"2020-01-01T00:00:00Z"`
	FreelancerRating float64               `json:"freelancerRating" validate:"required" example:"4.8"`
	ClientRating     float64               `json:"clientRating" validate:"required" example:"4.7"`
}

type Response struct {
	Users []User `json:"users" validate:"required"`
	Total int    `json:"total" example:"250"`
}

// Handle - Return user list handler
//
//	@Summary		Return user list
//	@Description	Return user list
//	@Tags			admin
//	@Security		JWT
//	@Param			offset	query	int		false	"Offset"	default(0)	minimum(0)	example(0)
//	@Param			limit	query	int		false	"Limit"		default(10)	minimum(1)	maximum(200)	example(10)
//	@Param			request	body	Request	true	"fields in query"
//	@Produce		json
//	@Success		200	{object}	Response
//	@Failure		401	{object}	echo.HTTPError
//	@Failure		403	{object}	echo.HTTPError
//	@Failure		500	{object}	echo.HTTPError
//	@Router			/admin/users [get]
func (h *handler) Handle(c echo.Context, in Request) error {
	offset, limit := h.applyDefaults(in)

	out, err := h.usersService.Find(c.Request().Context(), users.FindIn{
		Limit:               limit,
		Offset:              offset,
		NameSearch:          in.NameSearch,
		EmailSearch:         in.EmailSearch,
		Roles:               in.Roles,
		MinFreelancerRating: in.MinFreelancerRating,
		MaxFreelancerRating: in.MaxFreelancerRating,
		MinClientRating:     in.MinClientRating,
		MaxClientRating:     in.MaxClientRating,
		MinCreatedAt:        in.MinCreatedAt,
		MaxCreatedAt:        in.MaxCreatedAt,
		MinBalance:          in.MinBalance,
		MaxBalance:          in.MaxBalance,
		SortBy:              in.SortBy,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	users := make([]User, 0, len(out.Users))
	for _, user := range out.Users {
		users = append(users, User{
			ID:               user.ID,
			DisplayName:      user.DisplayName,
			Email:            user.Email,
			SystemRole:       user.SystemRole,
			Balance:          user.Balance,
			CreatedAt:        user.CreatedAt,
			UpdatedAt:        user.UpdatedAt,
			FreelancerRating: user.Freelancer.Rating,
			ClientRating:     user.Client.Rating,
		})
	}

	return c.JSON(http.StatusOK, Response{Users: users, Total: out.Total})
}

func (h *handler) applyDefaults(in Request) (offset int, limit int) {
	if in.Offset != nil {
		offset = *in.Offset
	} else {
		offset = 0
	}

	if in.Limit != nil {
		limit = *in.Limit
	} else {
		limit = 10
	}

	return
}
