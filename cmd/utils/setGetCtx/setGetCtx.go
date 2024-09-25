package setterGetterCtx

import "github.com/gofiber/fiber/v2"

func SetDataInFiberCtx(c *fiber.Ctx, key string, value interface{}) {
    c.Locals(key, value) // Store the data in Fiber's local storage
}


func GetDataFromFiberCtx(c *fiber.Ctx, key string) interface{} {
    return c.Locals(key) // Retrieve the stored data
}
