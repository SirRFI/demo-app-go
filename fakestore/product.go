package fakestore

import (
	"errors"
	"strings"
)

type Product struct {
	Id          uint          `json:"id"`
	Title       string        `json:"title"`
	Price       float64       `json:"price"`
	Description string        `json:"description"`
	Category    string        `json:"category"`
	Image       string        `json:"image"`
	Rating      ProductRating `json:"rating"`
}

type ProductRating struct {
	Rate  float64 `json:"rate"`
	Count uint    `json:"count"`
}

// AddProductCommand is used for creating new Product.
// Unexported fields and NewAddProductCommand ensures that the command is created and passed in
// valid state and cannot be changed (immutability).
type AddProductCommand struct {
	title       string
	price       float64
	description string
	category    string
	image       string
}

// NewAddProductCommand creates AddProductCommand and validates the data.
// FakeStore does not have any requirements for the data, so the validation below is purely for demo purposes.
func NewAddProductCommand(
	title string,
	price float64,
	description string,
	category string,
	image string,
) (AddProductCommand, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return AddProductCommand{}, errors.New("title must not be empty")
	}
	if price < 0 {
		return AddProductCommand{}, errors.New("price must not be negative")
	}
	category = strings.TrimSpace(category)
	if category == "" {
		return AddProductCommand{}, errors.New("undefined category")
	}

	return AddProductCommand{
		title:       title,
		price:       price,
		description: description,
		category:    category,
		image:       image,
	}, nil
}

func (a AddProductCommand) Title() string {
	return a.title
}

func (a AddProductCommand) Price() float64 {
	return a.price
}

func (a AddProductCommand) Description() string {
	return a.description
}

func (a AddProductCommand) Category() string {
	return a.category
}

func (a AddProductCommand) Image() string {
	return a.image
}

// UpdateProductCommand is used for updating existing Product.
// Just like in AddProductCommand, fields are not exported to ensure data validity and immutability.
type UpdateProductCommand struct {
	id          uint
	title       string
	price       float64
	description string
	category    string
	image       string
}

// NewUpdateProductCommand creates UpdateProductCommand and validates the data.
// FakeStore does not have any requirements for the data, so the validation below is purely for demo purposes.
func NewUpdateProductCommand(
	id uint,
	title string,
	price float64,
	description string,
	category string,
	image string,
) (UpdateProductCommand, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return UpdateProductCommand{}, errors.New("title must not be empty")
	}
	if price < 0 {
		return UpdateProductCommand{}, errors.New("price must not be negative")
	}
	category = strings.TrimSpace(category)
	if category == "" {
		return UpdateProductCommand{}, errors.New("undefined category")
	}

	return UpdateProductCommand{
		id:          id,
		title:       title,
		price:       price,
		description: description,
		category:    category,
		image:       image,
	}, nil
}

func (u UpdateProductCommand) Id() uint {
	return u.id
}

func (u UpdateProductCommand) Title() string {
	return u.title
}

func (u UpdateProductCommand) Price() float64 {
	return u.price
}

func (u UpdateProductCommand) Description() string {
	return u.description
}

func (u UpdateProductCommand) Category() string {
	return u.category
}

func (u UpdateProductCommand) Image() string {
	return u.image
}
