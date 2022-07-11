package models

import (
	"fmt"
	"log"
	"strconv"

	"github.com/awesome-sphere/as-seating/utils"
	"github.com/shopspring/decimal"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/sharding"
)

var DB *gorm.DB

func InitDatabase() {
	dbUser := utils.GetenvOr("POSTGRES_USER", "pkinwza")
	dbPassword := utils.GetenvOr("POSTGRES_PASSWORD", "securepassword")
	dbHost := utils.GetenvOr("POSTGRES_HOST", "127.0.0.1")
	dbPort := utils.GetenvOr("POSTGRES_PORT", "5432")
	dbName := utils.GetenvOr("POSTGRES_DB", "as-cinema")
	theater_amount := utils.GetenvOr("THEATER_AMOUNT", "5")
	theater_number, err := strconv.Atoi(theater_amount)
	if err != nil {
		log.Fatal("Fail to convert type of THEATER_AMOUNT")
	}

	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
	)

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&SeatType{})
	db.AutoMigrate(&Seat{}, &Theater{})
	for i := 1; i <= theater_number; i++ {
		time_slot_table := fmt.Sprintf("time_slots_%01d", i)
		db.AutoMigrate(&TimeSlot{})
		db.Migrator().RenameTable("time_slots", time_slot_table)

		seat_info_table := fmt.Sprintf("seat_infos_%01d", i)
		db.AutoMigrate(&SeatInfo{})
		db.Migrator().RenameTable("seat_infos", seat_info_table)

	}
	db.Use(sharding.Register(sharding.Config{
		ShardingKey:         "theater_id",
		NumberOfShards:      uint(theater_number + 1),
		PrimaryKeyGenerator: sharding.PKPGSequence,
	}, "time_slots", "seat_infos"))
	DB = db
	// // Testing
	db.Create(&SeatType{ID: 1, Price: decimal.NewFromFloat(69.9), Type: Standard})
	db.Create(&Theater{ID: 2, Location: "Pattaya"})
	db.Create(&TimeSlot{ID: 2, MovieId: 1, TheaterId: 2})
	db.Create(&SeatInfo{TheaterId: 2, TimeSlotId: 2, SeatTypeId: 1})
}
