package services

import (
	"errors"

	dbconn "api-rect-go/db"
	jobholder "api-rect-go/modals"
	mysqlModels "api-rect-go/modals/mysql"
)

func UpdateByStatusDataDiri(
	idStatusDataDiri int32,
	newIDCustomer int32,
	newIDUsers int32,
	newServiceUsersID int32,
) error {

	// ===============================
	// 1️⃣ POSTGRESQL - tb_job_holder
	// ===============================
	var jobHolder jobholder.TbJobHolder

	if err := dbconn.DB.
		Where("id_status_data_diri = ?", idStatusDataDiri).
		Order("id DESC").
		Limit(1).
		First(&jobHolder).Error; err != nil {

		return errors.New("job holder terbaru tidak ditemukan")
	}

	if err := dbconn.DB.
		Model(&jobHolder).
		Updates(map[string]interface{}{
			"id_customer": newIDCustomer,
			"id_users":    newIDUsers,
		}).Error; err != nil {

		return err
	}

	// ===============================
	// 2️⃣ MYSQL - service_drivers
	// ===============================
	var serviceDriver mysqlModels.ServiceDriver

	if err := dbconn.DBMySQL.
		Where("id_recruitment = ?", idStatusDataDiri).
		Order("id DESC").
		Limit(1).
		First(&serviceDriver).Error; err != nil {

		return errors.New("service driver tidak ditemukan")
	}

	// ===============================
	// 3️⃣ MYSQL - service_details
	// ===============================
	var serviceDetail mysqlModels.ServiceDetail

	if err := dbconn.DBMySQL.
		Where("service_drivers_id = ?", serviceDriver.ID).
		Order("id DESC").
		Limit(1).
		First(&serviceDetail).Error; err != nil {

		return errors.New("service detail tidak ditemukan")
	}

	if err := dbconn.DBMySQL.
		Model(&serviceDetail).
		Update("service_users_id", newServiceUsersID).Error; err != nil {

		return err
	}

	return nil
}
