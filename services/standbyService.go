package services

import (
	"api-rect-go/db"
	models "api-rect-go/modals"
	"api-rect-go/modals/mysql"
	"fmt"

	// "strings"
	"sync"
)

type MasterRingkas struct {
	ID         int32  `json:"id" gorm:"column:id;primaryKey;autoIncrement:true"`
	Nama       string `json:"nama" gorm:"column:nama"`
	EmployeeID string `json:"employee_id" gorm:"column:employee_id"`
	NoTelp     string `json:"no_telp" gorm:"column:no_telp"`
	Form1ID    int32  `json:"form_1_id" gorm:"column:form_1_id"`
	Foto       string `json:"foto" gorm:"column:foto"`
	TglLahir   string `gorm:"column:tgl_lahir" json:"tgl_lahir"`
	StatusKawin string `gorm:"column:status_kawin" json:"status_kawin"`
}

type FormGabungan struct {
	Master                    MasterRingkas `json:"master"`
	PengalamanJepang          bool          `json:"pengalaman_jepang"`
	PendidikanTerakhir        string        `json:"pendidikan_terakhir"`
	Agama                     string        `json:"agama"`
	SimANomor                 string        `json:"sim_a_nomor"`
	SimAMasaBerlaku           string        `json:"sim_a_masa_berlaku"`
	SimB1Nomor                string        `json:"sim_b1_nomor"`
	SimB1MasaBerlaku          string        `json:"sim_b1_masa_berlaku"`
	SimB2Nomor                string        `json:"sim_b2_nomor"`
	SimB2MasaBerlaku          string        `json:"sim_b2_masa_berlaku"`
	PengalamanTahunMulai      int32         `json:"pengalaman_tahun_mulai"`
	PengalamanTahunSelesai    int32         `json:"pengalaman_tahun_selesai"`
	PengalamanNamaPerusahaan  string        `json:"pengalaman_nama_perusahaan"`
	PengalamanUserEkspat      string        `json:"pengalaman_user_ekspat"`
	PengalamanNegaraAsal      string        `json:"pengalaman_negara_asal"`
	PengalamanTahunMulai2     int32         `json:"pengalaman_tahun_mulai2"`
	PengalamanTahunSelesai2   int32         `json:"pengalaman_tahun_selesai2"`
	PengalamanNamaPerusahaan2 string        `json:"pengalaman_nama_perusahaan2"`
	PengalamanUserEkspat2     string        `json:"pengalaman_user_ekspat2"`
	PengalamanNegaraAsal2     string        `json:"pengalaman_negara_asal2"`
	PengalamanTahunMulai3     int32         `json:"pengalaman_tahun_mulai3"`
	PengalamanTahunSelesai3   int32         `json:"pengalaman_tahun_selesai3"`
	PengalamanNamaPerusahaan3 string        `json:"pengalaman_nama_perusahaan3"`
	PengalamanUserEkspat3     string        `json:"pengalaman_user_ekspat3"`
	PengalamanNegaraAsal3     string        `json:"pengalaman_negara_asal3"`
	NoKtp                     string        `gorm:"column:no_ktp" json:"no_ktp"`
	NoNpwp                    string        `gorm:"column:no_npwp" json:"no_npwp"`
	Pertanyaan6               string        `json:"pertanyaan_6" gorm:"column:pertanyaan_6;not null"`
	Kerapihan                 string        `json:"kerapihan" gorm:"column:kerapihan;not null"`
	KemampuanBahasaInggris    string        `json:"kemampuanBahasaInggris" gorm:"column:kemampuanBahasaInggris;not null"`
	StatusRecruitment         string        `json:"status_recruitment"`
}

type form2Data struct {
	Agama              string
	PendidikanTerakhir string
	SimANomor          string
	SimAMasaBerlaku    string
	SimB1Nomor         string
	SimB1MasaBerlaku   string
	SimB2Nomor         string
	SimB2MasaBerlaku   string
	PengalamanData     [3]pengalamanInfo
	NoKtp              string
	NoNpwp             string
}

