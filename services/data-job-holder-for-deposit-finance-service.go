package services

import (
	"api-rect-go/db"
	models "api-rect-go/modals/mysql"
	"context"
)

type MASTER struct {
	ID              int32  `json:"id"`
	Nama            string `json:"nama"`
	StatusKaryawan  string `json:"status_karyawan"`
	Foto            string `json:"foto"`
	// Form1ID         int32  `json:"form_1_id"`
	EmployeeID      string `json:"employee_id"`

	// job holder
	// JobHolderID      int32  `json:"job_holder_id"`
	TglJobHolder     string `json:"tgl_job_holder"`
	IDCustomer       int32  `json:"id_customer"`
	IDStatusDataDiri int32  `json:"id_status_data_diri"`
	JobEmployeeID    string `json:"job_employee_id"`

	// mysql company
	CompanyName string `json:"company_name"`

	//status diri 
	Status           string    `gorm:"column:status" json:"status"`
	IDMasterDataDiri int32     `gorm:"column:id_master_data_diri;not null" json:"id_master_data_diri"`
	Tanggal          string    `gorm:"column:tanggal" json:"tanggal"`
}

func DataJobHolderForDepositFinance() ([]MASTER, error) {
	var jobHolders []MASTER
	query := `
-- Status NAIK dari tb_job_holder
SELECT 
    m.id,
    m.nama,
    m.status_karyawan,
    m.foto,
    m.employee_id,
    'Naik' AS status,
    j.tgl_job_holder AS tanggal,
    m.id AS id_master_data_diri,
    j.tgl_job_holder,
    j.id_customer,
    j.id_status_data_diri,
    j.employee_id AS job_employee_id
FROM tb_master_data_diri m
JOIN tb_job_holder j 
    ON m.id = j.id_status_data_diri 
    OR m.employee_id = j.employee_id
WHERE j.tgl_job_holder IS NOT NULL

UNION ALL

-- Status TURUN dari tb_status_diri
SELECT 
    m.id,
    m.nama,
    m.status_karyawan,
    m.foto,
    m.employee_id,
    sd.status,
    sd.tanggal,
    m.id AS id_master_data_diri,
    j.tgl_job_holder,
    j.id_customer,
    j.id_status_data_diri,
    j.employee_id AS job_employee_id
FROM tb_master_data_diri m
JOIN tb_job_holder j 
    ON m.id = j.id_status_data_diri 
    OR m.employee_id = j.employee_id
JOIN tb_status_diri sd
    ON sd.id_master_data_diri = m.id
WHERE sd.status = 'Turun'
  AND sd.tanggal IS NOT NULL
  AND j.tgl_job_holder IS NOT NULL
  AND sd.tanggal > j.tgl_job_holder

ORDER BY id, tanggal ASC;
    `

	if err := db.DB.Raw(query).Scan(&jobHolders).Error; err != nil {
		return nil, err
	}

	if len(jobHolders) == 0 {
		return jobHolders, nil
	}

	customerIDs := make([]int32, 0, len(jobHolders))
	customerIDMap := make(map[int32]bool)

	for _, item := range jobHolders {
		if item.IDCustomer > 0 && !customerIDMap[item.IDCustomer] {
			customerIDs = append(customerIDs, item.IDCustomer)
			customerIDMap[item.IDCustomer] = true
		}
	}

	var companies []models.MasterCompany
	companyMap := make(map[int32]string)

	if len(customerIDs) > 0 {
		err := db.DBMySQL.
			Select("id, company_name").
			Where("id IN ?", customerIDs).
			Find(&companies).Error

		if err != nil {
			return nil, err
		}

		for _, comp := range companies {
			companyMap[comp.ID] = comp.CompanyName
		}
	}

	for i := range jobHolders {
		if name, exists := companyMap[jobHolders[i].IDCustomer]; exists {
			jobHolders[i].CompanyName = name
		} else {
			jobHolders[i].CompanyName = ""
		}
	}

	return jobHolders, nil
}
// BONUS: Versi dengan Context untuk timeout & cancellation
func DataJobHolderForDepositFinanceWithContext(ctx context.Context) ([]MASTER, error) {
	var jobHolders []MASTER

	query := `
        SELECT 
            m.id,
            m.nama,
            m.status_karyawan,
            m.foto,
            m.form_1_id,
            m.employee_id,
            j.id AS job_holder_id,
            j.tgl_job_holder,
            j.id_customer,
            j.id_status_data_diri,
            j.employee_id AS job_employee_id
        FROM tb_master_data_diri m
        JOIN tb_job_holder j 
            ON m.id = j.id_status_data_diri 
            OR m.employee_id = j.employee_id

    `

	if err := db.DB.WithContext(ctx).Raw(query).Scan(&jobHolders).Error; err != nil {
		return nil, err
	}

	if len(jobHolders) == 0 {
		return jobHolders, nil
	}

	customerIDs := make([]int32, 0, len(jobHolders))
	customerIDMap := make(map[int32]bool)

	for _, item := range jobHolders {
		if item.IDCustomer > 0 && !customerIDMap[item.IDCustomer] {
			customerIDs = append(customerIDs, item.IDCustomer)
			customerIDMap[item.IDCustomer] = true
		}
	}

	var companies []models.MasterCompany
	companyMap := make(map[int32]string)

	if len(customerIDs) > 0 {
		err := db.DBMySQL.WithContext(ctx).
			Select("id, company_name").
			Where("id IN ?", customerIDs).
			Find(&companies).Error

		if err != nil {
			return nil, err
		}

		for _, comp := range companies {
			companyMap[comp.ID] = comp.CompanyName
		}
	}

	for i := range jobHolders {
		if name, exists := companyMap[jobHolders[i].IDCustomer]; exists {
			jobHolders[i].CompanyName = name
		} else {
			jobHolders[i].CompanyName = ""
		}
	}

	return jobHolders, nil
}