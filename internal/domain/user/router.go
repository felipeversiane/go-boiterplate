package user

import (
	"github.com/felipeversiane/go-boiterplate/internal/infra/database"
	"github.com/gin-gonic/gin"
)

func UserRouter(g *gin.RouterGroup, db database.DatabaseInterface) *gin.RouterGroup {

	user := g.Group("/user")
	{

	}

	return user
}
