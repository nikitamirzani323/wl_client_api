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

const database_agenbank_local = configs.DB_tbl_mst_master_agen_bank

func Fetch_agenbankHome(idmasteragen string) (helpers.Responsemember, error) {
	var obj entities.Model_agenbank
	var arraobj []entities.Model_agenbank
	var objbanktype entities.Model_bankTypeshare
	var arraobjbanktype []entities.Model_bankTypeshare
	var res helpers.Responsemember
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
			idagenbank , type_agenbank, 
			idbanktype,  norekbank, nmownerbank, status_agenbank,
			create_agenbank, to_char(COALESCE(createdate_agenbank,now()), 'YYYY-MM-DD HH24:MI:SS'), 
			update_agenbank, to_char(COALESCE(updatedate_agenbank,now()), 'YYYY-MM-DD HH24:MI:SS') 
			FROM ` + database_agenbank_local + `  
			WHERE idmasteragen=$1 
			ORDER BY status_agenbank ASC   `

	row, err := con.QueryContext(ctx, sql_select, idmasteragen)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idagenbank_db                                                                          int
			type_agenbank_db, idbanktype_db, norekbank_db, nmownerbank_db, status_agenbank_db      string
			create_agenbank_db, createdate_agenbank_db, update_agenbank_db, updatedate_agenbank_db string
		)

		err = row.Scan(&idagenbank_db, &type_agenbank_db, &idbanktype_db, &norekbank_db, &nmownerbank_db, &status_agenbank_db,
			&create_agenbank_db, &createdate_agenbank_db, &update_agenbank_db, &updatedate_agenbank_db)

		helpers.ErrorCheck(err)
		create := ""
		update := ""
		status_css := configs.STATUS_CANCEL
		if create_agenbank_db != "" {
			create = create_agenbank_db + ", " + createdate_agenbank_db
		}
		if update_agenbank_db != "" {
			update = update_agenbank_db + ", " + updatedate_agenbank_db
		}
		if status_agenbank_db == "Y" {
			status_css = configs.STATUS_COMPLETE
		}

		obj.Agenbank_id = idagenbank_db
		obj.Agenbank_tipe = type_agenbank_db
		obj.Agenbank_idbanktype = idbanktype_db
		obj.Agenbank_norek = norekbank_db
		obj.Agenbank_nmrek = nmownerbank_db
		obj.Agenbank_status = status_agenbank_db
		obj.Agenbank_status_css = status_css
		obj.Agenbank_create = create
		obj.Agenbank_update = update
		arraobj = append(arraobj, obj)
		msg = "Success"
	}
	defer row.Close()

	sql_selectbanktype := `SELECT 
			B.nmcatebank, A.idbanktype  
			FROM ` + configs.DB_tbl_mst_banktype + ` as A 
			JOIN ` + configs.DB_tbl_mst_cate_bank + ` as B ON B.idcatebank = A.idcatebank 
			ORDER BY B.nmcatebank,A.idbanktype ASC    
	`
	rowbanktype, errbanktype := con.QueryContext(ctx, sql_selectbanktype)
	helpers.ErrorCheck(errbanktype)
	for rowbanktype.Next() {
		var (
			nmcatebank_db, idbanktype_db string
		)

		errbanktype = rowbanktype.Scan(&nmcatebank_db, &idbanktype_db)

		helpers.ErrorCheck(errbanktype)

		objbanktype.Catebank_name = nmcatebank_db
		objbanktype.Banktype_id = idbanktype_db
		arraobjbanktype = append(arraobjbanktype, objbanktype)
		msg = "Success"
	}
	defer rowbanktype.Close()

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = arraobj
	res.Listbank = arraobjbanktype
	res.Time = time.Since(start).String()

	return res, nil
}
func Fetch_agenbankList(idmasteragen string) (helpers.Response, error) {
	var obj entities.Model_agenbankshare
	var arraobj []entities.Model_agenbankshare
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	perpage := 50

	sql_select := ""
	sql_select += ""
	sql_select += "SELECT "
	sql_select += "idagenbank , idbanktype, norekbank, nmownerbank "
	sql_select += "FROM " + database_agenbank_local + "  "
	sql_select += "WHERE idmasteragen = '" + idmasteragen + "' "
	sql_select += "AND status_agenbank='Y' "
	sql_select += "ORDER BY idbanktype ASC   LIMIT " + strconv.Itoa(perpage)

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idagenbank_db                               int
			idbanktype_db, norekbank_db, nmownerbank_db string
		)

		err = row.Scan(&idagenbank_db, &idbanktype_db, &norekbank_db, &nmownerbank_db)

		helpers.ErrorCheck(err)

		obj.Agenbank_id = idagenbank_db
		obj.Agenbank_info = idbanktype_db + "-" + norekbank_db + "-" + nmownerbank_db
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

func Save_agenbank(admin, idmasteragen, tipe, idbank, norek, nmrek, status, sData string, idrecord int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()

	if sData == "New" {
		sql_insert := `
				insert into
				` + database_agenbank_local + ` (
					idagenbank , idmasteragen, type_agenbank,
					idbanktype, norekbank, nmownerbank, status_agenbank, 
					create_agenbank, createdate_agenbank    
				) values (
					$1, $2, $3,   
					$4, $5, $6, $7,   
					$8, $9
				)
			`
		field_column := database_agenbank_local + tglnow.Format("YYYY")
		idrecord_counter := Get_counter(field_column)
		flag_insert, msg_insert := Exec_SQL(sql_insert, database_agenbank_local, "INSERT",
			tglnow.Format("YY")+strconv.Itoa(idrecord_counter), idmasteragen, tipe,
			idbank, norek, nmrek, status,
			admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

		if flag_insert {
			msg = "Succes"
		} else {
			fmt.Println(msg_insert)
		}
	} else {
		sql_update := `
				UPDATE 
				` + database_agenbank_local + `  
				SET type_agenbank=$1, idbanktype=$2, norekbank=$3, nmownerbank=$4, 
				status_agenbank=$5,     
				update_agenbank=$6, updatedate_agenbank=$7       
				WHERE idagenbank=$8 AND idmasteragen=$9     
			`

		flag_update, msg_update := Exec_SQL(sql_update, database_agenbank_local, "UPDATE",
			tipe, idbank, norek, nmrek, status,
			admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idrecord, idmasteragen)

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
