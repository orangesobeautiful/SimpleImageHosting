package settingoperation

import (
	"SimpleImageHosting/common"
	"errors"
	"fmt"
	"os"

	"gopkg.in/ini.v1"
)

type SettingProperties struct {
	Server struct {
		Host        string
		Port        int
		FileSaveDir string
		LogFilePath string
		Stdout      bool
		DebugMode   bool
	}
	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
		DSN      string
	}
}

func ReadConfFile() (SettingProperties, error) {
	var settingFilePath string = "setting.conf"
	var cfg *ini.File
	var res SettingProperties
	var hasError bool = false

	// 檢查 setting.conf 是否存在
	if info, err := os.Stat(settingFilePath); err == nil {
		// 如果 setting.conf 存在
		if info.IsDir() {
			fmt.Println(settingFilePath, " 此文件夾與設定檔同名")
		} else {
			cfg, err := ini.Load(settingFilePath)
			if err != nil {
				fmt.Printf("Fail to read file: %v", err)
			}
			var serverSection *ini.Section = cfg.Section("Server")
			res.Server.Host = serverSection.Key("Host").String()
			if res.Server.Host == "" {
				fmt.Println("Note: Server host 為空白效果等同於 0.0.0.0")
			}
			res.Server.Port, err = serverSection.Key("Port").Int()
			if err != nil {
				hasError = true
				fmt.Println("Error: Server port 不是數字")
			} else if 1 > res.Server.Port || res.Server.Port > 65535 {
				hasError = true
				fmt.Println("Error: Server port 範圍錯誤")
			}
			res.Server.FileSaveDir = serverSection.Key("FileSaveDir").String()
			//FileSaveDir 為空字串時，代表當前目錄
			if res.Server.FileSaveDir != "" {
				//檢查目錄位置是否正確
				pathType, err := common.CheckPath(res.Server.FileSaveDir)

				switch pathType {
				case -1:
					hasError = true
					fmt.Println(err.Error())
					break
				case 0:
					hasError = true
					fmt.Println("Error: Server ImgSaveDir \"" + res.Server.FileSaveDir + "\" 不存在")
					break
				case 1:
					hasError = true
					fmt.Println("Error: Server ImgSaveDir \"" + res.Server.FileSaveDir + "\" 是個檔案")
					break
				case 2:
					// 正確 do nothing
					break
				}
			}

			res.Server.LogFilePath = serverSection.Key("Log").String()
			res.Server.Stdout, err = serverSection.Key("Stdout").Bool()
			if err != nil {
				fmt.Println("Warning: Server Stdout 不是正確布林值，使用預設值 false")
				res.Server.Stdout = false
			}
			res.Server.DebugMode, err = serverSection.Key("Debug").Bool()
			if err != nil {
				fmt.Println("Warning: Server debug 不是正確布林值，使用預設值 false")
				res.Server.DebugMode = false
			}

			var dbSection *ini.Section = cfg.Section("Database")
			res.Database.Host = dbSection.Key("Host").String()
			res.Database.Port, err = dbSection.Key("Port").Int()
			if err != nil {
				hasError = true
				fmt.Println("Error: Database port 不是數字")
			} else if 1 > res.Database.Port || res.Database.Port > 65535 {
				hasError = true
				fmt.Println("Error: Database port 範圍錯誤")
			}
			res.Database.User = dbSection.Key("User").String()
			if res.Database.User == "" {
				hasError = true
				fmt.Println("Error: Database user 未設定")
			}
			res.Database.Password = dbSection.Key("Password").String()
			res.Database.DBName = dbSection.Key("DBName").String()
			if res.Database.DBName == "" {
				hasError = true
				fmt.Println("Error: Database dbname 未設定")
			}

			res.Database.DSN = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				res.Database.User, res.Database.Password, res.Database.Host, res.Database.Port, res.Database.DBName)

		}

	} else if os.IsNotExist(err) {
		// 如果 setting.conf 不存在
		cfg = ini.Empty()
		var serverSection *ini.Section = cfg.Section("Server")
		serverSection.Key("Host").SetValue("127.0.0.1")
		serverSection.Key("Port").SetValue("5000")
		serverSection.Key("FileSaveDir").SetValue("")
		serverSection.Key("Log").SetValue("")
		serverSection.Key("Stdout").SetValue("false")
		serverSection.Key("Debug").SetValue("false")

		var dbSection *ini.Section = cfg.Section("Database")
		dbSection.Key("Host").SetValue("127.0.0.1")
		dbSection.Key("Port").SetValue("3306")
		dbSection.Key("User")
		dbSection.Key("Password")
		dbSection.Key("DBName")
		cfg.SaveTo(settingFilePath)
		fmt.Println("找不到檔案 setting.conf， 由程式自動生成，需要進行編輯")
		hasError = true
	} else {
		// 其他錯誤
		fmt.Println("check ", settingFilePath, " status err:", err)
	}

	if hasError {
		return res, errors.New("設定檔有錯誤")
	} else {
		return res, nil
	}

}
