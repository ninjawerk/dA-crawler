package routers

import (
	"github.com/VoidArtanis/dA-crawler/result-view/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
}
