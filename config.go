package testhelper

import (
	"os"

	"github.com/alrusov/config"
	"github.com/alrusov/misc"
)

//----------------------------------------------------------------------------------------------------------------------------//

// Загрузка конфига приложения
func LoadConfig(home string, fileName string, env misc.StringMap, cfg interface{}) (err error) {
	if home != "" {
		// Идем в заданную (домашнюю для приложения) директорию, если она задана
		os.Chdir(home)
	}

	// Загружаем переменные из стандартного environment (обычно из файла .env)
	err = misc.LoadEnv("")
	if err != nil {
		return
	}

	// Выставляем их в environment приложения
	for k, v := range env {
		err = os.Setenv(k, v)
		if err != nil {
			return
		}
	}

	// Загружаем конфиг
	err = config.LoadFile("$"+fileName, cfg)
	if err != nil {
		return
	}

	return
}

//----------------------------------------------------------------------------------------------------------------------------//
