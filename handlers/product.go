package handlers

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/s6352410016/go-fiber-gorm-rest-api-crud-mssql/database"
	"github.com/s6352410016/go-fiber-gorm-rest-api-crud-mssql/models"
)

func GetAll(c *fiber.Ctx) error {
	var products []models.Product
	database.DB.Find(&products)
	return c.Status(fiber.StatusOK).JSON(products)
}

func GetById(c *fiber.Ctx) error {
	productId := c.Params("id")
	var product models.Product
	database.DB.First(&product, productId)

	if product.ProductId == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Product Not Found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(product)
}

func Create(c *fiber.Ctx) error {
	product := new(models.Product)
	productName := c.FormValue("productName")
	productPrice := c.FormValue("productPrice")
	if productName == "" || productPrice == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Input Is Required",
		})
	}

	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "File Image Is Required",
		})
	}

	fileExt := filepath.Ext(file.Filename)
	allowExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}

	if !allowExts[fileExt] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid File Extension",
		})
	}

	newFileName := fmt.Sprintf("%s%s", uuid.New().String(), fileExt)
	if err := c.SaveFile(file, "./images/"+newFileName); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error While Save File",
		})
	}

	product.ProductName = productName

	value, err := strconv.ParseFloat(productPrice, 64)
	if err != nil {
		log.Fatal("Error Parsing String To Float64")
	}

	product.ProductPrice = value
	product.ProductImage = newFileName
	database.DB.Create(&product)
	return c.Status(fiber.StatusCreated).JSON(product)
}

func Update(c *fiber.Ctx) error {
	productId := c.Params("id")
	var product models.Product
	database.DB.First(&product, productId)

	if product.ProductId == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Product Not Found",
		})
	}

	productName := c.FormValue("productName")
	productPrice := c.FormValue("productPrice")
	if productName == "" || productPrice == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Input Is Required",
		})
	}

	file, _ := c.FormFile("image")

	// กรณีอัพรูป
	if file != nil {
		fileExt := filepath.Ext(file.Filename)
		allowExts := map[string]bool{
			".jpg":  true,
			".jpeg": true,
			".png":  true,
			".webp": true,
		}

		if !allowExts[fileExt] {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid File Extension",
			})
		}

		newFileName := fmt.Sprintf("%s%s", uuid.New().String(), fileExt)
		if err := c.SaveFile(file, "./images/"+newFileName); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error While Save File",
			})
		}

		if err := os.Remove("./images/" + product.ProductImage); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error While Unlink Current File",
			})
		}

		product.ProductName = productName

		value, err := strconv.ParseFloat(productPrice, 64)
		if err != nil {
			log.Fatal("Error Parsing String To Float64")
		}

		product.ProductPrice = value
		product.ProductImage = newFileName
		database.DB.Save(&product)
		return c.Status(fiber.StatusOK).JSON(product)
	}

	// กรณีไม่อัพรูป
	product.ProductName = productName

	value, err := strconv.ParseFloat(productPrice, 64)
	if err != nil {
		log.Fatal("Error Parsing String To Float64")
	}

	product.ProductPrice = value
	database.DB.Save(&product)
	return c.Status(fiber.StatusOK).JSON(product)
}

func Delete(c *fiber.Ctx) error {
	productId := c.Params("id")
	var product models.Product
	database.DB.First(&product, productId)

	if product.ProductId == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Product Not Found",
		})
	}

	if err := os.Remove("./images/" + product.ProductImage); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error While Unlink Current File",
		})
	}

	database.DB.Delete(&product)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product Deleted Successfully",
	})
}

func GetImage(c *fiber.Ctx) error {
	imageName := c.Params("filename")

	//ทำอ่านไฟล์จาก directory
	filePath := "./images/" + imageName
	file, err := os.Open(filePath)
	//กรณีไม่เจอรูปที่เรียกมา
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "File Image Not Found",
		})
	}
	defer file.Close()

	//ทำอานอ่านไฟล์รูปภาพและส่งไปแสดงผล
	fileData, err := io.ReadAll(file)
	//กรณีไม่เจอรูปที่เรียกมา
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "File Image Not Found",
		})
	}

	return c.Send(fileData)
}
