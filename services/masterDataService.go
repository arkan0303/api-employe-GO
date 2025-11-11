package services

import (
	"api-rect-go/db"
	models "api-rect-go/modals"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type ExternalDriver struct {
	FullName   string `json:"full_name"`
	ID         int    `json:"id"`
	EmployeeID string `json:"employee_id"`
}

type ExternalResponse struct {
	Success bool              `json:"success"`
	Code    int               `json:"code"`
	Data    []ExternalDriver  `json:"data"`
}



type DataToForm2 struct {
	ID      int32  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	NoKtp   string `gorm:"column:no_ktp" json:"no_ktp"`
	NoNpwp  string `gorm:"column:no_npwp" json:"no_npwp"`
	Form1ID int32  `gorm:"column:form_1_id" json:"form_1_id"`
	Bank    string `gorm:"column:bank" json:"bank"`
	Norek   string `gorm:"column:norek" json:"norek"`
}

type MasterWithForm2 struct {
	models.TbMasterDataDiri
	Form2 DataToForm2 `json:"form2"`
}

func GetMasterData() ([]MasterWithForm2, error) {
	var masterData []models.TbMasterDataDiri

	// Ambil semua data master
	if err := db.DB.
		Where("status_karyawan = ?", "Available").
		Find(&masterData).Error; err != nil {
		return nil, err
	}

	// Kumpulkan semua form_1_id untuk query batch
	form1IDs := make([]int32, 0, len(masterData))
	for _, master := range masterData {
		form1IDs = append(form1IDs, master.Form1ID)
	}

	// Query semua form2 sekaligus (1 query saja!)
	var form2List []DataToForm2
	if len(form1IDs) > 0 {
		if err := db.DB.Model(&models.Form2{}).
			Select("id", "no_ktp", "no_npwp", "form_1_id", "bank", "norek").
			Where("form_1_id IN ?", form1IDs).
			Find(&form2List).Error; err != nil {
			return nil, err
		}
	}

	// Buat map untuk akses cepat O(1)
	form2Map := make(map[int32]DataToForm2, len(form2List))
	for _, f2 := range form2List {
		form2Map[f2.Form1ID] = f2
	}

	// Gabungkan hasil
	result := make([]MasterWithForm2, 0, len(masterData))
	for _, master := range masterData {
		result = append(result, MasterWithForm2{
			TbMasterDataDiri: master,
			Form2:            form2Map[master.Form1ID], // Akses O(1)
		})
	}

	return result, nil
}


func CreateMasterData(masterData *models.TbMasterDataDiri) error {
	result := db.DB.Create(masterData)
	return result.Error
}

func EditMasterData(id int) error {
	log.Printf("Memulai EditMasterData untuk ID: %d", id)
	// ===== 1. Generate employee_id baru dulu =====
	now := time.Now()
	datePart := now.Format("0601")
	log.Printf("Membuat employee_id dengan format tanggal: %s", datePart)

	var count int64
	query := db.DB.Model(&models.TbMasterDataDiri{}).
	    Where("employee_id LIKE ?", fmt.Sprintf("DR%s%%", datePart))
	
	log.Printf("Menjalankan query count: %v", query.Statement.SQL.String())
	
	if err := query.Count(&count).Error; err != nil {
	    log.Printf("Error saat menghitung employee_id: %v", err)
	    return fmt.Errorf("gagal menghitung employee_id: %v", err)
	}
	log.Printf("Jumlah employee_id yang ada dengan format DR%s: %d", datePart, count)

	// Fungsi untuk memeriksa ketersediaan employee_id di semua sumber
checkEmployeeID := func(employeeID string) (bool, error) {
    log.Printf("Memeriksa ketersediaan employee_id: %s", employeeID)
    
    // 1. Cek di tb_job_holder
    var existsJobHolder bool
    query := "SELECT EXISTS(SELECT 1 FROM tb_job_holder WHERE employee_id = ?)"
    log.Printf("Menjalankan query: %s dengan parameter: %s", query, employeeID)
    
    if err := db.DB.Raw(query, employeeID).Scan(&existsJobHolder).Error; err != nil {
        log.Printf("Error saat mengecek di tb_job_holder: %v", err)
        return false, fmt.Errorf("error mengecek di tb_job_holder: %v", err)
    }
    
    log.Printf("Hasil pengecekan di tb_job_holder untuk %s: %v", employeeID, existsJobHolder)
    if existsJobHolder {
        return false, nil
    }

    // 2. Cek di tb_master_data_diri
    var existsMaster bool
    query = "SELECT EXISTS(SELECT 1 FROM tb_master_data_diri WHERE employee_id = ?)"
    log.Printf("Menjalankan query: %s dengan parameter: %s", query, employeeID)
    
    if err := db.DB.Raw(query, employeeID).Scan(&existsMaster).Error; err != nil {
        log.Printf("Error saat mengecek di tb_master_data_diri: %v", err)
        return false, fmt.Errorf("error mengecek di tb_master_data_diri: %v", err)
    }
    
    log.Printf("Hasil pengecekan di tb_master_data_diri untuk %s: %v", employeeID, existsMaster)
    if existsMaster {
        return false, nil
    }

    // 3. Cek di API eksternal
    log.Printf("Memeriksa ketersediaan di API eksternal untuk: %s", employeeID)
    resp, err := http.Get("https://backend.sigapdriver.com/api/all_driver")
    if err != nil {
        log.Printf("Gagal memanggil API eksternal: %v", err)
        return false, fmt.Errorf("gagal memanggil API eksternal: %v", err)
    }
    defer resp.Body.Close()

    var ext ExternalResponse
    if err := json.NewDecoder(resp.Body).Decode(&ext); err != nil {
        log.Printf("Gagal decode response API: %v", err)
        return false, fmt.Errorf("gagal decode response API: %v", err)
    }

    for _, driver := range ext.Data {
        if driver.EmployeeID == employeeID {
            log.Printf("Employee ID %s ditemukan di API eksternal", employeeID)
            return false, nil
        }
    }

    log.Printf("Employee ID %s tersedia di semua sumber", employeeID)
    return true, nil
}
// Cari employee_id yang tersedia
var newEmployeeID string
for i := 1; i <= 1000; i++ { // Batasi maksimal 1000 percobaan untuk menghindari infinite loop
    candidateID := fmt.Sprintf("DR%s%03d", datePart, i)
    
    available, err := checkEmployeeID(candidateID)
    if err != nil {
        return fmt.Errorf("gagal memeriksa ketersediaan employee_id: %v", err)
    }
    
    if available {
        newEmployeeID = candidateID
        break
    }
    
    if i == 1000 {
        return fmt.Errorf("tidak dapat menemukan employee_id yang tersedia")
    }
}

	// ===== 5. Ambil data master dan form2 untuk payload API =====
	var masterData models.TbMasterDataDiri
	if err := db.DB.First(&masterData, id).Error; err != nil {
		return fmt.Errorf("gagal mengambil data master: %v", err)
	}

	var form2Data models.Form2 // Sesuaikan dengan nama model form2 Anda
	db.DB.Where("form_1_id = ?", masterData.Form1ID).First(&form2Data)

	// ===== 6. Update data jika aman =====
	result := db.DB.Model(&models.TbMasterDataDiri{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status_karyawan": "Temporary",
			"employee_id":     newEmployeeID,
		})

	if result.Error != nil {
		return result.Error
	}

	// ===== 7. Kirim data ke API eksternal =====
	payload := map[string]interface{}{
		"full_name":             getStringOrNull(masterData.Nama),
		"phone_number":          getStringOrNull(masterData.NoTelp),
		"home_address":          getStringOrNull(masterData.HomeAddres),
		"domisili":              getStringOrNull(masterData.Domisili),
		"photo":                 getStringOrNull(masterData.Foto),
		"employee_id":           newEmployeeID,
		"no_berkas":             getStringOrNull(masterData.NoBerkas),
		"no_ktp":                getStringOrNull(form2Data.NoKtp),
		"no_npwp":               getStringOrNull(form2Data.NoNpwp),
		"religion":              getStringOrNull(masterData.Agama),
		"birthdate":             masterData.TglLahir,
		"ktp_address":           getStringOrNull(masterData.HomeAddres),
		"no_rekening":           getStringOrNull(form2Data.Norek),
		"nama_rekening":         getStringOrNull(form2Data.Bank),
		"id_status_data_diri":   fmt.Sprintf("%d", masterData.ID),
		"status_pelamar":        getStringOrNull(masterData.StatusKaryawan),
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("gagal marshal payload: %v", err)
	}

	
	// Kirim data ke API eksternal 1
	req, err := http.NewRequest("POST", "https://backend.sigapdriver.com/api/create_drivers_recruitment", bytes.NewBuffer(payloadBytes))
if err != nil {
    log.Printf("Gagal membuat request ke API 1: %v", err)
    return fmt.Errorf("gagal membuat request untuk API external 1: %v", err)
}
req.Header.Set("Content-Type", "application/json")

// Kirim data ke API eksternal 2
req2, err := http.NewRequest("POST", "https://api-invoice-mysql.sigapdriver.com/api/v1/create-drivers-recruitment", bytes.NewBuffer(payloadBytes))
if err != nil {
    log.Printf("Gagal membuat request ke API 2: %v", err)
    return fmt.Errorf("gagal membuat request untuk API external 2: %v", err)
}
req2.Header.Set("Content-Type", "application/json")

	
	// Kirim request ke API eksternal 1 secara asynchronous
	go func() {
    client := &http.Client{Timeout: 30 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        log.Printf("Gagal mengirim request ke API 1: %v", err)
        return
    }
    defer resp.Body.Close()
    
    body, _ := io.ReadAll(resp.Body)
    log.Printf("Response dari API 1 - Status: %d, Body: %s", resp.StatusCode, string(body))
}()
	
	
// Kirim request ke API eksternal 1 secara asynchronous
go func() {
    client := &http.Client{Timeout: 30 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        log.Printf("Gagal mengirim request ke API 1: %v", err)
        return
    }
    defer resp.Body.Close()
    
    body, _ := io.ReadAll(resp.Body)
    log.Printf("Response dari API 1 - Status: %d, Body: %s", resp.StatusCode, string(body))
}()

// Kirim request ke API eksternal 2 secara synchronous
client := &http.Client{Timeout: 30 * time.Second}
resp, err := client.Do(req2)
if err != nil {
    log.Printf("Gagal mengirim request ke API 2: %v", err)
    return fmt.Errorf("gagal mengirim request ke API 2: %v", err)
}
defer resp.Body.Close()

// Baca response body untuk logging
body, _ := io.ReadAll(resp.Body)
log.Printf("Response dari API 2 - Status: %d, Body: %s", resp.StatusCode, string(body))

// Periksa status code
if resp.StatusCode < 200 || resp.StatusCode >= 300 {
    return fmt.Errorf("API eksternal 2 mengembalikan error (status %d): %s", resp.StatusCode, string(body))
}

	return nil
}
// PostMasterDataExternal sends existing master data by ID to the external API without modifying local records.
func PostMasterDataExternal(id int) error {
    // 1. Ambil data master dan form2 untuk payload API
    var masterData models.TbMasterDataDiri
    if err := db.DB.First(&masterData, id).Error; err != nil {
        return fmt.Errorf("gagal mengambil data master: %v", err)
    }

    var form2Data models.Form2
    db.DB.Where("form_1_id = ?", masterData.Form1ID).First(&form2Data)

    // 2. Siapkan payload (gunakan employee_id yang sudah ada)
    payload := map[string]interface{}{
        "full_name":           getStringOrNull(masterData.Nama),
        "phone_number":        getStringOrNull(masterData.NoTelp),
        "home_address":        getStringOrNull(masterData.HomeAddres),
        "domisili":            getStringOrNull(masterData.Domisili),
        "photo":               getStringOrNull(masterData.Foto),
        "employee_id":         getStringOrNull(masterData.EmployeeID),
        "no_berkas":           getStringOrNull(masterData.NoBerkas),
        "no_ktp":              getStringOrNull(form2Data.NoKtp),
        "no_npwp":             getStringOrNull(form2Data.NoNpwp),
        "religion":            getStringOrNull(masterData.Agama),
        "birthdate":           masterData.TglLahir,
        "ktp_address":         getStringOrNull(masterData.HomeAddres),
        "no_rekening":         getStringOrNull(form2Data.Norek),
        "nama_rekening":       getStringOrNull(form2Data.Bank),
        "id_status_data_diri": fmt.Sprintf("%d", masterData.ID),
        "status_pelamar":      getStringOrNull(masterData.StatusKaryawan),
    }

    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        return fmt.Errorf("gagal marshal payload: %v", err)
    }

    req, err := http.NewRequest("POST", "https://backend.sigapdriver.com/api/create_drivers_recruitment", bytes.NewBuffer(payloadBytes))
    if err != nil {
        return fmt.Errorf("gagal membuat request: %v", err)
    }
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{Timeout: 30 * time.Second}
    apiResp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("gagal mengirim data ke API eksternal: %v", err)
    }
    defer apiResp.Body.Close()

    if apiResp.StatusCode != http.StatusOK && apiResp.StatusCode != http.StatusCreated {
        body, _ := io.ReadAll(apiResp.Body)
        return fmt.Errorf("API eksternal mengembalikan error (status %d): %s", apiResp.StatusCode, string(body))
    }

    return nil
}

// Helper function untuk menangani null values
func getStringOrNull(s string) string {
    if s == "" {
        return "null"
    }
    return s
}
