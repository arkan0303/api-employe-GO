package services

import (
	"api-rect-go/db"
	models "api-rect-go/modals"
)

	type Form5AndForm2AndForm6AndForm4AndForm1 struct {
		ID                        int32  `json:"id"`
		Status                    string `json:"status"`
		Channel                   string `json:"channel"`
		Form1ID                   int32  `json:"form_1_id"`
		NamaLengkap               string `json:"nama_lengkap"`
		NoTelphone                string `json:"no_telphone"`
		NamaKotaDomisili          string `json:"nama_kota_domisili"`
		TahunLahir                string `json:"tahun_lahir"`
		PengalamanJepang          bool   `json:"pengalaman_jepang"`
		PendidikanTerakhir        string `json:"pendidikan_terakhir"`
		Agama                     string `json:"agama"`
		SimANomor                 string `json:"sim_a_nomor"`
		SimAMasaBerlaku           string `json:"sim_a_masa_berlaku"`
		SimB1Nomor                string `json:"sim_b1_nomor"`
		SimB1MasaBerlaku          string `json:"sim_b1_masa_berlaku"`
		SimB2Nomor                string `json:"sim_b2_nomor"`
		SimB2MasaBerlaku          string `json:"sim_b2_masa_berlaku"`
		PengalamanTahunMulai      int32  `json:"pengalaman_tahun_mulai"`
		PengalamanTahunSelesai    int32  `json:"pengalaman_tahun_selesai"`
		PengalamanNamaPerusahaan  string `json:"pengalaman_nama_perusahaan"`
		PengalamanUserEkspat      string `json:"pengalaman_user_ekspat"`
		PengalamanNegaraAsal      string `json:"pengalaman_negara_asal"`
		PengalamanTahunMulai2     int32  `json:"pengalaman_tahun_mulai2"`
		PengalamanTahunSelesai2   int32  `json:"pengalaman_tahun_selesai2"`
		PengalamanNamaPerusahaan2 string `json:"pengalaman_nama_perusahaan2"`
		PengalamanUserEkspat2     string `json:"pengalaman_user_ekspat2"`
		PengalamanNegaraAsal2     string `json:"pengalaman_negara_asal2"`
		PengalamanTahunMulai3     int32  `json:"pengalaman_tahun_mulai3"`
		PengalamanTahunSelesai3   int32  `json:"pengalaman_tahun_selesai3"`
		PengalamanNamaPerusahaan3 string `json:"pengalaman_nama_perusahaan3"`
		PengalamanUserEkspat3     string `json:"pengalaman_user_ekspat3"`
		PengalamanNegaraAsal3     string `json:"pengalaman_negara_asal3"`
		SudahTerdaftar            bool   `json:"sudah_terdaftar"`
		Bank                      string `json:"bank"`
		Norek                     string `json:"norek"`
		StatusPernikahan          string `json:"status_pernikahan"`
		NamaKontakDarurat         string `json:"nama_kontak_darurat"`
		TelponeKontakDarurat      string `json:"telpone_kontak_darurat"`
		HubunganKontakDarurat     string `json:"hubungan_kontak_darurat"`
		AlamatTempatTinggal       string `json:"alamat_tempat_tinggal"`
		FotoCv                    string `json:"foto_cv"`
		FotoDiri                  string `json:"foto_diri"`
		NomorBerkas               string `json:"nomor_berkas"`
		Kerapihan                 string `json:"kerapihan"`
		KemampuanBahasaInggris    string `json:"kemampuan_bahasa_inggris"`
		Pertanyaan6               string `json:"pertanyaan_6"`
	}

	func GetForm5() ([]Form5AndForm2AndForm6AndForm4AndForm1, error) {
		// 1. Ambil semua Form5
		var form5List []models.Form5
		if err := db.DB.Order("id DESC").Find(&form5List).Error; err != nil {
			return nil, err
		}
		if len(form5List) == 0 {
			return []Form5AndForm2AndForm6AndForm4AndForm1{}, nil
		}

		// 2. Kumpulkan Form1ID yang unik
		form1IDSet := make(map[int32]bool)
		for _, f5 := range form5List {
			form1IDSet[f5.Form1ID] = true
		}

		form1IDs := make([]int32, 0, len(form1IDSet))
		for id := range form1IDSet {
			form1IDs = append(form1IDs, id)
		}

		// 3. Ambil semua Form1 sekaligus (1 query)
		var form1List []models.Form1
		if err := db.DB.Where("id IN ?", form1IDs).Find(&form1List).Error; err != nil {
			return nil, err
		}

		form1Map := make(map[int32]models.Form1, len(form1List))
		for _, f1 := range form1List {
			form1Map[f1.ID] = f1
		}

		// 4. Ambil semua Form2 sekaligus (1 query)
		var form2List []models.Form2
		if err := db.DB.Where("form_1_id IN ?", form1IDs).Find(&form2List).Error; err != nil {
			return nil, err
		}

		form2Map := make(map[int32]models.Form2, len(form2List))
		for _, f2 := range form2List {
			form2Map[f2.Form1ID] = f2
		}

		// 5. ✅ OPTIMASI: Ambil semua nama dari Form1 untuk dicek sekaligus
		namaList := make([]string, 0, len(form1List))
		for _, f1 := range form1List {
			if f1.NamaLengkap != "" {
				namaList = append(namaList, f1.NamaLengkap)
			}
		}

		// 6. ✅ Cek semua nama sekaligus dengan 1 query saja!
		var existingNames []struct {
			Nama string
		}
		if len(namaList) > 0 {
			if err := db.DB.Table("tb_master_data_diri").
				Select("nama").
				Where("nama IN ?", namaList).
				Find(&existingNames).Error; err != nil {
				return nil, err
			}
		}

		// 7. ✅ Buat map untuk lookup cepat O(1)
		namaExistsMap := make(map[string]bool, len(existingNames))
		for _, existing := range existingNames {
			namaExistsMap[existing.Nama] = true
		}

		// 7. Ambil semua Form4 sekaligus (1 query)
		var form4List []models.Form4
		if err := db.DB.Where("form_1_id IN ?", form1IDs).Find(&form4List).Error; err != nil {
			return nil, err
		}

		form4Map := make(map[int32]models.Form4, len(form4List))
		for _, f4 := range form4List {
			form4Map[f4.Form1ID] = f4
		}

		// 8. Ambil semua Form6 sekaligus (1 query)
		var form6List []models.Form6
		if err := db.DB.Where("form_1_id IN ?", form1IDs).Find(&form6List).Error; err != nil {
			return nil, err
		}

		form6Map := make(map[int32]models.Form6, len(form6List))
		for _, f6 := range form6List {
			form6Map[f6.Form1ID] = f6
		}

		// 8. Gabungkan semua data
		results := make([]Form5AndForm2AndForm6AndForm4AndForm1, 0, len(form5List))
		for _, f5 := range form5List {
			form1, exists := form1Map[f5.Form1ID]
			form2, exists2 := form2Map[f5.Form1ID]
			form4, exists4 := form4Map[f5.Form1ID]
			form6, exists6 := form6Map[f5.Form1ID]
			if !exists || !exists2 || !exists4 || !exists6 {
				continue
			}

			// ✅ Cek dari map (super cepat O(1))
			sudahAda := namaExistsMap[form1.NamaLengkap]

			results = append(results, Form5AndForm2AndForm6AndForm4AndForm1{
				ID:                        f5.ID,
				Status:                    f5.Status,
				Channel:                   f5.Channel,
				Form1ID:                   f5.Form1ID,
				NamaLengkap:               form1.NamaLengkap,
				NoTelphone:                form1.NoTelphone,
				NamaKotaDomisili:          form1.NamaKotaDomisili,
				TahunLahir:                form1.TahunLahir,
				PengalamanJepang:          form1.PengalamanJepang,
				PendidikanTerakhir:        form2.PendidikanTerakhir,
				Agama:                     form2.Agama,
				SimANomor:                 form2.SimANomor,
				SimAMasaBerlaku:           form2.SimAMasaBerlaku.Format("2006-01-02"),
				SimB1Nomor:                form2.SimB1Nomor,
				SimB1MasaBerlaku:          form2.SimB1MasaBerlaku.Format("2006-01-02"),
				SimB2Nomor:                form2.SimB2Nomor,
				SimB2MasaBerlaku:          form2.SimB2MasaBerlaku.Format("2006-01-02"),
				PengalamanTahunMulai:      form2.PengalamanTahunMulai,
				PengalamanTahunSelesai:    form2.PengalamanTahunSelesai,
				PengalamanNamaPerusahaan:  form2.PengalamanNamaPerusahaan,
				PengalamanUserEkspat:      form2.PengalamanUserEkspat,
				PengalamanNegaraAsal:      form2.PengalamanNegaraAsal,
				PengalamanTahunMulai2:     form2.PengalamanTahunMulai2,
				PengalamanTahunSelesai2:   form2.PengalamanTahunSelesai2,
				PengalamanNamaPerusahaan2: form2.PengalamanNamaPerusahaan2,
				PengalamanUserEkspat2:     form2.PengalamanUserEkspat2,
				PengalamanNegaraAsal2:     form2.PengalamanNegaraAsal2,
				PengalamanTahunMulai3:     form2.PengalamanTahunMulai3,
				PengalamanTahunSelesai3:   form2.PengalamanTahunSelesai3,
				PengalamanNamaPerusahaan3: form2.PengalamanNamaPerusahaan3,
				PengalamanUserEkspat3:     form2.PengalamanUserEkspat3,
				PengalamanNegaraAsal3:     form2.PengalamanNegaraAsal3,
				Bank:                      form2.Bank,
				Norek:                     form2.Norek,
				StatusPernikahan:          form2.StatusPernikahan,
				SudahTerdaftar:            sudahAda,
				NamaKontakDarurat:         form2.NamaKontakDarurat,
				TelponeKontakDarurat:      form2.TelponeKontakDarurat,
				HubunganKontakDarurat:     form2.HubunganKontakDarurat,
				AlamatTempatTinggal:       form2.AlamatTempatTinggal,
				FotoCv:                    form6.FotoCv,
				FotoDiri:                  form4.FotoDiri,
				NomorBerkas:               f5.NomorBerkas,
				Kerapihan:                 form6.Kerapihan,
				KemampuanBahasaInggris:    form6.KemampuanBahasaInggris,
				Pertanyaan6:               form6.Pertanyaan6,
			})
		}

		return results, nil
	}