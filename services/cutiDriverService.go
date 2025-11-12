package services

import (
	"fmt"
	"time"

	"api-rect-go/db"
)

type ReplacementData struct {
	ID                     int64      `json:"id"`
	NamaDriverCuti         string     `json:"nama_driver_cuti"`
	NamaDriverPengganti    string     `json:"nama_driver_pengganti"`
	TypeReplacement        string     `json:"type_replacement"`
	StartDate              *time.Time `json:"start_date"`
	EndDate                *time.Time `json:"end_date"`
	Reason                 string     `json:"reason"`
	DriverPengganti        string     `json:"driver_pengganti"`
	NamaPerusahaan         string     `json:"nama_perusahaan"`
	EmployeeIDCuti      string `json:"employee_id_driver_cuti" gorm:"column:employee_id_cuti"`
EmployeeIDPengganti string `json:"employee_id_driver_pengganti" gorm:"column:employee_id_pengganti"`
	NoRekeningPengganti    string     `json:"no_rekening_pengganti"`
	NamaRekeningPengganti  string     `json:"nama_rekening_pengganti"`
	NamaBankPengganti      string     `json:"nama_bank_pengganti"`
	TanggalGajian          *string    `json:"tanggal_gajian"`
	NamaPerusahaanInternal *string    `json:"nama_perusahaan_internal"`
	NoHPCuti               string     `json:"no_hp_cuti"`
	NoHPPengganti          string     `json:"no_hp_pengganti"`
}

func GetServiceReplacementData() ([]ReplacementData, error) {
	dbMySQL := db.DBMySQL

	var results []ReplacementData

	query := `
	SELECT 
		sr.id AS id,
		sdp.full_name AS nama_driver_pengganti,
		sdc.full_name AS nama_driver_cuti,
		CASE WHEN sr.type_replacement = 0 THEN 'Temporary Replacement' ELSE 'Permanent Replacement' END AS type_replacement,
		 CONVERT_TZ(FROM_UNIXTIME(sr.start_date), '+00:00', '+07:00') AS start_date,
    CONVERT_TZ(FROM_UNIXTIME(sr.end_date), '+00:00', '+07:00') AS end_date,
		sr.reason,
		sr.driver_pengganti,
		mc.company_name AS nama_perusahaan,
	sdc.employee_id AS employee_id_cuti,
sdp.employee_id AS employee_id_pengganti,
		sdp.no_rekening AS no_rekening_pengganti,
		sdp.nama_rekening AS nama_rekening_pengganti,
		sdp.nama_bank AS nama_bank_pengganti,
		sc.salary_date AS tanggal_gajian,
		tic.company_name AS nama_perusahaan_internal,
		sdc.phone_number AS no_hp_cuti,
		sdp.phone_number AS no_hp_pengganti
	FROM service_replacements sr
	LEFT JOIN service_drivers sdp ON sr.service_drivers_id = sdp.id
	LEFT JOIN service_drivers sdc ON sdc.id = sr.service_details_id
	LEFT JOIN service_details sd ON sd.service_drivers_id = sr.service_details_id
	LEFT JOIN service_users su ON su.id = sd.service_users_id
	LEFT JOIN iwo_templates iwo ON iwo.id = su.iwo_templates_id
	LEFT JOIN salary_cut_offs sc ON sc.id = iwo.tanggal_gajian
	LEFT JOIN master_companies mc ON mc.id = su.master_companies_id
	LEFT JOIN tb_internal_companies tic ON tic.id = mc.internal_company_id
	WHERE sdc.full_name IS NOT NULL 
	  AND sr.status = 1
	ORDER BY sr.id DESC;
	`

	err := dbMySQL.Raw(query).Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("gagal ambil data: %v", err)
	}

	// Format tambahan untuk NoHPPenggantiWA
	for i := range results {
		noHp := results[i].NoHPPengganti
		if len(noHp) > 0 && noHp[0] == '0' {
			results[i].NoHPPengganti = "62" + noHp[1:]
		}
	}

	return results, nil
}

// DeleteServiceReplacement performs a soft delete by setting status to 0
func DeleteServiceReplacement(id int64) error {
    dbMySQL := db.DBMySQL

    // First check if the record exists
    var count int64
    err := dbMySQL.Table("service_replacements").Where("id = ?", id).Count(&count).Error
    if err != nil {
        return fmt.Errorf("failed to check record existence: %v", err)
    }

    if count == 0 {
        return fmt.Errorf("record with id %d not found", id)
    }

    // Soft delete by updating status to 0 instead of actually deleting
    err = dbMySQL.Exec("UPDATE service_replacements SET status = 0 WHERE id = ?", id).Error
    if err != nil {
        return fmt.Errorf("failed to delete service replacement: %v", err)
    }

    return nil
}
