package user

import "github.com/gin-gonic/gin"

type userController struct {
	service UserUsecaseInterface
}

type UserControllerInterface interface {
	Create(c *gin.Context)
	Retrieve(c *gin.Context)
	List(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

func NewUserController(service UserUsecaseInterface) UserControllerInterface {
	return &userController{service}
}

func (controller *userController) Create(c *gin.Context) {

}

func (controller *userController) Retrieve(c *gin.Context) {

}

func (controller *userController) List(c *gin.Context) {

}

func (controller *userController) Update(c *gin.Context) {

}

func (controller *userController) Delete(c *gin.Context) {

}
