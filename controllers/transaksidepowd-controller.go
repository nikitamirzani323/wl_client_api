package controllers

import (
	"fmt"
	"log"
	"time"

	"github.com/buger/jsonparser"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nikitamirzani323/wl_agen_backend_api/entities"
	"github.com/nikitamirzani323/wl_agen_backend_api/helpers"
	"github.com/nikitamirzani323/wl_agen_backend_api/models"
)

const Fieldtransdpwd_home_redis = "LISTTRANSDPWD_CLIENT"

func Transdpwdhome(c *fiber.Ctx) error {
	user := c.Locals("jwt").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	temp_decp := helpers.Decryption(name)
	log.Println(user)
	_, client_idmasteragen, client_idmember := helpers.Parsing_Decry(temp_decp, "==")

	log.Println(client_idmasteragen)
	log.Println(client_idmember)

	var obj entities.Model_transdpwd
	var arraobj []entities.Model_transdpwd
	//BANK MEMBER
	var objbank entities.Model_memberbank
	var arraobjbank []entities.Model_memberbank
	//BANK AGEN
	var objagenbank entities.Model_agenbankshare
	var arraobjagenbank []entities.Model_agenbankshare
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldtransdpwd_home_redis + "_" + client_idmasteragen + "_" + client_idmember)
	jsonredis := []byte(resultredis)
	listbankmember_RD, _, _, _ := jsonparser.Get(jsonredis, "listbankmember")
	listbankagen_RD, _, _, _ := jsonparser.Get(jsonredis, "listbankagen")
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		transdpwd_id, _ := jsonparser.GetString(value, "transdpwd_id")
		transdpwd_date, _ := jsonparser.GetString(value, "transdpwd_date")
		transdpwd_idcurr, _ := jsonparser.GetString(value, "transdpwd_idcurr")
		transdpwd_tipedoc, _ := jsonparser.GetString(value, "transdpwd_tipedoc")
		transdpwd_bank_in, _ := jsonparser.GetInt(value, "transdpwd_bank_in")
		transdpwd_bank_out, _ := jsonparser.GetInt(value, "transdpwd_bank_out")
		transdpwd_bank_in_info, _ := jsonparser.GetString(value, "transdpwd_bank_in_info")
		transdpwd_bank_out_info, _ := jsonparser.GetString(value, "transdpwd_bank_out_info")
		transdpwd_amount, _ := jsonparser.GetFloat(value, "transdpwd_amount")
		transdpwd_note, _ := jsonparser.GetString(value, "transdpwd_note")
		transdpwd_status, _ := jsonparser.GetString(value, "transdpwd_status")
		transdpwd_status_css, _ := jsonparser.GetString(value, "transdpwd_status_css")
		transdpwd_create, _ := jsonparser.GetString(value, "transdpwd_create")
		transdpwd_update, _ := jsonparser.GetString(value, "transdpwd_update")

		obj.Transdpwd_id = transdpwd_id
		obj.Transdpwd_date = transdpwd_date
		obj.Transdpwd_idcurr = transdpwd_idcurr
		obj.Transdpwd_tipedoc = transdpwd_tipedoc
		obj.Transdpwd_bank_in = int(transdpwd_bank_in)
		obj.Transdpwd_bank_out = int(transdpwd_bank_out)
		obj.Transdpwd_bank_in_info = transdpwd_bank_in_info
		obj.Transdpwd_bank_out_info = transdpwd_bank_out_info
		obj.Transdpwd_amount = float64(transdpwd_amount)
		obj.Transdpwd_note = transdpwd_note
		obj.Transdpwd_status = transdpwd_status
		obj.Transdpwd_status_css = transdpwd_status_css
		obj.Transdpwd_create = transdpwd_create
		obj.Transdpwd_update = transdpwd_update
		arraobj = append(arraobj, obj)
	})
	jsonparser.ArrayEach(listbankmember_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		memberbank_id, _ := jsonparser.GetInt(value, "memberbank_id")
		memberbank_idbanktype, _ := jsonparser.GetString(value, "memberbank_idbanktype")
		memberbank_norek, _ := jsonparser.GetString(value, "memberbank_norek")
		memberbank_nmownerbank, _ := jsonparser.GetString(value, "memberbank_nmownerbank")

		objbank.Memberbank_id = int(memberbank_id)
		objbank.Memberbank_idbanktype = memberbank_idbanktype
		objbank.Memberbank_norek = memberbank_norek
		objbank.Memberbank_nmownerbank = memberbank_nmownerbank
		arraobjbank = append(arraobjbank, objbank)
	})
	jsonparser.ArrayEach(listbankagen_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		agenbank_id, _ := jsonparser.GetInt(value, "agenbank_id")
		agenbank_info, _ := jsonparser.GetString(value, "agenbank_info")

		objagenbank.Agenbank_id = int(agenbank_id)
		objagenbank.Agenbank_info = agenbank_info
		arraobjagenbank = append(arraobjagenbank, objagenbank)
	})
	if !flag {
		result, err := models.Fetch_transdpwdHome(client_idmasteragen, client_idmember)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldtransdpwd_home_redis+"_"+client_idmasteragen+"_"+client_idmember, result, 60*time.Minute)
		fmt.Println("TRANSDPWD CLIENT MYSQL")
		return c.JSON(result)
	} else {
		fmt.Println("TRANSDPWD CLIENT CACHE")
		return c.JSON(fiber.Map{
			"status":         fiber.StatusOK,
			"message":        "Success",
			"record":         arraobj,
			"listbankmember": arraobjbank,
			"listbankagen":   arraobjagenbank,
			"time":           time.Since(render_page).String(),
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
	client_idmaster, client_idmasteragen, client_idmember := helpers.Parsing_Decry(temp_decp, "==")

	// idmember, idrecord, idmasteragen, idmaster, tipedoc, note, ipaddress, timezone, sData string, bank_in, bank_out int, amount float32
	result, err := models.Save_transdpwd(
		client_idmember, client_idmasteragen, client_idmaster, client.Transdpwd_tipedoc,
		client.Transdpwd_note, client.Transdpwd_ipaddress, client.Transdpwd_timezone,
		client.Sdata, client.Transdpwd_bank_in, client.Transdpwd_bank_out, client.Transdpwd_amount)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	_deleteredis_transdpwd(client_idmasteragen, client_idmember)
	return c.JSON(result)
}
func _deleteredis_transdpwd(idmasteragen, idmember string) {
	val_client := helpers.DeleteRedis(Fieldtransdpwd_home_redis + "_" + idmasteragen + "_" + idmember)
	fmt.Printf("Redis Delete AGEN TRANSAKSI DEPO WD : %d\n", val_client)

	val_agen := helpers.DeleteRedis("LISTTRANSDPWD_AGEN_" + idmasteragen)
	fmt.Printf("Redis Delete AGEN TRANSAKSI DEPO WD : %d\n", val_agen)

}
