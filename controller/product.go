package controller

import (
	"eniqilo-store/model"
	"eniqilo-store/service"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/labstack/echo/v4"
)

type ProductController struct {
	ProductService service.ProductService
}

func NewProductController(productService service.ProductService) *ProductController {
	return &ProductController{
		ProductService: productService,
	}
}

func (ctr *ProductController) PostProduct(c echo.Context) error {
	var product model.CreatedProduct
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid product data"})
	}

	// Call the service to create a new product
	createdProduct, err := ctr.ProductService.CreateProduct(c.Request().Context(), product)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	// Prepare the response with the structured format
	response := model.PostProductResponse{
		Message: "success",
		Data: model.Data{
			ID:        createdProduct.ID.String(),
			CreatedAt: createdProduct.CreatedAt.Format(time.RFC3339),
		},
	}
	return c.JSON(http.StatusCreated, response)
}

func (ctr *ProductController) DeleteProduct(c echo.Context) error {
	idParam := c.Param("id") // Assuming you're using Echo framework and the ID is passed as a URL parameter
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid product ID format"})
	}

	err = ctr.ProductService.DeleteProduct(c.Request().Context(), id)
	if err != nil {
		if err.Error() == "no product found with the given ID" {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "Product not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Product successfully deleted"})
}
