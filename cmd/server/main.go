package main

import (
	"demo-app-go/fakestore"
	"demo-app-go/handlers"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	httpClient := &http.Client{Timeout: time.Minute}
	fakeStoreAPI := fakestore.NewAPI(os.Getenv("FAKESTOREAPI_BASEURL"), httpClient)
	productsHandler := handlers.ProductsHandler{FakeStoreAPI: fakeStoreAPI}

	e := echo.New()
	e.Validator = &RequestValidator{validator: validator.New()}

	e.GET("/", func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})
	e.GET("/products", productsHandler.GetProducts)
	e.POST("/products", productsHandler.AddProduct)
	e.GET("/products/:id", productsHandler.GetProduct)
	e.PUT("/products/:id", productsHandler.UpdateProduct)
	e.DELETE("/products/:id", productsHandler.DeleteProduct)

	e.Logger.Fatal(e.Start(":8000"))
}

type RequestValidator struct {
	validator *validator.Validate
}

func (cv *RequestValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
