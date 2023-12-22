package handlers

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/s6352410016/go-fiber-gorm-rest-api-crud-multiple-upload-files-postgresql/database"
	"github.com/s6352410016/go-fiber-gorm-rest-api-crud-multiple-upload-files-postgresql/models"
)

func Create(c *fiber.Ctx) error {
	product := new(models.Product)
	if err := c.BodyParser(product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request Data",
		})
	}

	if product.ProductName == "" || product.ProductPrice == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Input Is Required",
		})
	}

	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	files := form.File["images"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Files Is Required",
		})
	}

	allowExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}

	productImages := []string{}

	for _, file := range files {
		fileExt := filepath.Ext(file.Filename)

		if !allowExts[fileExt] {
			for _, fileName := range productImages {
				if err := os.Remove("./images/" + fileName); err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"message": err,
					})
				}
			}

			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid File Extension",
			})
		}

		newFileName := fmt.Sprintf("%s%s", uuid.New().String(), fileExt)
		if err := c.SaveFile(file, "./images/"+newFileName); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err,
			})
		}

		productImages = append(productImages, newFileName)
	}

	product.ProductImages = productImages
	result := database.DB.Create(&product)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": result.Error,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(product)
}

func GetAll(c *fiber.Ctx) error {
	var products []models.Product
	database.DB.Find(&products)
	return c.JSON(products)
}

func GetById(c *fiber.Ctx) error {
	productId, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Product Id Is Require Type Integer",
		})
	}

	var product models.Product
	database.DB.First(&product, productId)

	if product.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Product Not Found",
		})
	}

	return c.JSON(product)
}

func Update(c *fiber.Ctx) error {
	productId, err := c.ParamsInt("id")
	product := new(models.Product)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Product Id Is Require Type Integer",
		})
	}

	database.DB.First(&product, productId)
	if product.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Product Not Found",
		})
	}

	if err := c.BodyParser(product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request Data",
		})
	}

	if product.ProductName == "" || product.ProductPrice == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Input Is Required",
		})
	}

	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	files := form.File["images"]
	productImages := []string{}

	// กรณีอัพเดดรูป
	if len(files) != 0 {
		allowExts := map[string]bool{
			".jpg":  true,
			".jpeg": true,
			".png":  true,
			".webp": true,
		}

		for _, file := range files {
			fileExt := filepath.Ext(file.Filename)

			if !allowExts[fileExt] {
				for _, fileName := range productImages {
					os.Remove("./images/" + fileName)
				}

				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"message": "Invalid File Extension",
				})
			}

			newFileName := fmt.Sprintf("%s%s", uuid.New().String(), fileExt)

			if err := c.SaveFile(file, "./images/"+newFileName); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": err,
				})
			}

			productImages = append(productImages, newFileName)
		}

		for _, fileName := range product.ProductImages {
			os.Remove("./images/" + fileName)
		}

		product.ProductImages = productImages

		result := database.DB.Save(&product)
		if result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": result.Error,
			})
		}

		return c.JSON(product)
	}

	//กรณีไม่อัพเดดรูป
	result := database.DB.Save(&product)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": result.Error,
		})
	}

	return c.JSON(product)
}

func Delete(c *fiber.Ctx) error {
	productId, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Product Id Is Require Type Integer",
		})
	}

	var product models.Product
	database.DB.First(&product, productId)

	if product.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Product Not Found",
		})
	}

	for _, fileName := range product.ProductImages {
		os.Remove("./images/" + fileName)
	}

	database.DB.Delete(&product)

	return c.JSON(fiber.Map{
		"message": "Product Deleted Successfully",
	})
}

func GetImage(c *fiber.Ctx) error {
	fileName := c.Params("filename")

	file, err := os.Open("./images/" + fileName)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Image Not Found",
		})
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Image Not Found",
		})
	}

	return c.Send(fileData)
}
