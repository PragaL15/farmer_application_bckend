package Masterhandlers

import (
	"context"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/PragaL15/go_newBackend/go_backend/db"
	"github.com/go-playground/validator/v10"
)
func InsertMasterProduct(c *fiber.Ctx) error {
	type Request struct {
		CategoryID   int    `json:"category_id" validate:"required,min=1"`
		ProductName  string `json:"product_name" validate:"required,max=100"`
		Status       int    `json:"status" validate:"required"`
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	_, err := db.Pool.Exec(context.Background(), `
		CALL insert_master_product($1, $2, $3);
	`, req.CategoryID, req.ProductName, req.Status)

	if err != nil {
		log.Printf("Failed to insert product: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to insert product"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Product added successfully"})
}
func UpdateMasterProduct(c *fiber.Ctx) error {
	type Request struct {
		ProductID    int32    `json:"product_id" validate:"required,min=1"`
		CategoryID   int32    `json:"category_id" validate:"required,min=1"`
		ProductName  string `json:"product_name" validate:"required,max=100"`
		Status       int32   `json:"status" validate:"required"`
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	_, err := db.Pool.Exec(context.Background(), `
		CALL update_master_product($1, $2, $3, $4);
	`, req.ProductID, req.CategoryID, req.ProductName, req.Status)

	if err != nil {
		log.Printf("Failed to update product: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update product"})
	}

	return c.JSON(fiber.Map{"message": "Product updated successfully"})
}
func DeleteMasterProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID is required"})
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	_, err = db.Pool.Exec(context.Background(), `
		CALL delete_master_product($1);
	`, idInt)

	if err != nil {
		log.Printf("Failed to delete product: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete product"})
	}
	return c.JSON(fiber.Map{"message": "Product deleted successfully"})
}
func GetProducts(c *fiber.Ctx) error {
	rows, err := db.Pool.Query(context.Background(), "SELECT * FROM get_master_products()")
	if err != nil {
		log.Printf("Failed to fetch product records: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch products"})
	}
	defer rows.Close()

	var products []map[string]interface{}

	for rows.Next() {
		var productID, categoryID, status *int
		var productName, categoryName *string

		if err := rows.Scan(&productID, &categoryID, &categoryName, &productName, &status); err != nil {
			log.Printf("Error scanning row: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error processing data"})
		}

		products = append(products, map[string]interface{}{
			"product_id":    productID,
			"category_id":   categoryID,
			"category_name": categoryName,
			"product_name":  productName,
			"status":        status,
		})
	}

	return c.JSON(products)
}
