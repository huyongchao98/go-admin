package controller

import (
	"github.com/chenhg5/go-admin/modules/auth"
	"github.com/chenhg5/go-admin/plugins/admin/models"
	"github.com/chenhg5/go-admin/context"
	"bytes"
	"net/http"
	"github.com/chenhg5/go-admin/modules/menu"
	"strings"
	"path"
	"fmt"
	"github.com/chenhg5/go-admin/template"
	"github.com/chenhg5/go-admin/template/types"
)

// 显示列表
func ShowInfo(ctx *context.Context) {

	defer GlobalDeferHandler(ctx)

	user := ctx.UserValue["user"].(auth.User)

	prefix := ctx.Request.URL.Query().Get("prefix")

	page := ctx.Request.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}
	pageSize := ctx.Request.URL.Query().Get("pageSize")
	if pageSize == "" {
		pageSize = "10"
	}

	sortField := ctx.Request.URL.Query().Get("sort")
	if sortField == "" {
		sortField = "id"
	}
	sortType := ctx.Request.URL.Query().Get("sort_type")
	if sortType == "" {
		sortType = "desc"
	}

	thead, infoList, paninator, title, description := models.GlobalTableList[prefix].GetDataFromDatabase(map[string]string{
		"page":      page,
		"path":      ctx.Path(),
		"sortField": sortField,
		"sortType":  sortType,
		"prefix":    prefix,
		"pageSize":  pageSize,
	})

	var (
		editUrl string
		//newUrl  string
	)
	editUrl = Config.PREFIX + "/info/" + prefix + "/edit" + GetRouteParameterString(page, pageSize, sortType, sortField)
	newUrl := Config.PREFIX + "/info/" + prefix + "/new" + GetRouteParameterString(page, pageSize, sortType, sortField)
	deleteUrl := Config.PREFIX + "/delete/" + prefix

	tmpl, tmplName := template.Get("adminlte").GetTemplate(ctx.Request.Header.Get("X-PJAX") == "true")

	menu.GlobalMenu.SetActiveClass(ctx.Path())

	dataTable := template.Get(Config.THEME).DataTable().SetInfoList(infoList).SetThead(thead).SetEditUrl(editUrl).SetNewUrl(newUrl).SetDeleteUrl(deleteUrl)
	table := dataTable.GetContent()

	box := template.Get(Config.THEME).Box().
		SetBody(table).
		SetHeader(dataTable.GetDataTableHeader()).
		WithHeadBorder(false).
		SetFooter(paninator.GetContent()).
		GetContent()

	ctx.Response.Header.Add("Content-Type", "text/html; charset=utf-8")

	buf := new(bytes.Buffer)
	tmpl.ExecuteTemplate(buf, tmplName, types.Page{
		User: user,
		Menu: menu.GetGlobalMenu(user),
		System: types.SystemInfo{
			"0.0.1",
		},
		Panel: types.Panel{
			Content:     box,
			Description: description,
			Title:       title,
		},
		AssertRootUrl: Config.PREFIX,
		Title:         Config.TITLE,
		Logo:          Config.LOGO,
		MiniLogo:      Config.MINILOGO,
	})
	ctx.WriteString(buf.String())
}

func Assert(ctx *context.Context) {
	filepath := "template/adminlte/resource" + strings.Replace(ctx.Request.URL.Path, Config.PREFIX, "", 1)
	data, err := template.Get("adminlte").GetAsset(filepath)
	fileSuffix := path.Ext(filepath)
	fileSuffix = strings.Replace(fileSuffix, ".", "", -1)

	var contentType = ""
	if fileSuffix == "css" || fileSuffix == "js" {
		contentType = "text/" + fileSuffix + "; charset=utf-8"
	} else {
		contentType = "image/" + fileSuffix
	}

	if err != nil {
		fmt.Println("asset err", err)
		ctx.Write(http.StatusNotFound, map[string]string{}, "")
	} else {
		ctx.Write(http.StatusOK, map[string]string{
			"content-type": contentType,
		}, string(data))
	}
}

func GetRouteParameterString(page, pageSize, sortType, sortField string) string {
	return "?page=" + page + "&pageSize=" + pageSize + "&sort=" + sortField + "&sort_type=" + sortType
}
