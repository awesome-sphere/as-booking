package models

import (
	"log"
	"math/rand"
	"time"
)

var DONE_SEEDING bool

func SeedThreater() []int {
	theaters := []Theater{
		{Location: "Central Salaya"},
		{Location: "Central Westgate"},
		{Location: "Central Pinklao"},
		{Location: "Siam Paragon"},
		{Location: "The Crystal Rachaprunk"},
	}
	DB.Create(&theaters)
	var theater_ids []int
	for _, theater := range theaters {
		theater_ids = append(theater_ids, int(theater.ID))
	}
	return theater_ids

}

func SeedSeatType() []SeatType {
	seat_type := []SeatType{
		{Price: 200, Type: Standard},
		{Price: 450, Type: Premium},
	}
	DB.Create(seat_type)
	return seat_type
}

func SeedTimeSlot(theater_ids []int) map[int][]int {
	time_slots := make(map[int][]int)
	for _, theater_id := range theater_ids {
		for movie := 1; movie <= 5; movie++ {
			for day_count := 0; day_count < 5; day_count++ {
				for slot := 0; slot < 5; slot++ {
					hour := rand.Intn(4-2) + 2
					minute := rand.Intn(50-5) + 5
					time_slot := TimeSlot{
						MovieID:   movie,
						TheaterID: theater_id,
						Time:      time.Date(2022, 7, 22+day_count, 1+slot+hour, minute, 0, 0, time.UTC),
					}
					DB.Create(&time_slot)
					if _, ok := time_slots[theater_id]; ok {
						time_slots[theater_id] = append(time_slots[theater_id], int(time_slot.ID))
					} else {
						time_slots[theater_id] = []int{int(time_slot.ID)}

					}
				}
			}
		}
	}
	return time_slots
}

func SeedSeat(time_slots map[int][]int, seat_types []SeatType) {
	for theater_id, theater_time_slots := range time_slots {
		for _, seat_type := range seat_types {
			is_standard := seat_type.Type == Standard
			// for theater_id, theater_time_slots := range time_slots {
			for _, theater_time_slot := range theater_time_slots {
				amount := 0
				start := 1
				if is_standard {
					amount = 40
				} else {
					start = 41
					amount = 15
				}
				for i := start; i < start+amount; i++ {
					if is_standard {
						seat := SeatInfo{
							TimeSlotID: theater_time_slot,
							TheaterID:  theater_id,
							SeatTypeID: int(seat_type.ID),
							SeatNumber: i,
							Status:     Available,
						}
						DB.Create(&seat)
					} else if !is_standard {
						seat := SeatInfo{
							TimeSlotID: theater_time_slot,
							TheaterID:  theater_id,
							SeatTypeID: int(seat_type.ID),
							SeatNumber: i,
							Status:     Available,
						}
						DB.Create(&seat)
					}
				}

			}
		}
	}

}

func SeedData() {
	log.Println("================ Start Seeding Data ================")
	theater_ids := SeedThreater()
	seat_type := SeedSeatType()
	time_slots := SeedTimeSlot(theater_ids)
	SeedSeat(time_slots, seat_type)
	log.Println("================ Done Seeding ================")
	DONE_SEEDING = true
}
