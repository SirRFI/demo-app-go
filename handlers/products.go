package handlers

import (
	"demo-app-go/fakestore"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"strings"
)

type ProductsHandler struct {
	FakeStoreAPI fakestore.FakeStoreAPI
}

func (h *ProductsHandler) GetProducts(c echo.Context) error {
	products, err := h.FakeStoreAPI.GetProducts()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, products)
}

func (h *ProductsHandler) GetProduct(c echo.Context) error {
	id, err := getId(c)
	if err != nil {
		return err
	}

	product, err := h.FakeStoreAPI.GetProduct(id)
	if errors.Is(err, fakestore.ResourceNotFoundError) {
		return c.NoContent(http.StatusNotFound)
	}
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, product)
}

type productRequestBody struct {
	Title       string  `json:"title" validate:"required,min=1"`
	Price       float64 `json:"price" validate:"gte=0"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Image       string  `json:"image" validate:"required,min=1,url"`
}

func (h *ProductsHandler) AddProduct(c echo.Context) error {
	data := &productRequestBody{}
	err := c.Bind(data)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	data.Title = strings.TrimSpace(data.Title)
	data.Category = strings.TrimSpace(data.Category)
	err = c.Validate(data)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	command, err := fakestore.NewAddProductCommand(data.Title, data.Price, data.Description, data.Category, data.Image)
	if err != nil {
		return err
	}
	product, err := h.FakeStoreAPI.AddProduct(command)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, product)
}

func (h *ProductsHandler) UpdateProduct(c echo.Context) error {
	id, err := getId(c)
	if err != nil {
		return err
	}

	data := &productRequestBody{}
	err = c.Bind(data)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	data.Title = strings.TrimSpace(data.Title)
	data.Category = strings.TrimSpace(data.Category)
	err = c.Validate(data)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	command, err := fakestore.NewUpdateProductCommand(
		id,
		data.Title,
		data.Price,
		data.Description,
		data.Category,
		data.Image,
	)
	if err != nil {
		return err
	}

	product, err := h.FakeStoreAPI.UpdateProduct(command)
	if errors.Is(err, fakestore.ResourceNotFoundError) {
		return c.NoContent(http.StatusNotFound)
	}
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, product)
}

func (h *ProductsHandler) DeleteProduct(c echo.Context) error {
	id, err := getId(c)
	if err != nil {
		return err
	}

	err = h.FakeStoreAPI.DeleteProduct(id)
	if errors.Is(err, fakestore.ResourceNotFoundError) {
		return c.NoContent(http.StatusNotFound)
	}
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func getId(c echo.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
