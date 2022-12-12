package fakestore

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

const defaultCollectionSize = 20
const responseMaxSize = 102_400

// FakeStoreAPI documentation: https://fakestoreapi.com/docs
type FakeStoreAPI interface {
	GetProducts() ([]Product, error)
	GetProduct(id uint) (*Product, error)
	AddProduct(AddProductCommand) (*Product, error)
	UpdateProduct(UpdateProductCommand) (*Product, error)
	DeleteProduct(id uint) error
}

var ResourceNotFoundError = errors.New("resource not found")

type API struct {
	baseURL string
	client  *http.Client
}

func NewAPI(baseURL string, client *http.Client) *API {
	return &API{baseURL, client}
}

func (a API) GetProducts() ([]Product, error) {
	response, err := a.sendRequest(http.MethodGet, "/products", nil)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API responded with status %d", response.StatusCode)
	}

	// For simplicity of the demo, response data is not validated.
	products := make([]Product, 0, defaultCollectionSize)
	err = readResponseJSON(response.Body, responseMaxSize, &products)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (a API) GetProduct(id uint) (*Product, error) {
	response, err := a.sendRequest(http.MethodGet, fmt.Sprintf("/products/%d", id), nil)
	if err != nil {
		return nil, err
	}
	// This API does not return 404 Not Found. Instead, it returns empty response body.
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API responded with status %d", response.StatusCode)
	}
	data, err := readResponse(response.Body, responseMaxSize)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, ResourceNotFoundError
	}

	product := &Product{}
	err = json.Unmarshal(data, product)
	if err != nil {
		return nil, fmt.Errorf("unmarshal response json: %w", err)
	}

	return product, nil
}

type addProductRequestBody struct {
	Title       string  `json:"title"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Image       string  `json:"image"`
}

func (a API) AddProduct(command AddProductCommand) (*Product, error) {
	response, err := a.sendRequest(
		http.MethodPost,
		"/products",
		addProductRequestBody{
			Title:       command.Title(),
			Price:       command.Price(),
			Description: command.Description(),
			Category:    command.Category(),
			Image:       command.Image(),
		},
	)
	if err != nil {
		return nil, err
	}
	// This API does not return 201 Created status, neither validates data anyhow to return 400 Bad Request.
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API responded with status %d", response.StatusCode)
	}

	// The API also returns only assigned id, instead of entire object.
	// For simplicity of the demo, it's not validated.
	product := &Product{
		Title:       command.Title(),
		Price:       command.Price(),
		Description: command.Description(),
		Category:    command.Category(),
		Image:       command.Image(),
	}
	err = readResponseJSON(response.Body, responseMaxSize, product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

type updateProductRequestBody struct {
	Title       string  `json:"title"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Image       string  `json:"image"`
}

func (a API) UpdateProduct(command UpdateProductCommand) (*Product, error) {
	response, err := a.sendRequest(
		http.MethodPut,
		fmt.Sprintf("/products/%d", command.Id()),
		updateProductRequestBody{
			Title:       command.Title(),
			Price:       command.Price(),
			Description: command.Description(),
			Category:    command.Category(),
			Image:       command.Image(),
		},
	)
	if err != nil {
		return nil, err
	}
	// This API does not return 404 Not Found. Unlike GetProduct, it always returns id and nothing else.
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API responded with status %d", response.StatusCode)
	}

	// For simplicity of the demo, the returned is not validated.
	product := &Product{
		Id:          command.Id(),
		Title:       command.Title(),
		Price:       command.Price(),
		Description: command.Description(),
		Category:    command.Category(),
		Image:       command.Image(),
	}
	err = readResponseJSON(response.Body, responseMaxSize, product)
	if err != nil {
		return nil, err
	}
	if product.Id != command.Id() { // sanity check
		return nil, fmt.Errorf("expected product id %d, got %d", command.Id(), product.Id)
	}

	return product, nil
}

func (a API) DeleteProduct(id uint) error {
	response, err := a.sendRequest(http.MethodDelete, fmt.Sprintf("/products/%d", id), nil)
	if err != nil {
		return err
	}

	// This API does not return 204 No Content or 404 Not Found.
	// Instead, it returns the entire resource if it existed, and null when it doesn't.
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("API responded with status %d", response.StatusCode)
	}
	data, err := readResponse(response.Body, 4)
	if err != nil {
		return err
	}
	if string(data) == "null" {
		return ResourceNotFoundError
	}

	return nil
}

func (a API) sendRequest(method string, path string, data any) (*http.Response, error) {
	var payload io.Reader
	if data != nil {
		reqData, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("marshaling json for request: %w", err)
		}
		payload = bytes.NewReader(reqData)
	}

	request, err := http.NewRequest(method, a.baseURL+path, payload)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	if payload != nil {
		request.Header.Add("Content-Type", "application/json; charset=utf-8")
	}

	response, err := a.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	return response, nil
}

func readResponse(body io.ReadCloser, limit int64) ([]byte, error) {
	defer func(closer io.Closer) {
		err := closer.Close()
		if err != nil {
			log.Printf("Closing response body failed: %s", err)
		}
	}(body)
	data, err := io.ReadAll(io.LimitReader(body, limit))
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}
	return data, nil
}

func readResponseJSON(body io.ReadCloser, limit int64, object any) error {
	data, err := readResponse(body, limit)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &object)
	if err != nil {
		return fmt.Errorf("unmarshal response json: %w", err)
	}

	return nil
}
