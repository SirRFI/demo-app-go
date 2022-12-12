package fakestore_test

import (
	"demo-app-go/fakestore"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewAddProductCommand(t *testing.T) {
	t.Run("create successfully", func(t *testing.T) {
		// Sample data from the API should be considered valid.
		for _, sample := range products() {
			sample := sample
			t.Run(fmt.Sprintf("product id %d", sample.Id), func(t *testing.T) {
				t.Parallel()
				result, err := fakestore.NewAddProductCommand(
					sample.Title,
					sample.Price,
					sample.Description,
					sample.Category,
					sample.Image,
				)
				require.NoError(t, err, "unexpected error")

				require.Equal(t, sample.Title, result.Title())
				require.Equal(t, sample.Price, result.Price())
				require.Equal(t, sample.Description, result.Description())
				require.Equal(t, sample.Category, result.Category())
				require.Equal(t, sample.Image, result.Image())
			})
		}
	})
	t.Run("fails at validation", func(t *testing.T) {
		validSample := products()[16-1]
		// Map's key describes the scenario.
		// Input reuses valid sample, but in each test a copy of it is modified accordingly.
		// Lastly, the error message in included for comparison.
		samples := map[string]struct {
			input         fakestore.Product
			expectedError string
		}{
			"title is empty": {
				input: func(p fakestore.Product) fakestore.Product {
					p.Title = ""
					return p
				}(validSample),
				expectedError: "title must not be empty",
			},
			"title has only spaces": {
				input: func(p fakestore.Product) fakestore.Product {
					p.Title = "   "
					return p
				}(validSample),
				expectedError: "title must not be empty",
			},
			"negative price": {
				input: func(p fakestore.Product) fakestore.Product {
					p.Price = -1.23
					return p
				}(validSample),
				expectedError: "price must not be negative",
			},
			"no category": {
				input: func(p fakestore.Product) fakestore.Product {
					p.Category = ""
					return p
				}(validSample),
				expectedError: "undefined category",
			},
		}
		for name, sample := range samples {
			sample := sample
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				_, err := fakestore.NewAddProductCommand(
					sample.input.Title,
					sample.input.Price,
					sample.input.Description,
					sample.input.Category,
					sample.input.Image,
				)
				require.ErrorContains(t, err, sample.expectedError)
			})
		}
	})
}

func TestNewUpdateProductCommand(t *testing.T) {
	t.Run("create successfully", func(t *testing.T) {
		// Sample data from the API should be considered valid.
		for _, input := range products() {
			t.Run(fmt.Sprintf("product id %d", input.Id), func(t *testing.T) {
				t.Parallel()
				result, err := fakestore.NewUpdateProductCommand(
					input.Id,
					input.Title,
					input.Price,
					input.Description,
					input.Category,
					input.Image,
				)
				require.NoError(t, err, "unexpected error")

				require.Equal(t, input.Id, result.Id())
				require.Equal(t, input.Title, result.Title())
				require.Equal(t, input.Price, result.Price())
				require.Equal(t, input.Description, result.Description())
				require.Equal(t, input.Category, result.Category())
				require.Equal(t, input.Image, result.Image())
			})
		}
	})
	t.Run("fails at validation", func(t *testing.T) {
		validSample := products()[16-1]
		samples := map[string]struct {
			input         fakestore.Product
			expectedError string
		}{
			"title is empty": {
				input: func(p fakestore.Product) fakestore.Product {
					p.Title = ""
					return p
				}(validSample),
				expectedError: "title must not be empty",
			},
			"title has only spaces": {
				input: func(p fakestore.Product) fakestore.Product {
					p.Title = "   "
					return p
				}(validSample),
				expectedError: "title must not be empty",
			},
			"negative price": {
				input: func(p fakestore.Product) fakestore.Product {
					p.Price = -1.23
					return p
				}(validSample),
				expectedError: "price must not be negative",
			},
			"no category": {
				input: func(p fakestore.Product) fakestore.Product {
					p.Category = ""
					return p
				}(validSample),
				expectedError: "undefined category",
			},
		}
		for name, sample := range samples {
			sample := sample
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				_, err := fakestore.NewAddProductCommand(
					sample.input.Title,
					sample.input.Price,
					sample.input.Description,
					sample.input.Category,
					sample.input.Image,
				)
				require.ErrorContains(t, err, sample.expectedError)
			})
		}
	})
}
