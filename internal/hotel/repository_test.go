package hotel

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestSaveHotel(t *testing.T) {
	// Set up mock database connection with sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock database connection: %v", err)
	}
	defer db.Close()

	// Set expectation for sqlite_version query using a regular expression for flexibility
	mock.ExpectQuery(`(?i)^SELECT sqlite_version\(\)$`).WillReturnRows(sqlmock.NewRows([]string{"sqlite_version"}).AddRow("3.32.3"))

	// Open GORM DB from mock sql.DB
	gormDB, err := gorm.Open(sqlite.New(sqlite.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to initialize GORM: %v", err)
	}

	// Create an instance of hotelRepository
	repo := NewRepository(gormDB)

	// Sample hotel data with UUID
	hotelID := uuid.New()
	hotel := &Hotel{
		ID:           hotelID,
		OwnerName:    "Test Owner",
		OwnerSurname: "Test Surname",
		CompanyTitle: "Test Company",
		ContactInfos: []ContactInfo{},
	}

	// Expectation: a successful call to Create method with backticks around the table name
	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO `+"`hotels`"+` \(`).
		WithArgs(hotel.OwnerName, hotel.OwnerSurname, hotel.CompanyTitle, hotel.ID.String()). // Pass UUID as string
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Test Save method
	err = repo.Save(hotel)
	assert.NoError(t, err)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestDeleteHotel_Repository(t *testing.T) {
	// Set up mock database connection with sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock database connection: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`(?i)^SELECT sqlite_version\(\)$`).WillReturnRows(sqlmock.NewRows([]string{"sqlite_version"}).AddRow("3.32.3"))

	// Open GORM DB from mock sql.DB
	gormDB, err := gorm.Open(sqlite.New(sqlite.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to initialize GORM: %v", err)
	}

	// Create an instance of hotelRepository
	repo := NewRepository(gormDB)

	// Sample hotel ID
	hotelID := uuid.New()

	// Expectation: a successful call to Delete method
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM ` + "`hotels`" + ` WHERE id = ?`).
		WithArgs(hotelID.String()). // Pass UUID as string
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Test Delete method
	err = repo.Delete(hotelID)
	assert.NoError(t, err)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestAddContactInfo_Repository(t *testing.T) {
	// Set up mock database connection with sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock database connection: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`(?i)^SELECT sqlite_version\(\)$`).WillReturnRows(sqlmock.NewRows([]string{"sqlite_version"}).AddRow("3.32.3"))

	// Open GORM DB from mock sql.DB
	gormDB, err := gorm.Open(sqlite.New(sqlite.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to initialize GORM: %v", err)
	}

	// Create an instance of hotelRepository
	repo := NewRepository(gormDB)

	// Sample data
	hotelUUID := uuid.New()
	contact := &ContactInfo{
		ID:          uuid.New(),
		InfoType:    "phone",
		InfoContent: "1234567890",
	}

	// Expectation: a successful call to Create method for ContactInfo
	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO `+"`contact_infos`"+` \(`).
		WithArgs(hotelUUID.String(), contact.InfoType, contact.InfoContent, contact.ID.String()). // Fix order here
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Test AddContactInfo method
	err = repo.AddContactInfo(hotelUUID, contact)
	assert.NoError(t, err)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestListHotels_Repository(t *testing.T) {
	// Mock database connection setup
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock database connection: %v", err)
	}
	defer db.Close()

	// Expectation for sqlite_version query
	mock.ExpectQuery(`(?i)^SELECT sqlite_version\(\)$`).WillReturnRows(sqlmock.NewRows([]string{"sqlite_version"}).AddRow("3.32.3"))

	// Open GORM DB with the mock sql.DB
	gormDB, err := gorm.Open(sqlite.New(sqlite.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to initialize GORM: %v", err)
	}

	// Repository initialization
	repo := NewRepository(gormDB)

	// Sample hotel data
	hotels := []Hotel{
		{ID: uuid.New(), OwnerName: "Owner 1", OwnerSurname: "Surname 1", CompanyTitle: "Company 1"},
		{ID: uuid.New(), OwnerName: "Owner 2", OwnerSurname: "Surname 2", CompanyTitle: "Company 2"},
	}

	// Expectation for querying hotels
	mock.ExpectQuery(`(?i)^SELECT .* FROM ` + "`hotels`" + `.*`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "owner_name", "owner_surname", "company_title"}).
			AddRow(hotels[0].ID.String(), hotels[0].OwnerName, hotels[0].OwnerSurname, hotels[0].CompanyTitle).
			AddRow(hotels[1].ID.String(), hotels[1].OwnerName, hotels[1].OwnerSurname, hotels[1].CompanyTitle))

		// Expectation for querying contact_infos
	mock.ExpectQuery(`(?i)^SELECT .* FROM `+"`contact_infos`"+`.*`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()). // Match any hotel_id (UUID)
		WillReturnRows(sqlmock.NewRows([]string{"id", "hotel_id", "contact_info"}).
			AddRow(uuid.New().String(), hotels[0].ID.String(), "contact1").
			AddRow(uuid.New().String(), hotels[1].ID.String(), "contact2"))

	// Test ListHotels method
	result, err := repo.ListHotels()
	assert.NoError(t, err)   // No error should occur
	assert.Len(t, result, 2) // We should have 2 hotels

	// Verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestGetHotelOfficials(t *testing.T) {
	// Set up mock database connection with sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock database connection: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`(?i)^SELECT sqlite_version\(\)$`).WillReturnRows(sqlmock.NewRows([]string{"sqlite_version"}).AddRow("3.32.3"))

	// Open GORM DB from mock sql.DB
	gormDB, err := gorm.Open(sqlite.New(sqlite.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to initialize GORM: %v", err)
	}

	// Create an instance of hotelRepository
	repo := NewRepository(gormDB)

	// Sample hotel official data
	officials := []HotelOfficial{
		{OwnerName: "Owner 1", OwnerSurname: "Surname 1", CompanyTitle: "Company 1"},
		{OwnerName: "Owner 2", OwnerSurname: "Surname 2", CompanyTitle: "Company 2"},
	}

	// Expectation: a successful call to Find method to list hotel officials
	mock.ExpectQuery(`(?i)^SELECT .* FROM ` + "`hotels`" + `.*`).
		WillReturnRows(sqlmock.NewRows([]string{"owner_name", "owner_surname", "company_title"}).
			AddRow(officials[0].OwnerName, officials[0].OwnerSurname, officials[0].CompanyTitle).
			AddRow(officials[1].OwnerName, officials[1].OwnerSurname, officials[1].CompanyTitle))

	// Test GetHotelOfficials method
	result, err := repo.GetHotelOfficials()
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestGetHotelOfficials_NoRecords(t *testing.T) {
	// Set up mock database connection with sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock database connection: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`(?i)^SELECT sqlite_version\(\)$`).WillReturnRows(sqlmock.NewRows([]string{"sqlite_version"}).AddRow("3.32.3"))

	// Open GORM DB from mock sql.DB
	gormDB, err := gorm.Open(sqlite.New(sqlite.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to initialize GORM: %v", err)
	}

	// Create an instance of hotelRepository
	repo := NewRepository(gormDB)

	// Expectation: no records for hotel officials
	mock.ExpectQuery(`(?i)^SELECT .* FROM ` + "`hotels`" + `.*`).
		WillReturnRows(sqlmock.NewRows([]string{"owner_name", "owner_surname", "company_title"}))

	// Test GetHotelOfficials method when no data is found
	result, err := repo.GetHotelOfficials()
	assert.NoError(t, err)
	assert.Len(t, result, 0)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestGetHotelDetails_Repository(t *testing.T) {
	// Set up mock database connection with sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock database connection: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`(?i)^SELECT sqlite_version\(\)$`).WillReturnRows(sqlmock.NewRows([]string{"sqlite_version"}).AddRow("3.32.3"))

	// Open GORM DB from mock sql.DB
	gormDB, err := gorm.Open(sqlite.New(sqlite.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to initialize GORM: %v", err)
	}

	// Create an instance of hotelRepository
	repo := NewRepository(gormDB)

	// Sample hotel data
	hotel := &Hotel{
		ID:           uuid.New(),
		OwnerName:    "Owner Test",
		OwnerSurname: "Surname Test",
		CompanyTitle: "Company Test",
		ContactInfos: []ContactInfo{{ID: uuid.New(), InfoType: "phone", InfoContent: "1234567890"}},
	}

	// Expectation: querying hotel details by ID
	mock.ExpectQuery(`(?i)^SELECT .* FROM ` + "`hotels`" + `.*`).
		WithArgs(hotel.ID.String()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "owner_name", "owner_surname", "company_title"}).
			AddRow(hotel.ID.String(), hotel.OwnerName, hotel.OwnerSurname, hotel.CompanyTitle))

	// Expectation: querying contact_infos for the hotel
	mock.ExpectQuery(`(?i)^SELECT .* FROM ` + "`contact_infos`" + `.*`).
		WithArgs(hotel.ID.String()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "hotel_id", "info_type", "info_content"}).
			AddRow(hotel.ContactInfos[0].ID.String(), hotel.ID.String(), hotel.ContactInfos[0].InfoType, hotel.ContactInfos[0].InfoContent))

	// Test GetHotelDetails method
	result, err := repo.GetHotelDetails(hotel.ID)
	assert.NoError(t, err)
	assert.Equal(t, hotel.ID, result.ID)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}
