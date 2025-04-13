package utils

import "general-shop/models"

func CalculateTotal(items []models.CartItem) float64 {
	var total float64
	for _, item := range items {
		total += float64(item.Quantity) * item.Price
	}
	return total
}
