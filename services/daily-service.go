package services

import (
	"api-rect-go/db"
	models "api-rect-go/modals"
	"fmt"
	"strconv"
	"sync"
)

type MasterRingkasss struct {
	ID         int32  `json:"id" gorm:"column:id;primaryKey;autoIncrement:true"`
	Nama       string `json:"nama" gorm:"column:nama"`
	EmployeeID string `json:"employee_id" gorm:"column:employee_id"`
	NoTelp     string `json:"no_telp" gorm:"column:no_telp"`
	Form1ID    int32  `json:"form_1_id" gorm:"column:form_1_id"`
	Foto       string `json:"foto" gorm:"column:foto"`
	TglLahir   string `gorm:"column:tgl_lahir" json:"tgl_lahir"`
	StatusKawin string `gorm:"column:status_kawin" json:"status_kawin"`
	Agama       string    `gorm:"column:agama" json:"agama"`
	Bank                  string    `gorm:"column:bank" json:"bank"`
	Norek                 string    `gorm:"column:norek" json:"norek"`
	Domisili              string    `gorm:"column:domisili" json:"domisili"`
	HomeAddres            string    `gorm:"column:home_addres" json:"home_addres"`
}

type FormGabungannn struct {
	Master	                  MasterRingkasss `json:"master"`
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
	PendukungKtp              string        `json:"pendukung_ktp"`
	PendukungNpwp             string        `json:"pendukung_npwp"`
	PengalamanList            []PengalamanItem `json:"pengalaman_list"`
	IDPerusahaan              int32         `json:"id_perusahaan"`
	IDCustomer                int32         `json:"id_customer"`
	TglJobHolder              string        `json:"tgl_job_holder"`
}

// Struct helper untuk Form2 data
type form2Dataaa struct {
	Agama              string
	PendidikanTerakhir string
	SimANomor          string
	SimAMasaBerlaku    string
	SimB1Nomor         string
	SimB1MasaBerlaku   string
	SimB2Nomor         string
	SimB2MasaBerlaku   string
	PengalamanData     [3]pengalamanInfooo
	NoKtp              string
	NoNpwp             string
}

type pengalamanInfooo struct {
	TahunMulai     int32
	TahunSelesai   int32
	NamaPerusahaan string
	UserEkspat     string
	NegaraAsal     string
}

// Pengalaman dari tabel tb_pengalaman (dipisahkan dari Form2)
type PengalamanItemm struct {
	TahunMulai   int32  `json:"tahun_mulai"`
	TahunSelesai int32  `json:"tahun_selesai"`
	Perusahaan   string `json:"perusahaan"`
	Jabatan      string `json:"jabatan"`
}

