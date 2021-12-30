/*
Хелпер для теста приложения в целом
*/
package testhelper

import (
	"testing"

	"github.com/alrusov/config"
	"github.com/alrusov/log"
	"github.com/alrusov/misc"

	_ "github.com/alrusov/auth-basic"
	_ "github.com/alrusov/auth-jwt"
	_ "github.com/alrusov/auth-keycloak"
	_ "github.com/alrusov/auth-krb5"
	_ "github.com/alrusov/auth-url"
)

//----------------------------------------------------------------------------------------------------------------------------//

// Запуск теста
func Start(t *testing.T, home string, fileName string, env misc.StringMap, cfg config.App) (err error) {
	panicID := panic.ID()
	defer panic.SaveStackToLogEx(panicID)

	misc.TEST = true // В коде можно проверять эту переменную, чтобы понимать это тест или нормальная работа

	// Настраиваем логирование

	misc.Logger = log.StdLogger

	log.SetTestWriter(t)
	log.SetFile("-", "", false, 0, 0)
	log.SetLogLevels("TRACE4", misc.StringMap{}, log.FuncNameModeNone)

	// Загружаем конфиг приложения
	err = LoadConfig(home, fileName, env, cfg)
	if err != nil {
		return
	}

	// И проверяем загруженный конфиг
	err = cfg.Check()
	if err != nil {
		t.Fatal(err)
	}

	return
}

//----------------------------------------------------------------------------------------------------------------------------//

// Заканчиваем работу
func Stop(code int) {
	misc.StopApp(code)
}

//----------------------------------------------------------------------------------------------------------------------------//
