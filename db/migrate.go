package db

// Migrate database (if not already)
func AutoMigrateDatabase() {
	DB.Set("gorm:table_options", "ENGINE=InnoDB")
	users := &User{}
	// userInfo := &UserInfo{}
	parkingMeta := &ParkingLot{}
	slot := &ParkingSlot{}
	BookSlot := &BookSlot{}
	DB.AutoMigrate(users, parkingMeta, slot, BookSlot)
}
