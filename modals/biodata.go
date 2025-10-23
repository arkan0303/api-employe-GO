package modals

type Biodata struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Nama   string `json:"nama"`
	Umur   int    `json:"umur"`
	Alamat string `json:"alamat"`
}