package dev

import (
	"fmt"

	"github.com/Aiscom-LLC/meals-api/config"
	"github.com/Aiscom-LLC/meals-api/domain"
	"github.com/Aiscom-LLC/meals-api/repository/enums"
	"github.com/Aiscom-LLC/meals-api/utils"
)

const seedName string = "init admin"

// CreateAdmin init admin user
func CreateAdmin() {
	seedExists := config.DB.Where("name = ?", seedName).First(&domain.Seed{}).Error
	if seedExists != nil {
		seed := domain.Seed{
			Name: seedName,
		}

		superAdmin := domain.User{
			FirstName: "super",
			LastName:  "admin",
			Email:     "meals@aisnovations.com",
			Password:  utils.HashString("Password12!"),
			Status:    &enums.StatusTypesEnum.Active,
			Role:      enums.UserRoleEnum.SuperAdmin,
		}

		config.DB.Create(&superAdmin)
		config.DB.Create(&seed)
		fmt.Println("=== Admin seed created ===")
	} else {
		fmt.Printf("Seed `%s` already exists \n", seedName)
	}
}
