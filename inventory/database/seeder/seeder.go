package main

import (
	"fmt"
	"log"
	"synapsis/inventory/app/model"
	"synapsis/inventory/config"
	"synapsis/inventory/database/connection"

	"go.uber.org/dig"
	"gorm.io/gorm"
)

func main() {
	var err error
	// Initialize the application bootstrap
	container := dig.New()

	if err = container.Provide(func() config.AppConfig {
		return config.NewAppConfig()
	}); err != nil {
		panic(err)
	}

	// provide postgres connection
	if err = container.Provide(func(cfg config.AppConfig) connection.DBInstance {
		return connection.NewDatabaseInstance(cfg)
	}); err != nil {
		panic(err)
	}

	// run the migration
	if err = container.Invoke(func(dbInstance connection.DBInstance) error {
		db := dbInstance.Database()

		if err := runSeeder(db); err != nil {
			return err
		}

		return nil
	}); err != nil {
		panic(err)
	}
}

func runSeeder(db *gorm.DB) error {
	if err := seedProduct(db); err != nil {
		return err
	}

	return nil
}

func seedProduct(db *gorm.DB) error {
	var count int64
	if err := db.Model(&model.Product{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		log.Println("product already seeded. Skipping...")
		return nil
	}

	products := []model.Product{
		{
			Name:        "Laptop Lenovo 11 Inch",
			Sku:         "00001",
			Description: "Small Laptop From Lenovo",
			Price:       5000000,
		},
		{
			Name:        "Laptop Lenovo 12 Inch",
			Sku:         "00002",
			Description: "Medium Laptop From Lenovo",
			Price:       6000000,
		},
		{
			Name:        "Laptop Lenovo 13 Inch",
			Sku:         "00003",
			Description: "Big Laptop From Lenovo",
			Price:       7000000,
		},
		{
			Name:        "Laptop HP 11 Inch",
			Sku:         "00004",
			Description: "Small Laptop From HP",
			Price:       5500000,
		},
		{
			Name:        "Laptop HP 12 Inch",
			Sku:         "00005",
			Description: "Medium Laptop From HP",
			Price:       6500000,
		},
		{
			Name:        "Laptop HP 13 Inch",
			Sku:         "00006",
			Description: "Big Laptop From HP",
			Price:       7500000,
		},
	}

	for _, product := range products {
		if err := db.Create(&product).Error; err != nil {
			return err
		}

		stock := model.Stock{
			ProductId:      product.ID,
			TotalStock:     100,
			AvailableStock: 100,
			ReservedStock:  0,
		}
		if err := db.Create(&stock).Error; err != nil {
			return err
		}

		if err := db.Create(&model.StockMovement{
			ProductId:   product.ID,
			ChangeType:  "ADD",
			Quantity:    100,
			ReferenceId: product.ID,
			Note:        fmt.Sprintf("Initial stock for product %d", product.ID),
		}).Error; err != nil {
			return err
		}
	}

	return nil
}
