package controllers

import (
	"fmt"
	"time"

	"github.com/buger/jsonparser"
	"github.com/gofiber/fiber/v2"
	"github.com/nikitamirzani323/wl_agen_backend_api/entities"
	"github.com/nikitamirzani323/wl_agen_backend_api/helpers"
	"github.com/nikitamirzani323/wl_agen_backend_api/models"
)

const Fieldcurr_home_redis = "LISTCURR_AGEN"

func Currhome(c *fiber.Ctx) error {
	var obj entities.Model_curr
	var arraobj []entities.Model_curr
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldcurr_home_redis)
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		curr_id, _ := jsonparser.GetString(value, "curr_id")
		curr_multiplier, _ := jsonparser.GetFloat(value, "curr_multiplier")

		obj.Curr_id = curr_id
		obj.Curr_multiplier = float32(curr_multiplier)
		arraobj = append(arraobj, obj)
	})

	if !flag {
		result, err := models.Fetch_currHome()
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldcurr_home_redis, result, 60*time.Minute)
		fmt.Println("CURR AGEN MYSQL")
		return c.JSON(result)
	} else {
		fmt.Println("CURR AGEN CACHE")
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "Success",
			"record":  arraobj,
			"time":    time.Since(render_page).String(),
		})
	}
}
func _deleteredis_curr() {
	val_master := helpers.DeleteRedis(Fieldcurr_home_redis)
	fmt.Printf("Redis Delete AGEN CURR : %d", val_master)

}
