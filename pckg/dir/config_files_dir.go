package dir

import "os"

// Папка config
func config_files() string {
	dir, err := os.Getwd() // Рабочая директория(.\cmd)
	if err != nil {
		panic(err)
	}
	dir += `\config\`
	return dir
}

// config/config.json
func Json_file() string {
	res := config_files() + `config.json`
	return res
}

// config/.env
func Env_file() string {
	res := config_files() + `.env`
	return res
}
