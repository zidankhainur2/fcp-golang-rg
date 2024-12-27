package api

import (
	"a21hc3NpZ25tZW50/model"
	"a21hc3NpZ25tZW50/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserAPI interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	GetUserTaskCategory(c *gin.Context)
}

type userAPI struct {
	userService service.UserService
}

func NewUserAPI(userService service.UserService) *userAPI {
	return &userAPI{userService}
}

func (u *userAPI) Register(c *gin.Context) {
	var user model.UserRegister

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("invalid decode json"))
		return
	}

	if user.Email == "" || user.Password == "" || user.Fullname == "" {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("register data is empty"))
		return
	}

	var recordUser = model.User{
		Fullname: user.Fullname,
		Email:    user.Email,
		Password: user.Password,
	}

	recordUser, err := u.userService.Register(&recordUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse("error internal server"))
		return
	}

	c.JSON(http.StatusCreated, model.NewSuccessResponse("register success"))
}

func (u *userAPI) Login(c *gin.Context) {
    var user model.UserLogin

    // Parsing JSON request body
    if err := c.BindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, model.NewErrorResponse("invalid decode json")) // Konsisten dengan Register
        return
    }

    // Validasi email dan password kosong
    if user.Email == "" || user.Password == "" {
        c.JSON(http.StatusBadRequest, model.NewErrorResponse("email or password is empty"))
        return
    }

    // Login pengguna dan dapatkan token
    tokenString, err := u.userService.Login(&model.User{
        Email:    user.Email,
        Password: user.Password,
    })
    if err != nil {
        c.JSON(http.StatusUnauthorized, model.NewErrorResponse("wrong email or password"))
        return
    }

    // Set token ke cookie untuk session
    c.SetCookie("session_token", *tokenString, 3600, "/", "", false, true)

    // Response sukses
    c.JSON(http.StatusOK, gin.H{
        "message": "login success",
        "token":   *tokenString,
    })
}

// Fungsi untuk mendapatkan kategori task user
func (u *userAPI) GetUserTaskCategory(c *gin.Context) {
    taskCategories, err := u.userService.GetUserTaskCategory()
    if err != nil {
        c.JSON(http.StatusInternalServerError, model.NewErrorResponse("error fetching task categories"))
        fmt.Println("Invalid category_id:", c.Param("category_id"))
        return
    }

    c.JSON(http.StatusOK, taskCategories)
}
