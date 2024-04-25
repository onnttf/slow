package petrol

import (
	"slow/dal"
	"slow/util/input"
	"slow/util/output"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterRoutes(g *echo.Group) {
	g.GET("/price", getPrice)
}

func getPrice(c echo.Context) error {
	var i struct {
		Type int    `query:"type" validate:"required,oneof=0 92 95 98"`
		Area string `query:"area"`
	}
	err := input.BindAndValidate(c, &i)
	if err != nil {
		return output.Fail(c, err)
	}
	rows, err := dal.NewDao(dal.Slow).QueryList(c.Request().Context(), func(d *gorm.DB) *gorm.DB {
		d.Where("type = ?", i.Type)
		if len(i.Area) > 0 {
			d.Where("area = ?", i.Area)
		}
		return d
	})
	if err != nil {
		return output.Fail(c, err)
	}
	type price struct {
		Price       string `json:"price"`
		ReleaseDate string `json:"release_date"`
	}
	list := make([]price, 0, len(rows))
	for _, v := range rows {
		list = append(list, price{
			Price:       v.Price,
			ReleaseDate: v.ReleaseDate,
		})
	}
	return output.Success(c, map[string]interface{}{
		"list": list,
	})
}
