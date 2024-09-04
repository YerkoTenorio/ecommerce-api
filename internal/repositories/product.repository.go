package repositories

import (
	"ecommerce-api/internal/models"
	"ecommerce-api/pkg/utils"
	"errors"

	"gorm.io/gorm"
)

func GetProductById(productID uint) (*models.Product, error) {
	var product models.Product
	if err := utils.DB.Where("id = ?", productID).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("el producto no existe")
		}
		return nil, err
	}
	return &product, nil

}

func GetProductsByIds(productIDs []uint) ([]models.Product, error) {
	var products []models.Product
	if err := utils.DB.Where("id IN ?", productIDs).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}
