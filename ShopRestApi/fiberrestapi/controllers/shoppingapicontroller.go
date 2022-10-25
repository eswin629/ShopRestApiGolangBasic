package controllers

import (
	"ilmudata/fiberrestapi/database"
	"ilmudata/fiberrestapi/models"
	"strconv"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ShoppingAPIController struct {
	// declare variables
	Db *gorm.DB
}
func InitShoppingAPIController() *ShoppingAPIController {
	db := database.InitDb()
	// gorm
	db.AutoMigrate(&models.Shopping{})

	return &ShoppingAPIController{Db: db}
}

func (controller *ShoppingAPIController) Greeting(c *fiber.Ctx) error {
	
	return c.JSON(fiber.Map{
		"message": "welcome...",
	})
}
// GET /products
func (controller *ShoppingAPIController) GetAllShoppings(c *fiber.Ctx) error {
	// load all Shopping
	var shoppings []models.Shopping
	err := models.ReadShopping(controller.Db, &shoppings)
	if err!=nil {
		return c.SendStatus(500) // http 500 internal server error
	}
	return c.JSON(shoppings)
}



// RESOURCE API
// get all, detail  HTTP GET
// insert HTTP POST
// update HTTP PUT  / HTTP POST
// delete HTTP DELETE  / HTTP GET

// POST /products
func (controller *ShoppingAPIController) CreateShopping(c *fiber.Ctx) error {
	/*
	var shopping models.Shopping

	if err := c.BodyParser(&shopping); err != nil {
		return c.SendStatus(400) // bad request http
	}
	// save shopping
	err := models.CreateShopping(controller.Db, &shopping)
	if err!=nil {
		return c.SendStatus(500)
	}
	// if succeed
	return c.JSON(shopping)
	*/

	if form, err := c.MultipartForm(); err == nil {
		files := form.File["image"]
		
		for _, file := range files {
			var data models.Shopping
			fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])
				
			if err := c.BodyParser(&data); err != nil {
				// return c.Redirect("/shoppings")
				return c.SendStatus(400)
			}
			
			if err := c.SaveFile(file, fmt.Sprintf("./public/upload/%s", file.Filename)); err != nil {
				return err
			}

			data.Image = file.Filename
		
			err := models.CreateShopping(controller.Db, &data)
		
			if err != nil {
				return c.SendStatus(400)
			}

			return c.JSON(data)
		}
		// return c.JSON(fiber.Map{
			return c.SendStatus(500)
		// })
	}

	// return c.JSON(fiber.Map{
		return c.SendStatus(500)
	// })

}

// GET /products/productdetail?id=xxx
/*
func (controller *ProductAPIController) GetDetailProduct(c *fiber.Ctx) error {
	id := c.Query("id")
	idn,_ := strconv.Atoi(id)

	var product models.Product
	err := models.ReadProductById(controller.Db, &product, idn)
	if err!=nil {
		return c.SendStatus(500) // http 500 internal server error
	}
	return c.JSON(product)
}
*/
// GET /products/detail/xxx
func (controller *ShoppingAPIController) GetDetailShopping(c *fiber.Ctx) error {
	id := c.Params("id")
	idn,_ := strconv.Atoi(id)


	var shopping models.Shopping
	err := models.ReadShoppingById(controller.Db, &shopping, idn)
	if err!=nil {
		return c.SendStatus(500) // http 500 internal server error
	}
	return c.JSON(shopping)
}

/// PUT /products/:id
func (controller *ShoppingAPIController) EditShopping(c *fiber.Ctx) error {
	id := c.Params("id")
	idn,_ := strconv.Atoi(id)

	var shopping models.Shopping
	err := models.ReadShoppingById(controller.Db, &shopping, idn)
	if err!=nil {
		return c.SendStatus(500) // http 500 internal server error
	}
	var updateShopping models.Shopping

	if err := c.BodyParser(&updateShopping); err != nil {
		return c.SendStatus(400)
	}
	shopping.Name = updateShopping.Name
	shopping.Quantity = updateShopping.Quantity
	shopping.Image = updateShopping.Image
	shopping.Price = updateShopping.Price

	// save product
	models.UpdateShopping(controller.Db, &shopping)
	
	return c.JSON(shopping)	
}

/// DELETE /products/:id
func (controller *ShoppingAPIController) DeleteShopping(c *fiber.Ctx) error {
	id := c.Params("id")
	idn,_ := strconv.Atoi(id)

	var shopping models.Shopping
	models.DeleteShoppingById(controller.Db, &shopping, idn)

	//return c.JSON(product)	
	return c.JSON(fiber.Map{
		"message": "data was deleted",
	})
}
