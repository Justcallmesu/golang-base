package routes

import (
	"github.com/gin-gonic/gin"
	"justcallmesu.com/rest-api/internal/app"
	"justcallmesu.com/rest-api/internal/app/blogs"
)

func BlogDetailsRoutes(engine *gin.Engine, services *app.Services, middlewares *app.Middlewares) {
	blogDetailsRouter := engine.Group("/blog-details")

	//Handler
	blogDetailsHandler := blogs.NewBlogDetailsHandler(services.BlogDetailService)

	// Routes
	blogDetailsRouter.POST("/", blogDetailsHandler.CreateOne)
	blogDetailsRouter.PUT("/order", blogDetailsHandler.UpdateDetailsOrder)

}
