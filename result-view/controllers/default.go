package controllers

import (
	"github.com/astaxie/beego"
	"github.com/jinzhu/gorm"
	"fmt"
	"github.com/astaxie/beego/utils/pagination"
	"github.com/VoidArtanis/dA-crawler/result-view/models"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {

	db, dbErr := gorm.Open("postgres", "host=localhost user=postgres dbname=dA_Data sslmode=disable password=")
	if dbErr != nil {
		fmt.Println(dbErr)
	}
	defer db.Close()

	var count int
	db.Table("artworks").Count(&count)

	postsPerPage := 20
	paginator := pagination.SetPaginator(c.Ctx, postsPerPage, int64(count))
	var data []models.Artwork
	 db.Order("fav_count desc").Offset(paginator.Offset()).Limit(postsPerPage).Find(&data)
	// fetch the next 20 posts
	db.Table("artworks").Count(&count)
	c.Data["posts"] = data

	c.TplName = "index.tpl"


}
