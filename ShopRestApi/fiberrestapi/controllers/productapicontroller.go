package controllers

import (
	"ilmudata/fiberrestapi/database"
	"ilmudata/fiberrestapi/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ProductAPIController struct {
	// declare variables
	Db *gorm.DB
}
func InitProductAPIController() *ProductAPIController {
	db := database.InitDb()
	// gorm
	db.AutoMigrate(&models.Product{})

	return &ProductAPIController{Db: db}
}

func (controller *ProductAPIController) Greeting(c *fiber.Ctx) error {
	
	return c.JSON(fiber.Map{
		"message": "welcome...",
	})
}
// GET /products
func (controller *ProductAPIController) GetAllProducts(c *fiber.Ctx) error {
	// load all products
	var products []models.Product
	err := models.ReadProducts(controller.Db, &products)
	if err!=nil {
		return c.SendStatus(500) // http 500 internal server error
	}
	return c.JSON(products)
}



// RESOURCE API
// get all, detail  HTTP GET
// insert HTTP POST
// update HTTP PUT  / HTTP POST
// delete HTTP DELETE  / HTTP GET

// POST /products
func (controller *ProductAPIController) CreateProduct(c *fiber.Ctx) error {
	var product models.Product

	if err := c.BodyParser(&product); err != nil {
		return c.SendStatus(400) // bad request http
	}
	// save product
	err := models.CreateProduct(controller.Db, &product)
	if err!=nil {
		return c.SendStatus(500)
	}
	// if succeed
	return c.JSON(product)
}

// GET /products/productdetail?id=xxx
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
// GET /products/detail/xxx
func (controller *ProductAPIController) GetDetailProduct2(c *fiber.Ctx) error {
	id := c.Params("id")
	idn,_ := strconv.Atoi(id)


	var product models.Product
	err := models.ReadProductById(controller.Db, &product, idn)
	if err!=nil {
		return c.SendStatus(500) // http 500 internal server error
	}
	return c.JSON(product)
}

/// PUT /products/:id
func (controller *ProductAPIController) EditProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	idn,_ := strconv.Atoi(id)

	var product models.Product
	err := models.ReadProductById(controller.Db, &product, idn)
	if err!=nil {
		return c.SendStatus(500) // http 500 internal server error
	}
	var updateProduct models.Product

	if err := c.BodyParser(&updateProduct); err != nil {
		return c.SendStatus(400)
	}
	product.Name = updateProduct.Name
	product.Quantity = updateProduct.Quantity
	product.Price = updateProduct.Price

	// save product
	models.UpdateProduct(controller.Db, &product)
	
	return c.JSON(product)	
}

/// DELETE /products/:id
func (controller *ProductAPIController) DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	idn,_ := strconv.Atoi(id)

	var product models.Product
	models.DeleteProductById(controller.Db, &product, idn)

	//return c.JSON(product)	
	return c.JSON(fiber.Map{
		"message": "data was deleted",
	})
}
