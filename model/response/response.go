package response

import (
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func Result(data interface{}, msg string, c *fiber.Ctx) error {

	return c.Status(fiber.StatusOK).JSON(Response{
		data,
		msg,
	})
}

func FailWithValidate(err error, c *fiber.Ctx) error {
	return c.Status(fiber.StatusUnauthorized).JSON(Response{
		map[string]interface{}{},
		err.Error(),
	})
}

func FailWithError(err error, c *fiber.Ctx) error {
	return c.Status(fiber.StatusExpectationFailed).JSON(Response{
		map[string]interface{}{},
		err.Error(),
	})
}

func Success(c *fiber.Ctx) error {
	return Result(map[string]interface{}{}, "操作成功", c)
}

func SuccessWithData(data interface{}, c *fiber.Ctx) error {
	return Result(data, "操作成功", c)
}
