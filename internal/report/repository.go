package report

import (
	"errors"
	"fmt"

	"hotel-guide/internal/db"
	"hotel-guide/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ReportRepository arayüzü, rapor veritabanı işlemlerini tanımlar
type ReportRepository interface {
	Save(report *models.Report) error
	ListReports() ([]models.Report, error)
	GetReportByID(id uuid.UUID) (*models.Report, error)
	UpdateReportStatus(id uuid.UUID, status models.ReportStatus) error
	GetHotelStatsByLocation(location string, hotels *[]models.Hotel) error
	UpdateReportStats(reportID uuid.UUID, hotelCount, phoneCount int, status models.ReportStatus) error
}

// reportRepository struct'ı, ReportRepository arayüzünü uygular
type reportRepository struct {
	db *gorm.DB
}

// NewReportRepository, yeni bir reportRepository örneği döndürür
func NewRepository() ReportRepository {
	return &reportRepository{db: db.DB}
}

// Save yeni bir rapor kaydeder
func (r *reportRepository) Save(report *models.Report) error {
	return r.db.Create(report).Error
}

// ListReports tüm raporları listeler
func (r *reportRepository) ListReports() ([]models.Report, error) {
	var reports []models.Report
	err := r.db.Find(&reports).Error
	return reports, err
}

// GetReportByID belirli bir ID'ye sahip raporu getirir
func (r *reportRepository) GetReportByID(id uuid.UUID) (*models.Report, error) {
	var report models.Report
	err := r.db.First(&report, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // Rapor bulunamazsa nil döndür
	}
	return &report, err
}

// UpdateReportStatus raporun durumunu günceller
func (r *reportRepository) UpdateReportStatus(id uuid.UUID, status models.ReportStatus) error {
	return r.db.Model(&models.Report{}).Where("id = ?", id).Update("status", status).Error
}

func (r *reportRepository) GetHotelStatsByLocation(location string, hotels *[]models.Hotel) error {
	// ContactInfo'dan location bilgilerini sorguluyoruz
	err := r.db.Joins("JOIN contact_infos ON contact_infos.hotel_id = hotels.id").
		Where("contact_infos.info_type = ? AND contact_infos.info_content = ?", "location", location).
		Find(hotels).Error

	if err != nil {
		return fmt.Errorf("failed to get hotels by location: %w", err)
	}

	return nil
}

func (r *reportRepository) UpdateReportStats(reportID uuid.UUID, hotelCount, phoneCount int, status models.ReportStatus) error {
	return r.db.Model(&models.Report{}).
		Where("id = ?", reportID).
		Updates(map[string]interface{}{
			"hotel_count": hotelCount,
			"phone_count": phoneCount,
			"status":      status,
		}).Error
}
