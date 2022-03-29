package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Order struct {
	ID          uint32    `gorm:"primary_key;auto_increment" json:"id"`
	User        User      `json:"user"`
	UserID      uint32    `gorm:"size:255;not null;unique" json:"user_id"`
	Description string    `gorm:"size:255;not null;unique" json:"item_description"`
	Quantity    int64     `gorm:"not null" json:"item_quantity"`
	Price       int64     `gorm:"not null" json:"item_price"`
	TotalValue  int64     `gorm:"not null" json:"total_value"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (o *Order) Prepare() {
	o.User = User{}
	o.Description = html.EscapeString(strings.TrimSpace(o.Description))
	o.CreatedAt = time.Now()
	o.UpdatedAt = time.Now()
}

func (o *Order) Validate() error {

	if o.Description == "" {
		return errors.New("Required Description")

	}
	if o.Quantity == 0 {
		return errors.New("Required Quantity")

	}
	if o.Price == 0 {
		return errors.New("Required Price")

	}
	if o.UserID < 1 {
		return errors.New("Required UserID")
	}
	return nil
}

func (o *Order) SaveOrder(db *gorm.DB) (*Order, error) {
	var err error

	o.TotalValue = o.Quantity * o.Price

	err = db.Debug().Model(&Order{}).Create(&o).Error
	if err != nil {
		return &Order{}, err
	}
	if o.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", o.UserID).Take(&o.User).Error
		if err != nil {
			return &Order{}, err
		}
	}
	return o, nil
}

func (o *Order) FindAllOrders(db *gorm.DB) (*[]Order, error) {
	var err error
	orders := []Order{}
	err = db.Debug().Model(&Order{}).Limit(100).Order("created_at desc").Find(&orders).Error
	if err != nil {
		return &[]Order{}, err
	}
	if len(orders) > 0 {
		for i, _ := range orders {
			err := db.Debug().Model(&User{}).Where("id = ?", orders[i].UserID).Take(&orders[i].User).Error
			if err != nil {
				return &[]Order{}, err
			}
		}
	}
	return &orders, nil
}

func (o *Order) FindOrderByID(db *gorm.DB, oid uint64) (*Order, error) {
	var err error
	err = db.Debug().Model(&Order{}).Where("id = ?", oid).Take(&o).Error
	if err != nil {
		return &Order{}, err
	}
	if o.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", o.UserID).Take(&o.User).Error
		if err != nil {
			return &Order{}, err
		}
	}
	return o, nil
}

func (o *Order) UpdateAOrder(db *gorm.DB) (*Order, error) {

	var err error

	err = db.Debug().Model(&Order{}).Where("id = ?", o.ID).Updates(Order{Description: o.Description, Quantity: o.Quantity, Price: o.Price, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Order{}, err
	}
	if o.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", o.UserID).Take(&o.User).Error
		if err != nil {
			return &Order{}, err
		}
	}
	return o, nil
}

func (o *Order) DeleteAOrder(db *gorm.DB, oid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&Order{}).Where("id = ? and user_id = ?", oid, uid).Take(&Order{}).Delete(&Order{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Order not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
