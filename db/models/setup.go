package models

import (
	"fmt"
	"log"
	"strconv"

	"github.com/awesome-sphere/as-booking/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/sharding"
)

var DB *gorm.DB
var THEATER_AMOUNT int

func InitDatabase() {
	dbUser := utils.GetenvOr("POSTGRES_USER", "pkinwza")
	dbPassword := utils.GetenvOr("POSTGRES_PASSWORD", "securepassword")
	dbHost := utils.GetenvOr("POSTGRES_HOST", "127.0.0.1")
	dbPort := utils.GetenvOr("POSTGRES_PORT", "5432")
	dbName := utils.GetenvOr("POSTGRES_DB", "as-cinema")
	theater_amount := utils.GetenvOr("THEATER_AMOUNT", "5")
	theater_number, err := strconv.Atoi(theater_amount)
	if err != nil {
		log.Println("Fail to convert type of THEATER_AMOUNT")
	}
	THEATER_AMOUNT = theater_number

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
		log.Println(err)
	}
	var does_table_exist bool

	db.AutoMigrate(&SeatType{}, &Theater{})
	// db.AutoMigrate(&Seat{}, &Theater{})
	for i := 0; i <= theater_number; i++ {
		time_slot_table := fmt.Sprintf("time_slots_%01d", i)
		if !db.Migrator().HasTable(time_slot_table) {
			db.AutoMigrate(&TimeSlot{})
			db.Migrator().RenameTable("time_slots", time_slot_table)
			if !does_table_exist {
				does_table_exist = true
			}
		}

		seat_info_table := fmt.Sprintf("seat_infos_%01d", i)
		if !db.Migrator().HasTable(seat_info_table) {
			db.AutoMigrate(&SeatInfo{})
			db.Migrator().RenameTable("seat_infos", seat_info_table)
			if !does_table_exist {
				does_table_exist = true
			}
		}

	}
	db.Use(sharding.Register(sharding.Config{
		ShardingKey:         "theater_id",
		NumberOfShards:      uint(theater_number + 1),
		PrimaryKeyGenerator: sharding.PKPGSequence,
	}, "time_slots", "seat_infos"))
	DB = db
	if does_table_exist {
		go SeedData()
	} else {
		DONE_SEEDING = true
	}
}
