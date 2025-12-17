package dto

type Form5Response struct {
	ID          int32  `json:"id"`
	NamaLengkap string `json:"nama_lengkap"`
	NoTelphone  string `json:"no_telphone"`
	NomorBerkas string `json:"nomor_berkas"`
	Status      string `json:"status"`
	Tanggal     string `json:"tanggal"`
}

type Form2Response struct {
	Gender             string `json:"gender"`
	Agama              string `json:"agama"`
	NoKtp              string `json:"no_ktp"`
	PendidikanTerakhir string `json:"pendidikan_terakhir"`
}

type PelamarCombinedResponse struct {
	Form5 Form5Response `json:"form5"`
	Form2 Form2Response `json:"form2"`
}
