package model

type SystemSetting struct {
	Key   string `gorm:"primaryKey;size:128" json:"key"`
	Value string `gorm:"type:text" json:"value"`
}

func GetSetting(key string) string {
	var setting SystemSetting
	if err := DB.Where("`key` = ?", key).First(&setting).Error; err != nil {
		return ""
	}
	return setting.Value
}

func SetSetting(key, value string) error {
	setting := SystemSetting{Key: key, Value: value}
	return DB.Where("`key` = ?", key).Assign(SystemSetting{Value: value}).FirstOrCreate(&setting).Error
}

func GetAllSettings() (map[string]string, error) {
	var settings []SystemSetting
	err := DB.Find(&settings).Error
	if err != nil {
		return nil, err
	}
	result := make(map[string]string)
	for _, s := range settings {
		result[s.Key] = s.Value
	}
	return result, nil
}

func BatchSetSettings(settings map[string]string) error {
	tx := DB.Begin()
	for key, value := range settings {
		setting := SystemSetting{Key: key, Value: value}
		if err := tx.Where("`key` = ?", key).Assign(SystemSetting{Value: value}).FirstOrCreate(&setting).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}
