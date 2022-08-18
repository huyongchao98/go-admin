package gear

import (
	// add gin adapter

	"github.com/teambition/gear"
	"github.com/teambition/gear/middleware/static"

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

	"net/http"
	"os"

	"github.com/GoAdminGroup/themes/adminlte"
	"github.com/huyongchao98/go-admin/engine"
	"github.com/huyongchao98/go-admin/modules/config"
	"github.com/huyongchao98/go-admin/modules/language"
	"github.com/huyongchao98/go-admin/plugins/admin"
	"github.com/huyongchao98/go-admin/plugins/admin/modules/table"
	"github.com/huyongchao98/go-admin/template"
	"github.com/huyongchao98/go-admin/template/chartjs"
	"github.com/huyongchao98/go-admin/tests/tables"
)

func newHandler() http.Handler {
	app := gear.New()

	eng := engine.Default()

	template.AddComp(chartjs.NewChart())

	if err := eng.AddConfigFromJSON(os.Args[len(os.Args)-1]).
		AddGenerators(tables.Generators).
		AddGenerator("user", tables.GetUserTable).
		Use(app); err != nil {
		panic(err)
	}

	eng.HTML("GET", "/admin", tables.GetContent)

	app.Use(static.New(static.Options{Root: "./uploads", Prefix: "uploads"}))

	return app
}

func NewHandler(dbs config.DatabaseList, gens table.GeneratorList) http.Handler {
	app := gear.New()

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
		AddAdapter(admin.NewAdmin(gens)).
		AddGenerators(gens).
		Use(app); err != nil {
		panic(err)
	}

	eng.HTML("GET", "/admin", tables.GetContent)

	app.Use(static.New(static.Options{Root: "./uploads", Prefix: "uploads"}))

	return app
}
