package app

import (
	"github.com/moevm/nosql1h25-writer/backend/internal/api"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/get_admin"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/get_admin_export"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/get_health"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/get_orders"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/get_orders_id"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/get_users_id"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/patch_orders_id"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/patch_users_id"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/post_auth_login"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/post_auth_logout"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/post_auth_refresh"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/post_auth_register"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/post_balance_deposit"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/post_balance_withdraw"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/post_orders"
)

func (app *App) GetHealthHandler() api.Handler {
	if app.getHealthHandler != nil {
		return app.getHealthHandler
	}

	app.getHealthHandler = get_health.New(app.OrdersCollection())
	return app.getHealthHandler
}

func (app *App) GetUsersIDHandler() api.Handler {
	return get_users_id.New(app.UsersService())
}

func (app *App) PostAuthLoginHandler() api.Handler {
	if app.postAuthLoginHandler != nil {
		return app.postAuthLoginHandler
	}

	app.postAuthLoginHandler = post_auth_login.New(app.AuthService())
	return app.postAuthLoginHandler
}

func (app *App) PostAuthRefreshHandler() api.Handler {
	if app.postAuthRefreshHandler != nil {
		return app.postAuthRefreshHandler
	}

	app.postAuthRefreshHandler = post_auth_refresh.New(app.AuthService())
	return app.postAuthRefreshHandler
}

func (app *App) PostAuthLogoutHandler() api.Handler {
	if app.postAuthLogoutHandler != nil {
		return app.postAuthLogoutHandler
	}

	app.postAuthLogoutHandler = post_auth_logout.New(app.AuthService())
	return app.postAuthLogoutHandler
}

func (app *App) GetAdminHandler() api.Handler {
	if app.getAdminHandler != nil {
		return app.getAdminHandler
	}

	app.getAdminHandler = get_admin.New()
	return app.getAdminHandler
}

func (app *App) GetAdminExportHandler() api.Handler {
	if app.getAdminExportHandler != nil {
		return app.getAdminExportHandler
	}

	app.getAdminExportHandler = get_admin_export.New()
	return app.getAdminExportHandler
}

func (app *App) PostBalanceDepositHandler() api.Handler {
	if app.postBalanceDepositHandler != nil {
		return app.postBalanceDepositHandler
	}

	app.postBalanceDepositHandler = post_balance_deposit.New(app.UsersService())
	return app.postBalanceDepositHandler
}

func (app *App) PostBalanceWithdrawHandler() api.Handler {
	if app.postBalanceWithdrawHandler != nil {
		return app.postBalanceWithdrawHandler
	}

	app.postBalanceWithdrawHandler = post_balance_withdraw.New(app.UsersService())
	return app.postBalanceWithdrawHandler
}

func (app *App) GetOrdersHandler() api.Handler {
	if app.getOrdersHandler != nil {
		return app.getOrdersHandler
	}

	app.getOrdersHandler = get_orders.New(app.OrdersService())
	return app.getOrdersHandler
}

func (app *App) GetOrdersIDHandler() api.Handler {
	if app.getOrdersIDHandler != nil {
		return app.getOrdersIDHandler
	}

	app.getOrdersIDHandler = get_orders_id.New(app.OrdersService(), app.UsersService())
	return app.getOrdersIDHandler
}

func (app *App) PostAuthRegisterHandler() api.Handler {
	if app.postAuthRegisterHandler != nil {
		return app.postAuthRegisterHandler
	}

	app.postAuthRegisterHandler = post_auth_register.New(app.AuthService())
	return app.postAuthRegisterHandler
}

func (app *App) PostOrdersHandler() api.Handler {
	if app.postOrdersHandler != nil {
		return app.postOrdersHandler
	}

	app.postOrdersHandler = post_orders.New(app.OrdersService())
	return app.postOrdersHandler
}

func (app *App) PatchUsersIDHandler() api.Handler {
	if app.patchUsersIDHandler != nil {
		return app.patchUsersIDHandler
	}

	app.patchUsersIDHandler = patch_users_id.New(app.UsersService())
	return app.patchUsersIDHandler
}

func (app *App) PatchOrdersIDHandler() api.Handler {
	if app.patchOrdersIDHandler != nil {
		return app.patchOrdersIDHandler
	}

	app.patchOrdersIDHandler = patch_orders_id.New(app.OrdersService())
	return app.patchOrdersIDHandler
}
