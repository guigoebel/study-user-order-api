package seed

import (
	"log"

	"github.com/guigoebel/user-order-api/user-api/api/models"

	"github.com/jinzhu/gorm"
)

var users = []models.User{
	models.User{
		Name:        "Ezrael Pabla",
		Email:       "ezrael@gmail.com",
		Password:    "password",
		Cpf:         "00000000000",
		PhoneNumber: "999999999",
	},
	models.User{
		Name:        "Luther Lex",
		Email:       "lex@gmail.com",
		Password:    "password",
		Cpf:         "00000000000",
		PhoneNumber: "999999999",
	},
}

var orders = []models.Order{
	models.Order{
		Description: "Title 1",
		Quantity:    1,
		Price:       3,
	},
	models.Order{
		Description: "Title 2",
		Quantity:    2,
		Price:       4,
	},
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.Order{}, &models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}, &models.Order{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	err = db.Debug().Model(&models.Order{}).AddForeignKey("user_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	for i, _ := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		orders[i].UserID = users[i].ID

		err = db.Debug().Model(&models.Order{}).Create(&orders[i]).Error
		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		}
	}
}
