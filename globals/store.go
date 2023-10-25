package globals

import "github.com/gofiber/fiber/v2/middleware/session"

var Store *session.Store

func InitStore() {
	Store = session.New()
}
