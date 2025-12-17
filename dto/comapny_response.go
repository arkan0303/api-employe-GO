package dto

type MasterCompany struct {
	ID          int32  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	CompanyName string `gorm:"column:company_name;not null" json:"company_name"`
}