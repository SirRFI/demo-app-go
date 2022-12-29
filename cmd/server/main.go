package main

import (
	"demo-app-go/fakestore"
	"demo-app-go/handlers"
	"demo-app-go/storage"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
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

	db, err := sqlx.Connect("mysql", os.Getenv("DATABASE_DSN"))
	if err != nil {
		log.Fatalf("Failed connecting to the database: %s", err)
	}
	taskRepository := storage.NewTaskRepository(db)
	taskHandler := handlers.NewTaskHandler(taskRepository)

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
	e.GET("/tasks", taskHandler.List)
	e.POST("/tasks", taskHandler.Add)
	e.GET("/tasks/:id", taskHandler.Get)
	e.PUT("/tasks/:id", taskHandler.Update)
	e.DELETE("/tasks/:id", taskHandler.Delete)

	e.Logger.Fatal(e.Start(":8000"))
}

type RequestValidator struct {
	validator *validator.Validate
}

func (cv *RequestValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
