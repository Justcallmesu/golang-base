package routes

import (
	"github.com/gin-gonic/gin"
	"justcallmesu.com/rest-api/internal/app"
	"justcallmesu.com/rest-api/internal/app/blogs"
)

func BlogRoutes(engine *gin.Engine, services *app.Services, middlewares *app.Middlewares) {
	blogRouter := engine.Group("/blogs")

	//Handler
	blogHandler := blogs.NewBlogHandler(services.BlogService)

	// Routes
	blogRouter.GET("/", blogHandler.FindMany)
	blogRouter.GET("/:id", blogHandler.FindOne)
	blogRouter.POST("/", blogHandler.Create)
	blogRouter.PUT("/", blogHandler.Update)
	blogRouter.POST("/:id/upload", blogHandler.Upload)
	blogRouter.DELETE("/:id", blogHandler.Delete)
}
