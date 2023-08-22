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

const database_admin_local = configs.DB_tbl_mst_master_agen_admin
const database_adminrule_local = configs.DB_tbl_mst_master_agen_admin_rule

func Fetch_adminHome(idmasteragen string) (helpers.ResponseAdmin, error) {
	var obj entities.Model_admin
	var arraobj []entities.Model_admin
	var res helpers.ResponseAdmin
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()
	sql_select := `SELECT 
			A.idagenadmin, A.idagenadminrule, A.tipeagen_admin,
			B.nmagenadminrule, A.usernameagen_admin , A.nameagen_admin, A.phone1agen_admin, A.phone2agen_admin, 
			A.statusagenadmin, to_char(COALESCE(A.lastloginagen_admin,now()), 'YYYY-MM-DD HH24:MI:SS'), A.ipaddress_admin, 
			A.createagenadmin, to_char(COALESCE(A.createdateagenadmin,now()), 'YYYY-MM-DD HH24:MI:SS'), 
			A.updateagenadmin, to_char(COALESCE(A.updatedateagenadmin,now()), 'YYYY-MM-DD HH24:MI:SS')   
			FROM ` + database_admin_local + ` as A 
			JOIN ` + database_adminrule_local + ` as B on B.idagenadminrule = A.idagenadminrule 
			WHERE A.idmasteragen=$1 
			ORDER BY A.lastloginagen_admin DESC 
		`
	row, err := con.QueryContext(ctx, sql_select, idmasteragen)

	var no int = 0
	helpers.ErrorCheck(err)
	for row.Next() {
		no += 1
		var (
			idagenadminrule_db                                                                                                                        int
			idagenadmin_db, tipeagen_admin_db, nmagenadminrule_db, usernameagen_admin_db, nameagen_admin_db, phone1agen_admin_db, phone2agen_admin_db string
			statusagenadmin_db, lastloginagen_admin_db, ipaddress_admin_db                                                                            string
			createagenadmin_db, createdateagenadmin_db, updateagenadmin_db, updatedateagenadmin_db                                                    string
		)

		err = row.Scan(
			&idagenadmin_db, &idagenadminrule_db, &tipeagen_admin_db, &nmagenadminrule_db,
			&usernameagen_admin_db, &nameagen_admin_db, &phone1agen_admin_db, &phone2agen_admin_db,
			&statusagenadmin_db, &lastloginagen_admin_db, &ipaddress_admin_db,
			&createagenadmin_db, &createdateagenadmin_db, &updateagenadmin_db, &updatedateagenadmin_db)

		helpers.ErrorCheck(err)
		create := ""
		update := ""
		status_css := configs.STATUS_CANCEL
		if createagenadmin_db != "" {
			create = createagenadmin_db + ", " + createdateagenadmin_db
		}
		if updateagenadmin_db != "" {
			update = updateagenadmin_db + ", " + updatedateagenadmin_db
		}
		if statusagenadmin_db == "Y" {
			status_css = configs.STATUS_COMPLETE
		}
		obj.Admin_id = idagenadmin_db
		obj.Admin_idrule = idagenadminrule_db
		obj.Admin_tipe = tipeagen_admin_db
		obj.Admin_nmrule = nmagenadminrule_db
		obj.Admin_username = usernameagen_admin_db
		obj.Admin_nama = nameagen_admin_db
		obj.Admin_phone1 = phone1agen_admin_db
		obj.Admin_phone2 = phone2agen_admin_db
		obj.Admin_lastipaddres = ipaddress_admin_db
		obj.Admin_lastlogin = lastloginagen_admin_db
		obj.Admin_status = statusagenadmin_db
		obj.Admin_status_css = status_css
		obj.Admin_create = create
		obj.Admin_update = update
		arraobj = append(arraobj, obj)
		msg = "Success"
	}
	defer row.Close()

	var objRule entities.Model_adminrule
	var arraobjRule []entities.Model_adminrule
	sql_listrule := `SELECT 
		idagenadminrule, nmagenadminrule  	
		FROM ` + configs.DB_tbl_mst_master_agen_admin_rule + ` 
		WHERE idmasteragen=$1 AND nmagenadminrule!='master'
	`
	row_listrule, err_listrule := con.QueryContext(ctx, sql_listrule, idmasteragen)

	helpers.ErrorCheck(err_listrule)
	for row_listrule.Next() {
		var (
			idagenadminrule_db int
			nmagenadminrule_db string
		)

		err = row_listrule.Scan(&idagenadminrule_db, &nmagenadminrule_db)

		helpers.ErrorCheck(err)

		objRule.Adminrule_idruleadmin = int(idagenadminrule_db)
		objRule.Adminrule_nmruleadmin = nmagenadminrule_db
		arraobjRule = append(arraobjRule, objRule)
		msg = "Success"
	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = arraobj
	res.Listrule = arraobjRule
	res.Time = time.Since(start).String()

	return res, nil
}

func Save_adminHome(admin, idrecord, idmasteragen, username, password, nama, phone1, phone2, status, sData string, idrule int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()
	flag := false

	if sData == "New" {
		flag = CheckDBTwoField(database_admin_local, "idmasteragen", idmasteragen, "usernameagen_admin", username)
		if !flag {
			sql_insert := `
				insert into
				` + database_admin_local + ` (
					idagenadmin, idagenadminrule, idmasteragen , tipeagen_admin, usernameagen_admin, passwordagen_admin, lastloginagen_admin,   
					nameagen_admin, phone1agen_admin, phone2agen_admin , statusagenadmin, 
					createagenadmin, createdateagenadmin    
				) values (
					$1, $2, $3, $4, $5, $6, $7,   
					$8, $9, $10, $11,     
					$12, $13  
				)
			`
			field_column := idmasteragen + database_admin_local + tglnow.Format("YY")
			idrecord_counter := Get_counter(field_column)
			hashpass := helpers.HashPasswordMD5(password)
			create_date := tglnow.Format("YYYY-MM-DD HH:mm:ss")
			flag_insert, msg_insert := Exec_SQL(sql_insert, database_admin_local, "INSERT",
				idmasteragen+"-"+tglnow.Format("YY")+strconv.Itoa(idrecord_counter), idrule, idmasteragen, "ADMIN", username, hashpass, create_date,
				nama, phone1, phone2, status,
				admin, create_date)

			if flag_insert {
				msg = "Succes"
			} else {
				fmt.Println(msg_insert)
			}
		} else {
			msg = "Duplicate Entry"
		}
	} else {
		if password == "" {
			sql_update := `
				UPDATE 
				` + database_admin_local + `  
				SET idagenadminrule=$1, nameagen_admin=$2, phone1agen_admin=$3, phone2agen_admin=$4, statusagenadmin=$5,  
				updateagenadmin=$6, updatedateagenadmin=$7         
				WHERE idmasteragen=$8 AND idagenadmin=$9          
			`

			flag_update, msg_update := Exec_SQL(sql_update, database_admin_local, "UPDATE",
				idrule, nama, phone1, phone2, status,
				admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idmasteragen, idrecord)

			if flag_update {
				msg = "Succes"
			} else {
				fmt.Println(msg_update)
			}
		} else {
			hashpass := helpers.HashPasswordMD5(password)
			sql_update2 := `
				UPDATE 
				` + database_admin_local + `   
				SET idagenadminrule=$1, passwordagen_admin=$2, nameagen_admin=$3, phone1agen_admin=$4, phone2agen_admin=$5, statusagenadmin=$6,  
				updateagenadmin=$7, updatedateagenadmin=$8          
				WHERE idmasteragen=$9 AND idagenadmin=$10         
			`
			flag_update, msg_update := Exec_SQL(sql_update2, database_admin_local, "UPDATE",
				idrule, hashpass, nama, phone1, phone2, status,
				admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idmasteragen, idrecord)

			if flag_update {
				msg = "Succes"
			} else {
				fmt.Println(msg_update)
			}
		}
	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = nil
	res.Time = time.Since(render_page).String()

	return res, nil
}
