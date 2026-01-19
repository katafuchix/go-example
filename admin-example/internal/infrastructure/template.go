package infrastructure

import (
	"io"

	"github.com/labstack/echo/v4"
	"github.com/yosssi/ace"
)

type TemplateRenderer struct{}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	// 複雑なことはせず、単にロードして実行するだけ
	tpl, err := ace.Load("views/layout/base", "views/"+name, nil)
	if err != nil {
		return err
	}
	return tpl.Execute(w, data)
}