type pengalamanInfo struct {
	TahunMulai     int32
	TahunSelesai   int32
	NamaPerusahaan string
	UserEkspat     string
	NegaraAsal     string
}

// UpdateMasterData updates the master data with all related form data
// Including photo upload support
func UpdateMasterData(id int32, updateData map[string]interface{}) error {
	// 1. Handle photo upload if exists
	if foto, ok := updateData["foto"].(string); ok && foto != "" {
		// If the photo is a base64 string, you can handle the upload here
		// For example, you might want to save it to a file or cloud storage
		// and update the foto field with the URL or file path
		// This is a placeholder - implement according to your storage solution
		// fotoURL, err := savePhotoToStorage(foto)
		// if err != nil {
		//     return fmt.Errorf("gagal mengunggah foto: %w", err)
		// }
		// updateData["foto"] = fotoURL

		// For now, we'll just pass through the photo value as is
		// Make sure your client sends the correct URL or file path
	}

	// Start a transaction
	tx := db.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 1. Ambil data master untuk mendapatkan form_1_id
	var masterData models.TbMasterDataDiri
	if err := tx.Model(&models.TbMasterDataDiri{}).
		Where("id = ?", id).
		First(&masterData).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal mengambil data master: %w", err)
	}

	// 2. Update TbMasterDataDiri (master data)
	masterFields := map[string]bool{
		"nama": true, "employee_id": true, "no_telp": true, "form_1_id": true,
		"foto": true, "tgl_lahir": true, "status_kawin": true, "status_karyawan": true,
	}
	masterUpdate := make(map[string]interface{})
	for key, value := range updateData {
		if masterFields[key] {
			masterUpdate[key] = value
		}
	}

	if len(masterUpdate) > 0 {
		if err := tx.Model(&models.TbMasterDataDiri{}).
			Where("id = ?", id).
			Updates(masterUpdate).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal memperbarui data master: %w", err)
		}
	}

	// 2. Update Form1 data (jika ada form_1_id)
	if form1ID, ok := updateData["form_1_id"].(float64); ok && form1ID > 0 {
		form1Update := make(map[string]interface{})
		form1Fields := map[string]bool{
			"pengalaman_jepang": true,
		}

		for key, value := range updateData {
			if form1Fields[key] {
				form1Update[key] = value
			}
		}

		if len(form1Update) > 0 {
			if err := tx.Model(&models.Form1{}).
				Where("id = ?", int32(form1ID)).
				Updates(form1Update).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("gagal memperbarui form 1: %w", err)
			}
		}
	}

	// 3. Update Form2 data
	// Form2 fields will be handled in the pengalaman kerja section below

	// Handle pengalaman kerja (bisa multiple)
	for i := 1; i <= 3; i++ {
		prefix := fmt.Sprintf("pengalaman_%d_", i)
		pengalamanFields := map[string]string{
			"nama_perusahaan": fmt.Sprintf("pengalaman_nama_perusahaan%d", i),
			"negara_asal":     fmt.Sprintf("pengalaman_negara_asal%d", i),
			"tahun_mulai":     fmt.Sprintf("pengalaman_tahun_mulai%d", i),
			"tahun_selesai":   fmt.Sprintf("pengalaman_tahun_selesai%d", i),
			"user_ekspat":     fmt.Sprintf("pengalaman_user_ekspat%d", i),
		}

		hasPengalaman := false
		pengalamanData := make(map[string]interface{})

		for field, dbField := range pengalamanFields {
			key := prefix + field
			if value, exists := updateData[key]; exists {
				hasPengalaman = true
				pengalamanData[dbField] = value
			}
		}

		if hasPengalaman {
			// Update or create form2 record
			form1ID := masterData.Form1ID // Menggunakan form_1_id dari data master
			if form1ID == 0 {
				tx.Rollback()
				return fmt.Errorf("form_1_id tidak valid untuk data master dengan id %d", id)
			}
			pengalamanData["form_1_id"] = form1ID
			
			// Cek apakah record sudah ada
			var existingForm2 models.Form2
			result := tx.Model(&models.Form2{}).
				Where("form_1_id = ?", form1ID).
				First(&existingForm2)
			
			if result.Error != nil && result.Error.Error() != "record not found" {
				tx.Rollback()
				return fmt.Errorf("gagal memeriksa data form2: %w", result.Error)
			}
			
			if result.RowsAffected > 0 {
				// Update existing record
				if err := tx.Model(&models.Form2{}).
					Where("id = ?", existingForm2.ID).
					Updates(pengalamanData).Error; err != nil {
					tx.Rollback()
					return fmt.Errorf("gagal memperbarui data form2: %w", err)
				}
			} else {
				// Create new record
				if err := tx.Model(&models.Form2{}).
					Create(&pengalamanData).Error; err != nil {
					tx.Rollback()
					return fmt.Errorf("gagal membuat data form2 baru: %w", err)
				}
			}
		}
	}

	// 4. Update Form6 data
	form6Update := make(map[string]interface{})
	form6Fields := map[string]bool{
		"pertanyaan_6": true, "kerapihan": true, "kemampuan_bahasa_inggris": true,
	}

	for key, value := range updateData {
		if form6Fields[key] {
			form6Update[key] = value
		}
	}

	if len(form6Update) > 0 {
		// Gunakan form_1_id yang sama dengan yang digunakan di Form2
		if err := tx.Model(&models.Form6{}).
			Where("form_1_id = ?", masterData.Form1ID).
			Updates(form6Update).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal memperbarui form 6: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("gagal melakukan commit transaksi: %w", err)
	}

	return nil
}

