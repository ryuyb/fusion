package route

import (
	"github.com/gofiber/fiber/v3"
	"go.uber.org/fx"
)

type Router interface {
	RegisterRouters(router fiber.Router)
}

type RouterRegistry struct {
	routers []Router
}

type RouterRegistryParams struct {
	fx.In
	Routers []Router `group:"routers"`
}

func NewRouterRegistry(params RouterRegistryParams) *RouterRegistry {
	return &RouterRegistry{
		routers: params.Routers,
	}
}

func (r *RouterRegistry) RegisterAllRoutes(app *fiber.App) {
	for _, router := range r.routers {
		router.RegisterRouters(app)
	}
}

func asRoute(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(Router)),
		fx.ResultTags(`group:"routers"`),
	)
}

var Module = fx.Module("route",
	fx.Provide(
		asRoute(NewSwaggerRoute),
		asRoute(NewHealthRoute),
		asRoute(NewUserRoute),

		NewRouterRegistry,
	),
	fx.Invoke(func(app *fiber.App, routerRegistry *RouterRegistry) {
		routerRegistry.RegisterAllRoutes(app)
	}),
)
