package gofiber

import (
	// add fasthttp adapter
	ada "github.com/huyongchao98/go-admin/adapter/gofiber"
	// add mysql driver
	_ "github.com/huyongchao98/go-admin/modules/db/drivers/mysql"
	// add postgresql driver
	_ "github.com/huyongchao98/go-admin/modules/db/drivers/postgres"
	// add sqlite driver
	_ "github.com/huyongchao98/go-admin/modules/db/drivers/sqlite"
	// add mssql driver
	_ "github.com/huyongchao98/go-admin/modules/db/drivers/mssql"
	// add adminlte ui theme
	_ "github.com/GoAdminGroup/themes/adminlte"

	"os"

	"github.com/GoAdminGroup/themes/adminlte"
	"github.com/gofiber/fiber/v2"
	"github.com/huyongchao98/go-admin/engine"
	"github.com/huyongchao98/go-admin/modules/config"
	"github.com/huyongchao98/go-admin/modules/language"
	"github.com/huyongchao98/go-admin/plugins/admin"
	"github.com/huyongchao98/go-admin/plugins/admin/modules/table"
	"github.com/huyongchao98/go-admin/template"
	"github.com/huyongchao98/go-admin/template/chartjs"
	"github.com/huyongchao98/go-admin/tests/tables"
	"github.com/valyala/fasthttp"
)

func newHandler() fasthttp.RequestHandler {
	app := fiber.New(fiber.Config{
		ServerHeader: "Fiber",
	})

	eng := engine.Default()

	adminPlugin := admin.NewAdmin(tables.Generators).AddDisplayFilterXssJsFilter()
	adminPlugin.AddGenerator("user", tables.GetUserTable)

	template.AddComp(chartjs.NewChart())

	if err := eng.AddConfigFromJSON(os.Args[len(os.Args)-1]).
		AddPlugins(adminPlugin).
		Use(app); err != nil {
		panic(err)
	}

	eng.HTML("GET", "/admin", tables.GetContent)

	return app.Handler()
}

func NewHandler(dbs config.DatabaseList, gens table.GeneratorList) fasthttp.RequestHandler {
	app := fiber.New(fiber.Config{
		ServerHeader: "Fiber",
	})

	eng := engine.Default()

	template.AddComp(chartjs.NewChart())

	if err := eng.AddConfig(&config.Config{
		Databases: dbs,
		UrlPrefix: "admin",
		Store: config.Store{
			Path:   "./uploads",
			Prefix: "uploads",
		},
		Language:    language.EN,
		IndexUrl:    "/",
		Debug:       true,
		ColorScheme: adminlte.ColorschemeSkinBlack,
	}).
		AddAdapter(new(ada.Gofiber)).
		AddGenerators(gens).
		Use(app); err != nil {
		panic(err)
	}

	eng.HTML("GET", "/admin", tables.GetContent)

	return app.Handler()
}