func GetMasterDataDaily() ([]FormGabungannn, error) {
	var (
		masters        []MasterRingkasss
		form1Map       map[int32]bool
		form2Map       map[int32]form2Dataaa
		form6Map       map[int32]struct {
			Kerapihan              string
			KemampuanBahasaInggris string
			Pertanyaan6            string
		}
		pendukungMap  map[int32]models.TbPendukungTbMasterDataDiri
		pengalamanMap map[int32][]models.TbPengalaman
		statusMap     map[int32]models.TbStatusDiri
		jobHolderMap  map[int32]models.TbJobHolder
		wg     sync.WaitGroup
		mu     sync.Mutex
	)

	errChan := make(chan error, 4)

	// Fetch master data synchronously so we can safely build masterIDs before spawning other goroutines
	if err := db.DB.Model(&models.TbMasterDataDiri{}).
		Select("id", "nama", "employee_id", "no_telp", "form_1_id", "foto", "tgl_lahir", "status_kawin", "agama", "bank", "norek", "domisili", "home_addres").
		Where("status_karyawan = ?", "daily").
		Scan(&masters).Error; err != nil {
		return nil, fmt.Errorf("master data fetch failed: %w", err)
	}


	fmt.Printf("[INFO] Found %d master records\n", len(masters))

	// ===== STEP 2-3: Siapkan ID untuk query terkait =====
	form1IDs := make([]int32, 0, len(masters))
	masterIDs := make([]int32, 0, len(masters))
	for _, master := range masters {
		form1IDs = append(form1IDs, master.Form1ID)
		masterIDs = append(masterIDs, master.ID)
	}

	wg.Add(6)
	errChan = make(chan error, 6)

	// Query Form1
	go func() {
		defer wg.Done()
		var form1Results []struct {
			ID               int32 `gorm:"column:id"`
			PengalamanJepang bool  `gorm:"column:pengalaman_jepang"`
		}
		if err := db.DB.Model(&models.Form1{}).
			Select("id", "pengalaman_jepang").
			Where("id IN ?", form1IDs).
			Scan(&form1Results).Error; err != nil {
			errChan <- err
			return
		}

		m := make(map[int32]bool, len(form1Results))
		for _, f := range form1Results {
			m[f.ID] = f.PengalamanJepang
		}
		mu.Lock()
		form1Map = m
		mu.Unlock()
	}()

	// Query Form2
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
		if err := db.DB.Model(&models.Form2{}).
			Select("form_1_id", "agama", "pendidikan_terakhir", "sim_a_nomor", "sim_a_masa_berlaku",
				"sim_b1_nomor", "sim_b1_masa_berlaku", "sim_b2_nomor", "sim_b2_masa_berlaku",
				"pengalaman_tahun_mulai", "pengalaman_tahun_selesai", "pengalaman_nama_perusahaan",
				"pengalaman_user_ekspat", "pengalaman_negara_asal", "pengalaman_tahun_mulai2",
				"pengalaman_tahun_selesai2", "pengalaman_nama_perusahaan2", "pengalaman_user_ekspat2",
				"pengalaman_negara_asal2", "pengalaman_tahun_mulai3", "pengalaman_tahun_selesai3",
				"pengalaman_nama_perusahaan3", "pengalaman_user_ekspat3", "pengalaman_negara_asal3",
				"no_ktp", "no_npwp").
			Where("form_1_id IN ?", form1IDs).
			Scan(&form2Results).Error; err != nil {
			errChan <- err
			return
		}

		m := make(map[int32]form2Dataaa, len(form2Results))
		for _, f := range form2Results {
			m[f.Form1ID] = form2Dataaa{
				Agama:              f.Agama,
				PendidikanTerakhir: f.PendidikanTerakhir,
				SimANomor:          f.SimANomor,
				SimAMasaBerlaku:    f.SimAMasaBerlaku,
				SimB1Nomor:         f.SimB1Nomor,
				SimB1MasaBerlaku:   f.SimB1MasaBerlaku,
				SimB2Nomor:         f.SimB2Nomor,
				SimB2MasaBerlaku:   f.SimB2MasaBerlaku,
				PengalamanData: [3]pengalamanInfooo{
					{f.PengalamanTahunMulai, f.PengalamanTahunSelesai, f.PengalamanNamaPerusahaan, f.PengalamanUserEkspat, f.PengalamanNegaraAsal},
					{f.PengalamanTahunMulai2, f.PengalamanTahunSelesai2, f.PengalamanNamaPerusahaan2, f.PengalamanUserEkspat2, f.PengalamanNegaraAsal2},
					{f.PengalamanTahunMulai3, f.PengalamanTahunSelesai3, f.PengalamanNamaPerusahaan3, f.PengalamanUserEkspat3, f.PengalamanNegaraAsal3},
				},
				NoKtp:  f.NoKtp,
				NoNpwp: f.NoNpwp,
			}
		}
		mu.Lock()
		form2Map = m
		mu.Unlock()
	}()

	// Query Form6
	go func() {
		defer wg.Done()
		var form6Results []struct {
			Form1ID                int32  `gorm:"column:form_1_id"`
			Kerapihan              string `gorm:"column:kerapihan"`
			KemampuanBahasaInggris string `gorm:"column:kemampuanBahasaInggris"`
			Pertanyaan6            string `gorm:"column:pertanyaan_6"`
		}
		if err := db.DB.Model(&models.Form6{}).
			Select("form_1_id", "kerapihan", "kemampuanBahasaInggris", "pertanyaan_6").
			Where("form_1_id IN ?", form1IDs).
			Scan(&form6Results).Error; err != nil {
			errChan <- err
			return
		}

		m := make(map[int32]struct {
			Kerapihan              string
			KemampuanBahasaInggris string
			Pertanyaan6            string
		}, len(form6Results))
		for _, f := range form6Results {
			m[f.Form1ID] = struct {
				Kerapihan              string
				KemampuanBahasaInggris string
				Pertanyaan6            string
			}{f.Kerapihan, f.KemampuanBahasaInggris, f.Pertanyaan6}
		}
		mu.Lock()
		form6Map = m
		mu.Unlock()
	}()

	// Query TbStatusDiri (id_perusahaan) dan TbJobHolder (id_customer, tgl_job_holder)
	go func() {
		defer wg.Done()
		var statuses []models.TbStatusDiri
		if err := db.DB.Model(&models.TbStatusDiri{}).
			Where("id_master_data_diri IN ?", masterIDs).
			Scan(&statuses).Error; err != nil {
			errChan <- err
			return
		}
		mStatus := make(map[int32]models.TbStatusDiri, len(statuses))
		for _, s := range statuses {
			mStatus[s.IDMasterDataDiri] = s
		}
		var jobHolders []models.TbJobHolder
		if err := db.DB.Model(&models.TbJobHolder{}).
			Where("id_status_data_diri IN ?", masterIDs).
			Scan(&jobHolders).Error; err != nil {
			errChan <- err
			return
		}
		mJob := make(map[int32]models.TbJobHolder)
		for _, j := range jobHolders {
			mJob[j.IDStatusDataDiri] = j
		}
		mu.Lock()
		statusMap = mStatus
		jobHolderMap = mJob
		mu.Unlock()
	}()

    // Query TbPendukungTbMasterDataDiri berdasarkan id_master_data_diri
    go func() {
        defer wg.Done()
        var pendukungs []models.TbPendukungTbMasterDataDiri
        if err := db.DB.Model(&models.TbPendukungTbMasterDataDiri{}).
            Where("id_master_data_diri IN ?", masterIDs).
            Scan(&pendukungs).Error; err != nil {
            errChan <- err
            return
        }
        m := make(map[int32]models.TbPendukungTbMasterDataDiri, len(pendukungs))
        for _, p := range pendukungs {
            m[p.IDMasterDataDiri] = p
        }
        mu.Lock()
        pendukungMap = m
        mu.Unlock()
    }()

    // Query TbPengalaman (bisa multiple per master) berdasarkan id_master_data_diri
    go func() {
        defer wg.Done()
        var pengalamans []models.TbPengalaman
        if err := db.DB.Model(&models.TbPengalaman{}).
            Where("id_master_data_diri IN ?", masterIDs).
            Scan(&pengalamans).Error; err != nil {
            errChan <- err
            return
        }
        m := make(map[int32][]models.TbPengalaman)
        for _, p := range pengalamans {
            m[p.IDMasterDataDiri] = append(m[p.IDMasterDataDiri], p)
        }
        mu.Lock()
        pengalamanMap = m
        mu.Unlock()
    }()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	// ===== STEP 4: Gabungkan Data dengan Validasi yang Ketat =====
	results := make([]FormGabungannn, 0, len(masters))
	
	// Counter untuk tracking
	var sudahAda, belumAda int

	for _, master := range masters {
		form1 := form1Map[master.Form1ID]
		form2 := form2Map[master.Form1ID]
		form6 := form6Map[master.Form1ID]

		// Ambil data pendukung (KTP/NPWP) berdasarkan master.ID
		pend := pendukungMap[master.ID]

		// Ambil pengalaman dari tb_pengalaman (list terpisah)
		expList := pengalamanMap[master.ID]
		pengalamanItems := make([]PengalamanItem, 0, len(expList))
		for i := 0; i < len(expList); i++ {
			var tMulai, tSelesai int32
			if v, err := strconv.Atoi(expList[i].TahunAwal); err == nil {
				tMulai = int32(v)
			}
			if v, err := strconv.Atoi(expList[i].TahunAkhir); err == nil {
				tSelesai = int32(v)
			}
			pengalamanItems = append(pengalamanItems, PengalamanItem{
				TahunMulai:   tMulai,
				TahunSelesai: tSelesai,
				Perusahaan:   expList[i].Perusahaan,
				Jabatan:      expList[i].Jabatan,
			})
		}

		results = append(results, FormGabungannn{
			Master:                  	master,
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
			PendukungKtp:              pend.Ktp,
			PendukungNpwp:             pend.Npwp,
			PengalamanList:            pengalamanItems,
			IDPerusahaan:              statusMap[master.ID].IDPerusahaan,
			IDCustomer:                jobHolderMap[master.ID].IDCustomer,
			TglJobHolder:              jobHolderMap[master.ID].TglJobHolder,
			Kerapihan:                 form6.Kerapihan,
			KemampuanBahasaInggris:    form6.KemampuanBahasaInggris,
			Pertanyaan6:               form6.Pertanyaan6,
		})
	}

	fmt.Printf("[INFO] Status recruitment - Sudah ada: %d, Belum ada: %d\n", sudahAda, belumAda)

	return results, nil
}
