package gf

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/huyongchao98/go-admin/tests/common"
)

func TestGf(t *testing.T) {
	common.ExtraTest(httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(newHandler()),
			Jar:       httpexpect.NewJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
	}))
}
