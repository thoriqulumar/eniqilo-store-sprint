package controller

import (
	"eniqilo-store/model"
	"eniqilo-store/service"
	cerr "eniqilo-store/utils/error"
	"net/http"
	"net/url"
	"strconv"
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

func (ctr *ProductController) GetProduct(c echo.Context) error {
	// get query param
	value, err := c.FormParams()
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "params not valid"})
	}

	// query to service
	data, err := ctr.ProductService.GetProduct(c.Request().Context(), parseGetProductParams(value))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	// compose response
	return c.JSON(http.StatusOK, model.GetProductResponse{
		Message: "success",
		Data:    data,
	})
}

func (ctr *ProductController) GetProductCustomer(c echo.Context) error {
	// get query param
	value, err := c.FormParams()
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "params not valid"})
	}

	// query to service
	data, err := ctr.ProductService.GetProduct(c.Request().Context(), parseGetProductParams(value))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	// compose response
	return c.JSON(http.StatusOK, model.GetProductResponse{
		Message: "success",
		Data:    data,
	})
}

func (ctr *ProductController) PostProduct(c echo.Context) error {
	var product model.Product
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid product data"})
	}

	// Call the service to create a new product
	createdProduct, err := ctr.ProductService.CreateProduct(c.Request().Context(), product)
	if err != nil {
		return c.JSON(cerr.GetCode(err), echo.Map{"error": err.Error()})
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

func (ctr *ProductController) UpdateProduct(c echo.Context) error {
	idParam := c.Param("id") // Assuming you're using Echo framework and the ID is passed as a URL parameter
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid product ID format"})
	}

	var product model.Product
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid product data"})
	}
	product.ID = id

	// Call the service to create a new product
	up, err := ctr.ProductService.UpdateProduct(c.Request().Context(), product)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	// Prepare the response with the structured format
	response := model.UpdateProductResponse{
		Message: "success",
		Data:    up,
	}
	return c.JSON(http.StatusOK, response)
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

func parseGetProductParams(params url.Values) model.GetProductParam {
	var result model.GetProductParam

	for key, values := range params {
		switch key {
		case "id":
			id, err := uuid.Parse(values[0])
			if err == nil {
				result.ID = &id
			}
		case "limit":
			limit, err := strconv.Atoi(values[0])
			if err == nil {
				result.Limit = &limit
			}
		case "offset":
			offset, err := strconv.Atoi(values[0])
			if err == nil {
				result.Offset = &offset
			}
		case "name":
			result.Name = &values[0]
		case "isAvailable":
			isAvailable, err := strconv.ParseBool(values[0])
			if err == nil {
				result.IsAvailable = &isAvailable
			}
		case "category":
			// Assuming Category is a string
			cat := model.Category(values[0])
			result.Category = &cat
		case "sku":
			result.SKU = &values[0]
		case "inStock":
			inStock, err := strconv.ParseBool(values[0])
			if err == nil {
				result.InStock = &inStock
			}
		// param sorting in set
		case "price":
			result.Sort.Price = &values[0]
		case "createdAt":
			result.Sort.CreatedAt = &values[0]
		}

	}

	return result
}
