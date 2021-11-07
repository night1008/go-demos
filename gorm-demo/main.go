package main

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Dog struct {
	ID   int
	Name string
	Toys []Toy `gorm:"polymorphic:Owner;"`
}

func (d *Dog) BeforeDelete(tx *gorm.DB) error {
	// 支持事务
	fmt.Println("===> BeforeDelete ", d)
	tx.Where("owner_type", "dogs").Where("owner_id", 1).Where("name", "toy3").Delete(&Toy{})
	return nil
	// return fmt.Errorf("some delete error")
}

type Cat struct {
	ID   int
	Name string
	Toys []Toy `gorm:"polymorphic:Owner;"`
}

type Toy struct {
	// ID        int
	OwnerType string `gorm:"primaryKey"`
	OwnerID   int    `gorm:"primaryKey"`
	Name      string `gorm:"primaryKey;index"`
}

func main2() {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db = db.Debug()
	db.AutoMigrate(&Dog{}, &Cat{}, &Toy{})

	// db.Create(&Cat{ID: 1, Name: "cat1", Toys: []Toy{{Name: "toy1"}, {Name: "toy2"}}})
	// db.Create(&Dog{ID: 1, Name: "dog1", Toys: []Toy{{Name: "toy3"}, {Name: "toy4"}}})

	// 测试连表查询
	// var dogs []Dog
	// db.Preload("Toys").Joins("left join toys on dogs.id = toys.owner_id and toys.owner_type = 'dogs'").Where("toys.name = 'toy1'").Find(&dogs)
	// fmt.Println(dogs)

	// 测试 BeforeDelete hook 生效情况
	// 需要先查询出来再删除 BeforeDelete hook 才能拿到实例
	dog := Dog{ID: 1}
	if err := db.First(&dog).Error; err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(dog)
	db.Delete(&dog)
}
