package modeltests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/guigoebel/user-order-api/user-api/api/controllers"
	"github.com/guigoebel/user-order-api/user-api/api/models"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

var server = controllers.Server{}
var userInstance = models.User{}
var orderInstance = models.Order{}

func TestMain(m *testing.M) {
	var err error
	err = godotenv.Load(os.ExpandEnv("../../.env"))
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	}
	Database()
	os.Exit(m.Run())
}

func Database() {
	var err error
	TestDbDriver := os.Getenv("TestDbDriver")
	if TestDbDriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("TestDbUser"), os.Getenv("TestDbPassword"), os.Getenv("TestDbHost"), os.Getenv("TestDbPort"), os.Getenv("TestDbName"))
		server.DB, err = gorm.Open(TestDbDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database\n", TestDbDriver)
		}
	}

	if TestDbDriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", os.Getenv("TestDbHost"), os.Getenv("TestDbPort"), os.Getenv("TestDbUser"), os.Getenv("TestDbName"), os.Getenv("TestDbPassword"))
		server.DB, err = gorm.Open(TestDbDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database\n", TestDbDriver)
		}
	}
}

func refreshUserTable() error {
	err := server.DB.DropTableIfExists(&models.User{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed table")
	return nil
}

func seedOneUser() (models.User, error) {

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user := models.User{
		Name:     "Pet",
		Email:    "email@mail.com",
		Password: "password",
	}

	if err := server.DB.Model(&models.User{}).Create(&user).Error; err != nil {
		log.Fatalf("cannot seed users table: %v", err)
	}

	return user, err
}

func seedUsers() error {

	users := []models.User{
		models.User{
			Name:     "Joao Zackaria",
			Email:    "zackaria@mail.com",
			Password: "password",
		},
		models.User{
			Name:     "Joao Pedro",
			Email:    "pedro@mail.com",
			Password: "password",
		},
	}

	for i, _ := range users {
		err := server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func refreshUserAndOrderTable() error {
	err := server.DB.DropTableIfExists(&models.User{}, &models.Order{}).Error
	if err != nil {
		return err
	}

	err = server.DB.AutoMigrate(&models.User{}, &models.Order{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed table")
	return nil
}

func seedOneUserAndOneOrder() (models.Order, error) {

	err := refreshUserAndOrderTable()
	if err != nil {
		log.Fatal(err)
	}

	user := models.User{
		Name:     "Petcovish Pet",
		Email:    "pet@mail.com",
		Password: "password",
	}

	err = server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return models.Order{}, err
	}

	order := models.Order{
		Description: "Order description",
		Quantity:    10,
		Price:       1000,
		UserID:      user.ID,
	}

	err = server.DB.Model(&models.Order{}).Create(&order).Error
	if err != nil {
		return models.Order{}, err
	}

	return order, nil

}

func seedUsersAndOrders() ([]models.User, []models.Order, error) {

	var err error

	if err != nil {
		return []models.User{}, []models.Order{}, err
	}
	var users = []models.User{
		models.User{
			Name:     "Joao Zackaria",
			Email:    "zackaria@mail.com",
			Password: "password",
		},
		models.User{
			Name:     "Joao Pedro",
			Email:    "pedro@mail.com",
			Password: "password",
		},
	}

	var orders = []models.Order{
		models.Order{
			Description: "Order description",
			Quantity:    10,
			Price:       1000,
		},
		models.Order{
			Description: "Order description 2",
			Quantity:    5,
			Price:       500,
		},
	}

	for i, _ := range users {
		err = server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		orders[i].UserID = users[i].ID

		err = server.DB.Model(&models.Order{}).Create(&orders[i]).Error
		if err != nil {
			log.Fatalf("cannot seed orders table: %v", err)
		}
	}
	return users, orders, nil
}
