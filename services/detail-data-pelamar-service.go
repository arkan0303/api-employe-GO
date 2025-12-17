package services

import (
	"api-rect-go/db"
	"api-rect-go/dto"
	models "api-rect-go/modals"
)

func GetDataPelamar() ([]dto.PelamarCombinedResponse, error) {
    // Get all Form5
    var form5Data []models.Form5
    result := db.DB.Find(&form5Data)
    if result.Error != nil {
        return nil, result.Error
    }

    // Collect all form1_ids
    var form1IDs []int32
    for _, form5 := range form5Data {
        if form5.Form1ID != 0 {
            form1IDs = append(form1IDs, form5.Form1ID)
        }
    }

    // Get Form2
    var form2Data []models.Form2
    if len(form1IDs) > 0 {
        result = db.DB.Where("form_1_id IN ?", form1IDs).Find(&form2Data)
        if result.Error != nil {
            return nil, result.Error
        }
    }

    // Map Form2 by Form1ID
    form2Map := make(map[int32]models.Form2)
    for _, form2 := range form2Data {
        form2Map[form2.Form1ID] = form2
    }

    // FINAL RESULT (with smaller struct)
    var resultData []dto.PelamarCombinedResponse

    for _, form5 := range form5Data {
        f5 := dto.Form5Response{
            ID:          form5.ID,
            NamaLengkap: form5.NamaLengkap,
            NoTelphone:  form5.NoTelphone,
            NomorBerkas: form5.NomorBerkas,
            Status:      form5.Status,
            Tanggal:     form5.Tanggal.Format("2006-01-02"),
        }

        f2 := dto.Form2Response{} // default

        if form2, exists := form2Map[form5.Form1ID]; exists {
            f2 = dto.Form2Response{
                Gender:             form2.Gender,
                Agama:              form2.Agama,
                NoKtp:              form2.NoKtp,
                PendidikanTerakhir: form2.PendidikanTerakhir,
            }
        }

        resultData = append(resultData, dto.PelamarCombinedResponse{
            Form5: f5,
            Form2: f2,
        })
    }

    return resultData, nil
}
