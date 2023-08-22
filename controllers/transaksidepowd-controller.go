package controllers

import (
	"fmt"
	"time"

	"github.com/buger/jsonparser"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nikitamirzani323/wl_agen_backend_api/entities"
	"github.com/nikitamirzani323/wl_agen_backend_api/helpers"
	"github.com/nikitamirzani323/wl_agen_backend_api/models"
)

const Fieldtransdpwd_home_redis = "LISTTRANSDPWD_AGEN"

func Transdpwdhome(c *fiber.Ctx) error {
	user := c.Locals("jwt").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	temp_decp := helpers.Decryption(name)
	_, client_idmasteragen, _, _, _ := helpers.Parsing_Decry(temp_decp, "==")

	var obj entities.Model_transdpwd
	var arraobj []entities.Model_transdpwd
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldtransdpwd_home_redis + "_" + client_idmasteragen)
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		transdpwd_id, _ := jsonparser.GetString(value, "transdpwd_id")
		transdpwd_date, _ := jsonparser.GetString(value, "transdpwd_date")
		transdpwd_idcurr, _ := jsonparser.GetString(value, "transdpwd_idcurr")
		transdpwd_tipeuserdoc, _ := jsonparser.GetString(value, "transdpwd_tipeuserdoc")
		transdpwd_tipedoc, _ := jsonparser.GetString(value, "transdpwd_tipedoc")
		transdpwd_tipeakun, _ := jsonparser.GetString(value, "transdpwd_tipeakun")
		transdpwd_idmember, _ := jsonparser.GetString(value, "transdpwd_idmember")
		transdpwd_nmmember, _ := jsonparser.GetString(value, "transdpwd_nmmember")
		transdpwd_bank_in, _ := jsonparser.GetInt(value, "transdpwd_bank_in")
		transdpwd_bank_out, _ := jsonparser.GetInt(value, "transdpwd_bank_out")
		transdpwd_bank_in_info, _ := jsonparser.GetString(value, "transdpwd_bank_in_info")
		transdpwd_bank_out_info, _ := jsonparser.GetString(value, "transdpwd_bank_out_info")
		transdpwd_amount, _ := jsonparser.GetFloat(value, "transdpwd_amount")
		transdpwd_before, _ := jsonparser.GetFloat(value, "transdpwd_before")
		transdpwd_after, _ := jsonparser.GetFloat(value, "transdpwd_after")
		transdpwd_ipaddress, _ := jsonparser.GetString(value, "transdpwd_ipaddress")
		transdpwd_timezone, _ := jsonparser.GetString(value, "transdpwd_timezone")
		transdpwd_note, _ := jsonparser.GetString(value, "transdpwd_note")
		transdpwd_status, _ := jsonparser.GetString(value, "transdpwd_status")
		transdpwd_status_css, _ := jsonparser.GetString(value, "transdpwd_status_css")
		transdpwd_create, _ := jsonparser.GetString(value, "transdpwd_create")
		transdpwd_update, _ := jsonparser.GetString(value, "transdpwd_update")

		obj.Transdpwd_id = transdpwd_id
		obj.Transdpwd_date = transdpwd_date
		obj.Transdpwd_idcurr = transdpwd_idcurr
		obj.Transdpwd_tipedoc = transdpwd_tipedoc
		obj.Transdpwd_tipeuserdoc = transdpwd_tipeuserdoc
		obj.Transdpwd_tipeakun = transdpwd_tipeakun
		obj.Transdpwd_idmember = transdpwd_idmember
		obj.Transdpwd_nmmember = transdpwd_nmmember
		obj.Transdpwd_bank_in = int(transdpwd_bank_in)
		obj.Transdpwd_bank_out = int(transdpwd_bank_out)
		obj.Transdpwd_bank_in_info = transdpwd_bank_in_info
		obj.Transdpwd_bank_out_info = transdpwd_bank_out_info
		obj.Transdpwd_amount = float64(transdpwd_amount)
		obj.Transdpwd_before = float64(transdpwd_before)
		obj.Transdpwd_after = float64(transdpwd_after)
		obj.Transdpwd_ipaddress = transdpwd_ipaddress
		obj.Transdpwd_timezone = transdpwd_timezone
		obj.Transdpwd_note = transdpwd_note
		obj.Transdpwd_status = transdpwd_status
		obj.Transdpwd_status_css = transdpwd_status_css
		obj.Transdpwd_create = transdpwd_create
		obj.Transdpwd_update = transdpwd_update
		arraobj = append(arraobj, obj)
	})

	if !flag {
		result, err := models.Fetch_transdpwdHome(client_idmasteragen)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldtransdpwd_home_redis+"_"+client_idmasteragen, result, 60*time.Minute)
		fmt.Println("TRANSDPWD AGEN MYSQL")
		return c.JSON(result)
	} else {
		fmt.Println("TRANSDPWD AGEN CACHE")
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "Success",
			"record":  arraobj,
			"time":    time.Since(render_page).String(),
		})
	}
}
func TransdpwdSave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_transdpwdsave)
	validate := validator.New()
	if err := c.BodyParser(client); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	err := validate.Struct(client)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element helpers.ErrorResponse
			element.Field = err.StructField()
			element.Tag = err.Tag()
			errors = append(errors, &element)
		}
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "validation",
			"record":  errors,
		})
	}
	user := c.Locals("jwt").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	temp_decp := helpers.Decryption(name)
	client_idmaster, client_idmasteragen, client_admin, _, _ := helpers.Parsing_Decry(temp_decp, "==")

	// admin, idrecord, idmasteragen, idmaster, tipedoc, idmember, note_dpwd, status, sData string, bank_in, bank_out int, amount float32
	result, err := models.Save_transdpwd(
		client_admin,
		client.Transdpwd_id, client_idmasteragen, client_idmaster, client.Transdpwd_tipedoc, client.Transdpwd_idmember,
		client.Transdpwd_note, client.Transdpwd_status,
		client.Sdata, client.Transdpwd_bank_in, client.Transdpwd_bank_out, client.Transdpwd_amount)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	_deleteredis_transdpwd(client_idmasteragen)
	return c.JSON(result)
}
func _deleteredis_transdpwd(idmasteragen string) {
	val_master := helpers.DeleteRedis(Fieldtransdpwd_home_redis + "_" + idmasteragen)
	fmt.Printf("Redis Delete AGEN TRANSAKSI DEPO WD : %d", val_master)

}
