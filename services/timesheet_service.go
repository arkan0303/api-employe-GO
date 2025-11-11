package services

import (
	"api-rect-go/db"
	"fmt"

	"gorm.io/gorm"
)



func NewTimesheetService(db *gorm.DB) *TimesheetService {
	return &TimesheetService{DB: db}
}

func GetMergedDatas(bulan int, tahun int, idCutOff int) ([]map[string]interface{}, error) {
query := fmt.Sprintf(`
	SELECT 
		td.*, 
		td.id AS id_timesheets,
		ct.id AS id_cuti_timesheets,
		ct.fullname, 
		ct.company_name, 
		ct.created_at AS created_at_cuti, 
		CASE 
			WHEN ct.photo_timesheets1 IS NOT NULL 
				AND ct.tahun = '%d' 
				AND ct.bulan = '%d'
			THEN ct.photo_timesheets1 
			ELSE NULL 
		END AS photo_timesheets1,
		CASE 
			WHEN ct.photo_timesheets2 IS NOT NULL 
				AND ct.tahun = '%d' 
				AND ct.bulan = '%d'
			THEN ct.photo_timesheets2 
			ELSE NULL 
		END AS photo_timesheets2,
		ct.bulan, 
		ct.tahun, 
		ct.start_cutoff, 
		ct.end_cutoff, 
		ct.tanggal_gajian, 
		ct.periode_timesheets,
		CASE 
			WHEN ct.photo_timesheets1 IS NOT NULL 
				AND ct.end_cutoff IS NOT NULL  
				AND ct.tahun = '%d' 
				AND ct.bulan = '%d'
			THEN 1
			ELSE 0
		END AS status,
		co.periode
	FROM 
		public.timesheet_driver td
	LEFT JOIN 
		public."Cuti_TImesheets" ct 
		ON td.employee_id = ct.employee_id
	LEFT JOIN 
		public."cut_off" co  
		ON td.id_cut_off = co.id
	WHERE 
		(ct.bulan = '%d' AND ct.tahun = '%d' AND td.id_cut_off = %d AND td.status_data IS NULL)
		OR (td.id_cut_off = %d AND td.status_data IS NULL)
	ORDER BY 
		ct.id DESC
`, tahun, bulan, tahun, bulan, tahun, bulan, bulan, tahun, idCutOff, idCutOff)


	var results []map[string]interface{}
	rows, err := db.DB.Raw(query).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	for rows.Next() {
		colsData := make([]interface{}, len(cols))
		colsDataPtrs := make([]interface{}, len(cols))
		for i := range colsData {
			colsDataPtrs[i] = &colsData[i]
		}

		if err := rows.Scan(colsDataPtrs...); err != nil {
			return nil, err
		}

		rowMap := make(map[string]interface{})
		for i, colName := range cols {
			val := colsData[i]
			rowMap[colName] = val
		}

		results = append(results, rowMap)
	}

	return results, nil
}
