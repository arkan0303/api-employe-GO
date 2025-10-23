package services

import (
	"api-rect-go/db"
	"api-rect-go/modals"

	"gorm.io/gorm"
)

type TimesheetService struct {
	DB *gorm.DB
}

type MergedData struct {
	ID                   int32   `json:"id"`
	EmployeeID           string  `json:"employee_id"`
	Perusahaan           string  `json:"perusahaan"`
	CustomerID           string  `json:"customer_id"`
	Customer             string  `json:"customer"`
	PembayaranGaji       string  `json:"pembayaran_gaji"`
	JumlahDriver         int     `json:"jumlah_driver"`
	Periode              *string `json:"periode"`
	TotalUploadTimesheet int     `json:"total_upload_timesheets"`
}

// GetMergedData menghasilkan data gabungan seperti query raw sebelumnya
func GetMergedData(bulan string, tahun int32) ([]MergedData, error) {
	var timesheets []modals.TimesheetDriver
	if err := db.DB.Order("id DESC").Find(&timesheets).Error; err != nil {
		return nil, err
	}

	// 1️⃣ Hitung jumlah driver per customer
	type DriverCount struct {
		Customer     string
		JumlahDriver int
	}
	var driverCounts []DriverCount
	db.DB.Table("timesheet_driver").
		Select("customer, COUNT(DISTINCT employee_id) as jumlah_driver").
		Group("customer").
		Scan(&driverCounts)

	driverCountMap := make(map[string]int)
	for _, d := range driverCounts {
		driverCountMap[d.Customer] = d.JumlahDriver
	}

	// 2️⃣ Hitung total upload timesheet per customer
	type UploadCount struct {
		Customer string
		Total    int
	}
	var uploadCounts []UploadCount
db.DB.Raw(`
    SELECT 
        ts.customer, 
        COUNT(DISTINCT ct.employee_id) as total
    FROM "Cuti_TImesheets" ct
    INNER JOIN timesheet_driver ts ON ts.employee_id = ct.employee_id
    WHERE 
        ct.photo_timesheets1 IS NOT NULL 
        AND TRIM(ct.photo_timesheets1) <> ''
        AND ct.bulan = ? 
        AND ct.tahun = ?
        AND ts.customer IS NOT NULL
        AND ts.customer <> ''
    GROUP BY ts.customer
`, bulan, tahun).Scan(&uploadCounts)


	uploadCountMap := make(map[string]int)
	for _, u := range uploadCounts {
		uploadCountMap[u.Customer] = u.Total
	}

	// 3️⃣ Ambil data cut_off (id -> periode)
	var cutOffs []modals.CutOff
	db.DB.Find(&cutOffs)
	cutOffMap := make(map[int32]string)
	for _, c := range cutOffs {
		cutOffMap[c.ID] = c.Periode
	}

	// 4️⃣ Gabungkan data (distinct customer)
	seen := make(map[string]bool)
	var merged []MergedData

	for _, ts := range timesheets {
		if seen[ts.Customer] {
			continue
		}
		seen[ts.Customer] = true

		var periode *string
		if val, ok := cutOffMap[ts.IDCutOff]; ok {
			periode = &val
		}

		merged = append(merged, MergedData{
			ID:                   ts.ID,
			EmployeeID:           ts.EmployeeID,
			Perusahaan:           ts.Perusahaan,
			CustomerID:           ts.CustomerID,
			Customer:             ts.Customer,
			PembayaranGaji:       ts.PembayaranGaji,
			JumlahDriver:         driverCountMap[ts.Customer],
			Periode:              periode,
			TotalUploadTimesheet: uploadCountMap[ts.Customer],
		})
	}

	return merged, nil
}
