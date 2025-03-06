package dir

import "os"

// Папка config
func configFiles() string {
	dir, err := os.Getwd() // Рабочая директория(.\cmd)
	if err != nil {
		panic(err)
	}
	dir = dir[:len(dir)-3]
	dir += `\config\`
	return dir
}

// config/config.json
func JsonFile() string {
	res := configFiles() + `config.json`
	return res
}

// config/.env
func EnvFile() string {
	res := configFiles() + `.env`
	return res
}
