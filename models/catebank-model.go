package models

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nikitamirzani323/wl_agen_backend_api/configs"
	"github.com/nikitamirzani323/wl_agen_backend_api/db"
	"github.com/nikitamirzani323/wl_agen_backend_api/entities"
	"github.com/nikitamirzani323/wl_agen_backend_api/helpers"
	"github.com/nleeper/goment"
)

const database_catebank_local = configs.DB_tbl_mst_cate_bank
const database_banktype_local = configs.DB_tbl_mst_banktype

func Fetch_catebankHome() (helpers.Response, error) {
	var obj entities.Model_catebank
	var arraobj []entities.Model_catebank
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
			idcatebank , nmcatebank, statuscatebank,  
			createcatebank, to_char(COALESCE(createdatecatebank,now()), 'YYYY-MM-DD HH24:MI:SS'), 
			updatecatebank, to_char(COALESCE(updatedatecatebank,now()), 'YYYY-MM-DD HH24:MI:SS') 
			FROM ` + database_catebank_local + `  
			ORDER BY createdatecatebank DESC   `

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idcatebank_db                                                                      int
			nmcatebank_db, statuscatebank_db                                                   string
			createcatebank_db, createdatecatebank_db, updatecatebank_db, updatedatecatebank_db string
		)

		err = row.Scan(&idcatebank_db, &nmcatebank_db, &statuscatebank_db,
			&createcatebank_db, &createdatecatebank_db, &updatecatebank_db, &updatedatecatebank_db)

		helpers.ErrorCheck(err)
		create := ""
		update := ""
		status_css := configs.STATUS_CANCEL
		if createcatebank_db != "" {
			create = createcatebank_db + ", " + createdatecatebank_db
		}
		if updatecatebank_db != "" {
			update = updatecatebank_db + ", " + updatedatecatebank_db
		}
		if statuscatebank_db == "Y" {
			status_css = configs.STATUS_COMPLETE
		}

		//BANKTYPE
		var objbanktype entities.Model_bankType
		var arraobjbanktype []entities.Model_bankType
		sql_selectbanktype := `SELECT 
			idbanktype, nmbanktype, imgbanktype, statusbanktype,  
			createbanktype, to_char(COALESCE(createdatebanktype,now()), 'YYYY-MM-DD HH24:MI:SS'), 
			updatebanktype, to_char(COALESCE(updatedatebanktype,now()), 'YYYY-MM-DD HH24:MI:SS') 
			FROM ` + database_banktype_local + ` 
			WHERE idcatebank = $1   
		`
		row_banktype, err := con.QueryContext(ctx, sql_selectbanktype, idcatebank_db)
		helpers.ErrorCheck(err)
		for row_banktype.Next() {
			var (
				idbanktype_db, nmbanktype_db, imgbanktype_db, statusbanktype_db                    string
				createbanktype_db, createdatebanktype_db, updatebanktype_db, updatedatebanktype_db string
			)
			err = row_banktype.Scan(&idbanktype_db, &nmbanktype_db, &imgbanktype_db, &statusbanktype_db,
				&createbanktype_db, &createdatebanktype_db, &updatebanktype_db, &updatedatebanktype_db)

			create_detail := ""
			update_detail := ""
			status_css_detail := configs.STATUS_CANCEL
			if createbanktype_db != "" {
				create_detail = createbanktype_db + ", " + createdatebanktype_db
			}
			if updatebanktype_db != "" {
				update_detail = updatebanktype_db + ", " + updatedatebanktype_db
			}
			if statusbanktype_db == "Y" {
				status_css_detail = configs.STATUS_COMPLETE
			}

			objbanktype.Banktype_id = idbanktype_db
			objbanktype.Banktype_name = nmbanktype_db
			objbanktype.Banktype_img = imgbanktype_db
			objbanktype.Banktype_status = statusbanktype_db
			objbanktype.Banktype_status_css = status_css_detail
			objbanktype.Banktype_create = create_detail
			objbanktype.Banktype_update = update_detail
			arraobjbanktype = append(arraobjbanktype, objbanktype)
		}

		obj.Catebank_id = idcatebank_db
		obj.Catebank_name = nmcatebank_db
		obj.Catebank_status = statuscatebank_db
		obj.Catebank_status_css = status_css
		obj.Catebank_list = arraobjbanktype
		obj.Catebank_create = create
		obj.Catebank_update = update
		arraobj = append(arraobj, obj)
		msg = "Success"
	}
	defer row.Close()

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = arraobj
	res.Time = time.Since(start).String()

	return res, nil
}
func Save_catebank(admin, name, status, sData string, idrecord int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()
	if sData == "New" {
		sql_insert := `
				insert into
				` + database_catebank_local + ` (
					idcatebank , nmcatebank, statuscatebank, 
					createcatebank, createdatecatebank  
				) values (
					$1, $2, $3,   
					$4, $5
				)
			`
		field_column := database_catebank_local + tglnow.Format("YYYY")
		idrecord_counter := Get_counter(field_column)
		flag_insert, msg_insert := Exec_SQL(sql_insert, database_catebank_local, "INSERT",
			tglnow.Format("YY")+strconv.Itoa(idrecord_counter), name, status,
			admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

		if flag_insert {
			msg = "Succes"
		} else {
			fmt.Println(msg_insert)
		}
	} else {
		sql_update := `
				UPDATE 
				` + database_catebank_local + `  
				SET nmcatebank=$1, statuscatebank=$2,   
				updatecatebank=$3, updatedatecatebank=$4    
				WHERE idcatebank=$5   
			`

		flag_update, msg_update := Exec_SQL(sql_update, database_catebank_local, "UPDATE",
			name, status,
			admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idrecord)

		if flag_update {
			msg = "Succes"
		} else {
			fmt.Println(msg_update)
		}
	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = nil
	res.Time = time.Since(render_page).String()

	return res, nil
}
func Save_banktype(admin, idrecord, name, img, status, sData string, idcatebank int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()
	flag := false

	if sData == "New" {
		flag = CheckDB(database_banktype_local, "idbanktype", idrecord)
		if !flag {
			sql_insert := `
				insert into
				` + database_banktype_local + ` (
					idbanktype, idcatebank, nmbanktype, imgbanktype, statusbanktype, 
					createbanktype, createdatebanktype  
				) values (
					$1, $2, $3, $4, $5,  
					$6, $7
				)
			`
			flag_insert, msg_insert := Exec_SQL(sql_insert, database_banktype_local, "INSERT",
				idrecord, idcatebank, name, img, status,
				admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

			if flag_insert {
				msg = "Succes"
			} else {
				fmt.Println(msg_insert)
			}
		} else {
			msg = "Duplicate Entry"
		}
	} else {
		sql_update := `
				UPDATE 
				` + database_banktype_local + `  
				SET nmbanktype=$1, imgbanktype=$2, statusbanktype=$3,    
				updatebanktype=$4, updatedatebanktype=$5    
				WHERE idbanktype=$6   
			`

		flag_update, msg_update := Exec_SQL(sql_update, database_banktype_local, "UPDATE",
			name, img, status,
			admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idrecord)

		if flag_update {
			msg = "Succes"
		} else {
			fmt.Println(msg_update)
		}
	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = nil
	res.Time = time.Since(render_page).String()

	return res, nil
}
