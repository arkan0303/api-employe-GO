package modals

type Book struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Nama  string `json:"nama"`
	Harga int    `json:"harga"`
}