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

const Fieldagenbank_home_redis = "LISTAGENBANK_AGEN"
const Fieldagenbankshare_home_redis = "LISTAGENBANKSHARE_AGEN"

func Agenbankhome(c *fiber.Ctx) error {
	user := c.Locals("jwt").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	temp_decp := helpers.Decryption(name)
	_, client_idmasteragen, _, _, _ := helpers.Parsing_Decry(temp_decp, "==")

	var obj entities.Model_agenbank
	var arraobj []entities.Model_agenbank
	var obj_listbanktype entities.Model_bankTypeshare
	var arraobj_listbanktype []entities.Model_bankTypeshare
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldagenbank_home_redis + "_" + client_idmasteragen)
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	listbank_RD, _, _, _ := jsonparser.Get(jsonredis, "listbank")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		agenbank_id, _ := jsonparser.GetInt(value, "agenbank_id")
		agenbank_tipe, _ := jsonparser.GetString(value, "agenbank_tipe")
		agenbank_idbanktype, _ := jsonparser.GetString(value, "agenbank_idbanktype")
		agenbank_norek, _ := jsonparser.GetString(value, "agenbank_norek")
		agenbank_nmrek, _ := jsonparser.GetString(value, "agenbank_nmrek")
		agenbank_status, _ := jsonparser.GetString(value, "agenbank_status")
		agenbank_status_css, _ := jsonparser.GetString(value, "agenbank_status_css")
		agenbank_create, _ := jsonparser.GetString(value, "agenbank_create")
		agenbank_update, _ := jsonparser.GetString(value, "agenbank_update")

		obj.Agenbank_id = int(agenbank_id)
		obj.Agenbank_tipe = agenbank_tipe
		obj.Agenbank_idbanktype = agenbank_idbanktype
		obj.Agenbank_norek = agenbank_norek
		obj.Agenbank_nmrek = agenbank_nmrek
		obj.Agenbank_status = agenbank_status
		obj.Agenbank_status_css = agenbank_status_css
		obj.Agenbank_create = agenbank_create
		obj.Agenbank_update = agenbank_update
		arraobj = append(arraobj, obj)
	})
	jsonparser.ArrayEach(listbank_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		catebank_name, _ := jsonparser.GetString(value, "catebank_name")
		banktype_id, _ := jsonparser.GetString(value, "banktype_id")

		obj_listbanktype.Catebank_name = catebank_name
		obj_listbanktype.Banktype_id = banktype_id
		arraobj_listbanktype = append(arraobj_listbanktype, obj_listbanktype)
	})
	if !flag {
		result, err := models.Fetch_agenbankHome(client_idmasteragen)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldagenbank_home_redis+"_"+client_idmasteragen, result, 60*time.Minute)
		fmt.Println("CATEBANK MYSQL")
		return c.JSON(result)
	} else {
		fmt.Println("CATEBANK CACHE")
		return c.JSON(fiber.Map{
			"status":   fiber.StatusOK,
			"message":  "Success",
			"record":   arraobj,
			"listbank": arraobj_listbanktype,
			"time":     time.Since(render_page).String(),
		})
	}
}
func Agenbanklist(c *fiber.Ctx) error {
	user := c.Locals("jwt").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	temp_decp := helpers.Decryption(name)
	_, client_idmasteragen, _, _, _ := helpers.Parsing_Decry(temp_decp, "==")

	var obj entities.Model_agenbankshare
	var arraobj []entities.Model_agenbankshare
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldagenbankshare_home_redis + "_" + client_idmasteragen)
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		agenbank_id, _ := jsonparser.GetInt(value, "agenbank_id")
		agenbank_info, _ := jsonparser.GetString(value, "agenbank_info")

		obj.Agenbank_id = int(agenbank_id)
		obj.Agenbank_info = agenbank_info
		arraobj = append(arraobj, obj)
	})

	if !flag {
		result, err := models.Fetch_agenbankList(client_idmasteragen)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldagenbankshare_home_redis+"_"+client_idmasteragen, result, 60*time.Minute)
		fmt.Println("LIST MEMBER BANK MYSQL")
		return c.JSON(result)
	} else {
		fmt.Println("LIST MEMBER BANK CACHE")
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "Success",
			"record":  arraobj,
			"time":    time.Since(render_page).String(),
		})
	}
}
func AgenbankSave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_agenbanksave)
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
	_, client_idmasteragen, client_admin, _, _ := helpers.Parsing_Decry(temp_decp, "==")

	// admin, idmasteragen, tipe, idbank, norek, nmrek, status, sData string, idrecord int
	result, err := models.Save_agenbank(
		client_admin, client_idmasteragen,
		client.Agenbank_tipe, client.Agenbank_idbanktype, client.Agenbank_norek, client.Agenbank_nmrek, client.Agenbank_status,
		client.Sdata, client.Agenbank_id)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	_deleteredis_agenbank(client_idmasteragen)
	return c.JSON(result)
}
func _deleteredis_agenbank(idmasteragen string) {
	val_master := helpers.DeleteRedis(Fieldagenbank_home_redis + "_" + idmasteragen)
	fmt.Printf("Redis Delete AGEN BANK : %d", val_master)
}
