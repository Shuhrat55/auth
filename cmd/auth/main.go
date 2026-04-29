// @title           Auth Service API
// @version         1.0
// @description     API сервис аутентификации и управления пользователями
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8081
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Используйте формат: Bearer {token}

package main

import "github.com/Shuhrat55/auth/internal/app"

func main() {
	 app.LoggerRun()
	go app.Run()
	go app.StartGRPCServer()
	select {}

}
