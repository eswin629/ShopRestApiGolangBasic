package controllers

import (
	"github.com/gofiber/fiber/v2"
	// "github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/crypto/bcrypt"
	"strconv"

	"ilmudata/fiberrestapi/database"
	"ilmudata/fiberrestapi/models"

	
	"gorm.io/gorm"

	
)

type LoginForm struct {
	Username     string  `form:"username" json:"username" validate:"required"`
	Password	 string  `form:"password" json:"password" validate:"required"`
}

type AuthController struct {
	Db *gorm.DB
	// store *session.Store
}
//func InitAuthController(s *session.Store) *AuthController {
func InitAuthController() *AuthController {
	db := database.InitDb()

	db.AutoMigrate(&models.User{})

	// return &AuthController{Db: db, store: s}
	return &AuthController{Db: db}
}
// get /login
func (controller *AuthController) Login(c *fiber.Ctx) error {
	// load all User
	var users []models.User
	err := models.ReadUser(controller.Db, &users)
	if err!=nil {
		return c.SendStatus(500) // http 500 internal server error
	}
	return c.JSON(users)
}

//post
func (controller *AuthController) PostLogin(c *fiber.Ctx) error {
	/*
	sess, err := controller.store.Get(c)

	if err!=nil {
		panic(err)
	}
*/
	var myform models.User
	var data models.User

	if err := c.BodyParser(&myform); err != nil {
		return c.JSON(fiber.Map{"error": err})
	}

	username := myform.Username
	plainPassword := myform.Password

	err2 := models.ReadOneUser(controller.Db, &data, username)

	if err2 != nil {
		return c.Redirect("/login")
	}
	
	hashPassword := data.Password

	check := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(plainPassword))

	status := check == nil

	if status {
		/*
		sess.Set("username", username)
		sess.Save()
		*/
		return c.JSON(myform)
	} else {
		return c.SendStatus(500)
	}

	
}


/*
//register
func (controller *AuthController) Register(c *fiber.Ctx) error {
	return c.Render("register", fiber.Map{
		"Title": "Register",
	})
}
*/


func (controller *AuthController) CreateUser(c *fiber.Ctx) error {

	//register
	var register models.User

		if err := c.BodyParser(&register); err != nil {
			return c.SendStatus(400) // bad request http
		}

		bytes, _ := bcrypt.GenerateFromPassword([]byte(register.Password), 8)
		sHash := string(bytes)
		
		register.Password = sHash

		err := models.Register(controller.Db, &register)

		if err != nil {
			return c.SendStatus(400) // bad request http
		}
		
		// if succeed
		return c.JSON(register)

	
}

// GET /products/productdetail?id=xxx
/*
func (controller *AuthController) GetDetailUser(c *fiber.Ctx) error {
	id := c.Query("id")
	idn,_ := strconv.Atoi(id)

	var user models.User
	err := models.ReadUserById(controller.Db, &user, idn)
	if err!=nil {
		return c.SendStatus(500) // http 500 internal server error
	}
	return c.JSON(user)
}
*/

func (controller *AuthController) GetDetailUser(c *fiber.Ctx) error {
	id := c.Params("id")
	idn,_ := strconv.Atoi(id)


	var user models.User
	err := models.ReadUserById(controller.Db, &user, idn)
	if err!=nil {
		return c.SendStatus(500) // http 500 internal server error
	}
	return c.JSON(user)
}


/// DELETE /user/:id
func (controller *AuthController) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	idn,_ := strconv.Atoi(id)

	var user models.User
	models.DeleteUserById(controller.Db, &user, idn)

	//return c.JSON(user)	
	return c.JSON(fiber.Map{
		"message": "data was deleted",
	})
}