// fetchRecruitmentIDs dengan query yang lebih efisien
func fetchRecruitmentIDs() (map[int32]bool, error) {
	var recruitmentIDs []int32

	if err := db.DBMySQL.Model(&mysql.ServiceDriver{}).
		Select("DISTINCT id_recruitment").
		Where("id_recruitment IS NOT NULL AND id_recruitment != 0").
		Pluck("id_recruitment", &recruitmentIDs).Error; err != nil {
		return nil, fmt.Errorf("gagal mengambil data recruitment: %w", err)
	}

	// Pre-allocate map dengan kapasitas yang tepat
	recruitmentMap := make(map[int32]bool, len(recruitmentIDs))
	for _, id := range recruitmentIDs {
		recruitmentMap[id] = true
	}

	return recruitmentMap, nil
}

func GetMasterDataAvailableWithForms() ([]FormGabungan, error) {
	// Struktur untuk menampung hasil parallel queries
	type queryResults struct {
		masters        []MasterRingkas
		recruitmentMap map[int32]bool
		form1Map       map[int32]bool
		form2Map       map[int32]form2Data
		form6Map       map[int32]struct {
			Kerapihan              string
			KemampuanBahasaInggris string
			Pertanyaan6            string
		}
		err error
	}

	resultChan := make(chan queryResults, 1)

	go func() {
		var result queryResults
		var wg sync.WaitGroup
		
		// Mutex untuk thread-safe map writes
		var mu sync.Mutex
		errorList := make([]error, 0, 4)

		// Step 1: Fetch recruitment IDs dan master data secara parallel
		wg.Add(2)

		// Goroutine 1: Fetch recruitment IDs
		go func() {
			defer wg.Done()
			recMap, err := fetchRecruitmentIDs()
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				errorList = append(errorList, err)
				result.recruitmentMap = make(map[int32]bool)
			} else {
				result.recruitmentMap = recMap
			}
		}()

		// Goroutine 2: Fetch master data
		go func() {
			defer wg.Done()
			var masters []MasterRingkas
			err := db.DB.Model(&models.TbMasterDataDiri{}).
				Select("id", "nama", "employee_id", "no_telp", "form_1_id", "foto", "tgl_lahir", "status_kawin").
				Where("status_karyawan IN (?, ?)", "Temporary", "standby").
				Scan(&masters).Error
			
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				errorList = append(errorList, err)
			} else {
				result.masters = masters
			}
		}()

		wg.Wait()

		// Cek error di step 1
		if len(errorList) > 0 {
			result.err = errorList[0]
			resultChan <- result
			return
		}

		if len(result.masters) == 0 {
			resultChan <- result
			return
		}

		// Siapkan form1IDs dengan pre-allocated capacity
		form1IDs := make([]int32, 0, len(result.masters))
		for i := range result.masters {
			form1IDs = append(form1IDs, result.masters[i].Form1ID)
		}

		// Step 2: Fetch all forms secara parallel
		wg.Add(3)
		errorList = errorList[:0] // Reset error list

		// Query Form1 - Optimized
		go func() {
			defer wg.Done()
			var form1Results []struct {
				ID               int32 `gorm:"column:id"`
				PengalamanJepang bool  `gorm:"column:pengalaman_jepang"`
			}
			
			err := db.DB.Model(&models.Form1{}).
				Select("id", "pengalaman_jepang").
				Where("id IN ?", form1IDs).
				Scan(&form1Results).Error
			
			if err != nil {
				mu.Lock()
				errorList = append(errorList, err)
				mu.Unlock()
				return
			}

			// Build map tanpa lock di setiap iterasi
			localMap := make(map[int32]bool, len(form1Results))
			for i := range form1Results {
				localMap[form1Results[i].ID] = form1Results[i].PengalamanJepang
			}
			
			mu.Lock()
			result.form1Map = localMap
			mu.Unlock()
		}()

		// Query Form2 - Optimized
		go func() {
			defer wg.Done()
			var form2Results []struct {
				Form1ID                   int32  `gorm:"column:form_1_id"`
				Agama                     string `gorm:"column:agama"`
				PendidikanTerakhir        string `gorm:"column:pendidikan_terakhir"`
				SimANomor                 string `gorm:"column:sim_a_nomor"`
				SimAMasaBerlaku           string `gorm:"column:sim_a_masa_berlaku"`
				SimB1Nomor                string `gorm:"column:sim_b1_nomor"`
				SimB1MasaBerlaku          string `gorm:"column:sim_b1_masa_berlaku"`
				SimB2Nomor                string `gorm:"column:sim_b2_nomor"`
				SimB2MasaBerlaku          string `gorm:"column:sim_b2_masa_berlaku"`
				PengalamanTahunMulai      int32  `gorm:"column:pengalaman_tahun_mulai"`
				PengalamanTahunSelesai    int32  `gorm:"column:pengalaman_tahun_selesai"`
				PengalamanNamaPerusahaan  string `gorm:"column:pengalaman_nama_perusahaan"`
				PengalamanUserEkspat      string `gorm:"column:pengalaman_user_ekspat"`
				PengalamanNegaraAsal      string `gorm:"column:pengalaman_negara_asal"`
				PengalamanTahunMulai2     int32  `gorm:"column:pengalaman_tahun_mulai2"`
				PengalamanTahunSelesai2   int32  `gorm:"column:pengalaman_tahun_selesai2"`
				PengalamanNamaPerusahaan2 string `gorm:"column:pengalaman_nama_perusahaan2"`
				PengalamanUserEkspat2     string `gorm:"column:pengalaman_user_ekspat2"`
				PengalamanNegaraAsal2     string `gorm:"column:pengalaman_negara_asal2"`
				PengalamanTahunMulai3     int32  `gorm:"column:pengalaman_tahun_mulai3"`
				PengalamanTahunSelesai3   int32  `gorm:"column:pengalaman_tahun_selesai3"`
				PengalamanNamaPerusahaan3 string `gorm:"column:pengalaman_nama_perusahaan3"`
				PengalamanUserEkspat3     string `gorm:"column:pengalaman_user_ekspat3"`
				PengalamanNegaraAsal3     string `gorm:"column:pengalaman_negara_asal3"`
				NoKtp                     string `gorm:"column:no_ktp"`
				NoNpwp                    string `gorm:"column:no_npwp"`
			}
			
			err := db.DB.Model(&models.Form2{}).
				Select("form_1_id", "agama", "pendidikan_terakhir", "sim_a_nomor", "sim_a_masa_berlaku",
					"sim_b1_nomor", "sim_b1_masa_berlaku", "sim_b2_nomor", "sim_b2_masa_berlaku",
					"pengalaman_tahun_mulai", "pengalaman_tahun_selesai", "pengalaman_nama_perusahaan",
					"pengalaman_user_ekspat", "pengalaman_negara_asal", "pengalaman_tahun_mulai2",
					"pengalaman_tahun_selesai2", "pengalaman_nama_perusahaan2", "pengalaman_user_ekspat2",
					"pengalaman_negara_asal2", "pengalaman_tahun_mulai3", "pengalaman_tahun_selesai3",
					"pengalaman_nama_perusahaan3", "pengalaman_user_ekspat3", "pengalaman_negara_asal3",
					"no_ktp", "no_npwp").
				Where("form_1_id IN ?", form1IDs).
				Scan(&form2Results).Error
			
			if err != nil {
				mu.Lock()
				errorList = append(errorList, err)
				mu.Unlock()
				return
			}

			localMap := make(map[int32]form2Data, len(form2Results))
			for i := range form2Results {
				f := &form2Results[i]
				localMap[f.Form1ID] = form2Data{
					Agama:              f.Agama,
					PendidikanTerakhir: f.PendidikanTerakhir,
					SimANomor:          f.SimANomor,
					SimAMasaBerlaku:    f.SimAMasaBerlaku,
					SimB1Nomor:         f.SimB1Nomor,
					SimB1MasaBerlaku:   f.SimB1MasaBerlaku,
					SimB2Nomor:         f.SimB2Nomor,
					SimB2MasaBerlaku:   f.SimB2MasaBerlaku,
					PengalamanData: [3]pengalamanInfo{
						{f.PengalamanTahunMulai, f.PengalamanTahunSelesai, f.PengalamanNamaPerusahaan, f.PengalamanUserEkspat, f.PengalamanNegaraAsal},
						{f.PengalamanTahunMulai2, f.PengalamanTahunSelesai2, f.PengalamanNamaPerusahaan2, f.PengalamanUserEkspat2, f.PengalamanNegaraAsal2},
						{f.PengalamanTahunMulai3, f.PengalamanTahunSelesai3, f.PengalamanNamaPerusahaan3, f.PengalamanUserEkspat3, f.PengalamanNegaraAsal3},
					},
					NoKtp:  f.NoKtp,
					NoNpwp: f.NoNpwp,
				}
			}
			
			mu.Lock()
			result.form2Map = localMap
			mu.Unlock()
		}()

		// Query Form6 - Optimized
		go func() {
			defer wg.Done()
			var form6Results []struct {
				Form1ID                int32  `gorm:"column:form_1_id"`
				Kerapihan              string `gorm:"column:kerapihan"`
				KemampuanBahasaInggris string `gorm:"column:kemampuanBahasaInggris"`
				Pertanyaan6            string `gorm:"column:pertanyaan_6"`
			}
			
			err := db.DB.Model(&models.Form6{}).
				Select("form_1_id", "kerapihan", "kemampuanBahasaInggris", "pertanyaan_6").
				Where("form_1_id IN ?", form1IDs).
				Scan(&form6Results).Error
			
			if err != nil {
				mu.Lock()
				errorList = append(errorList, err)
				mu.Unlock()
				return
			}

			localMap := make(map[int32]struct {
				Kerapihan              string
				KemampuanBahasaInggris string
				Pertanyaan6            string
			}, len(form6Results))
			
			for i := range form6Results {
				f := &form6Results[i]
				localMap[f.Form1ID] = struct {
					Kerapihan              string
					KemampuanBahasaInggris string
					Pertanyaan6            string
				}{f.Kerapihan, f.KemampuanBahasaInggris, f.Pertanyaan6}
			}
			
			mu.Lock()
			result.form6Map = localMap
			mu.Unlock()
		}()

		wg.Wait()

		if len(errorList) > 0 {
			result.err = errorList[0]
		}

		resultChan <- result
	}()

	// Wait for results
	queryResult := <-resultChan

	if queryResult.err != nil {
		return nil, queryResult.err
	}

	if len(queryResult.masters) == 0 {
		return []FormGabungan{}, nil
	}

	// Step 3: Gabungkan data dengan pre-allocated slice
	results := make([]FormGabungan, 0, len(queryResult.masters))

	for i := range queryResult.masters {
		master := &queryResult.masters[i]
		form1 := queryResult.form1Map[master.Form1ID]
		form2 := queryResult.form2Map[master.Form1ID]
		form6 := queryResult.form6Map[master.Form1ID]

		// Determine recruitment status (inline tanpa variable temp)
		statusRecruitment := "belum_ada"
		if queryResult.recruitmentMap[master.ID] {
			statusRecruitment = "sudah_ada"
		}

		results = append(results, FormGabungan{
			Master:                    *master,
			PengalamanJepang:          form1,
			Agama:                     form2.Agama,
			PendidikanTerakhir:        form2.PendidikanTerakhir,
			SimANomor:                 form2.SimANomor,
			SimAMasaBerlaku:           form2.SimAMasaBerlaku,
			SimB1Nomor:                form2.SimB1Nomor,
			SimB1MasaBerlaku:          form2.SimB1MasaBerlaku,
			SimB2Nomor:                form2.SimB2Nomor,
			SimB2MasaBerlaku:          form2.SimB2MasaBerlaku,
			PengalamanTahunMulai:      form2.PengalamanData[0].TahunMulai,
			PengalamanTahunSelesai:    form2.PengalamanData[0].TahunSelesai,
			PengalamanNamaPerusahaan:  form2.PengalamanData[0].NamaPerusahaan,
			PengalamanUserEkspat:      form2.PengalamanData[0].UserEkspat,
			PengalamanNegaraAsal:      form2.PengalamanData[0].NegaraAsal,
			PengalamanTahunMulai2:     form2.PengalamanData[1].TahunMulai,
			PengalamanTahunSelesai2:   form2.PengalamanData[1].TahunSelesai,
			PengalamanNamaPerusahaan2: form2.PengalamanData[1].NamaPerusahaan,
			PengalamanUserEkspat2:     form2.PengalamanData[1].UserEkspat,
			PengalamanNegaraAsal2:     form2.PengalamanData[1].NegaraAsal,
			PengalamanTahunMulai3:     form2.PengalamanData[2].TahunMulai,
			PengalamanTahunSelesai3:   form2.PengalamanData[2].TahunSelesai,
			PengalamanNamaPerusahaan3: form2.PengalamanData[2].NamaPerusahaan,
			PengalamanUserEkspat3:     form2.PengalamanData[2].UserEkspat,
			PengalamanNegaraAsal3:     form2.PengalamanData[2].NegaraAsal,
			NoKtp:                     form2.NoKtp,
			NoNpwp:                    form2.NoNpwp,
			Kerapihan:                 form6.Kerapihan,
			KemampuanBahasaInggris:    form6.KemampuanBahasaInggris,
			Pertanyaan6:               form6.Pertanyaan6,
			StatusRecruitment:         statusRecruitment,
		})
	}

	return results, nil
}