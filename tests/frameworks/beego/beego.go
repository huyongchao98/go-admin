package beego

import (
	// add beego adapter
	"github.com/GoAdminGroup/themes/adminlte"
	_ "github.com/huyongchao98/go-admin/adapter/beego"
	"github.com/huyongchao98/go-admin/modules/config"
	"github.com/huyongchao98/go-admin/modules/language"
	"github.com/huyongchao98/go-admin/plugins/admin/modules/table"

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

	"github.com/astaxie/beego"
	"github.com/huyongchao98/go-admin/engine"
	"github.com/huyongchao98/go-admin/plugins/admin"
	"github.com/huyongchao98/go-admin/plugins/example"
	"github.com/huyongchao98/go-admin/template"
	"github.com/huyongchao98/go-admin/template/chartjs"
	"github.com/huyongchao98/go-admin/tests/tables"
)

func newHandler() http.Handler {

	app := beego.NewApp()

	eng := engine.Default()
	adminPlugin := admin.NewAdmin(tables.Generators)
	adminPlugin.AddGenerator("user", tables.GetUserTable)

	examplePlugin := example.NewExample()

	if err := eng.AddConfigFromJSON(os.Args[len(os.Args)-1]).
		AddPlugins(adminPlugin, examplePlugin).Use(app); err != nil {
		panic(err)
	}

	template.AddComp(chartjs.NewChart())

	eng.HTML("GET", "/admin", tables.GetContent)

	beego.BConfig.Listen.HTTPAddr = "127.0.0.1"
	beego.BConfig.Listen.HTTPPort = 9087

	return app.Handlers
}

func NewHandler(dbs config.DatabaseList, gens table.GeneratorList) http.Handler {

	app := beego.NewApp()

	eng := engine.Default()
	adminPlugin := admin.NewAdmin(gens)

	if err := eng.AddConfig(config.Config{
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
		AddPlugins(adminPlugin).Use(app); err != nil {
		panic(err)
	}

	template.AddComp(chartjs.NewChart())

	eng.HTML("GET", "/admin", tables.GetContent)

	beego.BConfig.Listen.HTTPAddr = "127.0.0.1"
	beego.BConfig.Listen.HTTPPort = 9087

	return app.Handlers
}
