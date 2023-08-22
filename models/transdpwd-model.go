package models

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nikitamirzani323/wl_agen_backend_api/configs"
	"github.com/nikitamirzani323/wl_agen_backend_api/db"
	"github.com/nikitamirzani323/wl_agen_backend_api/entities"
	"github.com/nikitamirzani323/wl_agen_backend_api/helpers"
	"github.com/nleeper/goment"
)

func Fetch_transdpwdHome(idmasteragen string) (helpers.Response, error) {
	var obj entities.Model_transdpwd
	var arraobj []entities.Model_transdpwd
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	tbl_trx_dpwd, _ := Get_mappingdatabase(idmasteragen)
	sql_select := `SELECT 
			iddpwd , date_dpwd, idcurr,  
			tipedocuser_dpwd, tipedoc_dpwd , tipeakun_dpwd, idagenmember,  ipaddress_dpwd, timezone_dpwd,  
			bank_in, bank_in_info , bank_out, bank_out_info, 
			round(amount_dpwd*multiplier_dpwd) as amount_dpwd , before_dpwd, after_dpwd,  status_dpwd,
			create_dpwd, to_char(COALESCE(createdate_dpwd,now()), 'YYYY-MM-DD HH24:MI:SS'), 
			update_dpwd, to_char(COALESCE(updatedate_dpwd,now()), 'YYYY-MM-DD HH24:MI:SS') 
			FROM ` + tbl_trx_dpwd + `  
			ORDER BY createdate_dpwd DESC   `

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			iddpwd_db, date_dpwd_db, idcurr_db                                                                                                                              string
			tipedocuser_dpwd_db, tipedoc_dpwd_db, tipeakun_dpwd_db, idagenmember_db, ipaddress_dpwd_db, timezone_dpwd_db, bank_in_info_db, bank_out_info_db, status_dpwd_db string
			bank_in_db, bank_out_db                                                                                                                                         int
			amount_dpwd_db, before_dpwd_db, after_dpwd_db                                                                                                                   float64
			create_dpwd_db, createdate_dpwd_db, update_dpwd_db, updatedate_dpwd_db                                                                                          string
		)

		err = row.Scan(&iddpwd_db, &date_dpwd_db, &idcurr_db,
			&tipedocuser_dpwd_db, &tipedoc_dpwd_db, &tipeakun_dpwd_db, &idagenmember_db, &ipaddress_dpwd_db, &timezone_dpwd_db,
			&bank_in_db, &bank_in_info_db, &bank_out_db, &bank_out_info_db,
			&amount_dpwd_db, &before_dpwd_db, &after_dpwd_db, &status_dpwd_db,
			&create_dpwd_db, &createdate_dpwd_db, &update_dpwd_db, &updatedate_dpwd_db)

		helpers.ErrorCheck(err)
		create := ""
		update := ""
		status_css := configs.STATUS_CANCEL
		if create_dpwd_db != "" {
			create = create_dpwd_db + ", " + createdate_dpwd_db
		}
		if update_dpwd_db != "" {
			update = update_dpwd_db + ", " + updatedate_dpwd_db
		}
		switch status_dpwd_db {
		case "PROCESS":
			status_css = configs.STATUS_RUNNING
		case "APPROVED":
			status_css = configs.STATUS_COMPLETE
		case "REJECT":
			status_css = configs.STATUS_CANCEL
		}

		obj.Transdpwd_id = iddpwd_db
		obj.Transdpwd_date = date_dpwd_db
		obj.Transdpwd_idcurr = idcurr_db
		obj.Transdpwd_tipedoc = tipedoc_dpwd_db
		obj.Transdpwd_tipeuserdoc = tipedocuser_dpwd_db
		obj.Transdpwd_tipeakun = tipeakun_dpwd_db
		obj.Transdpwd_idmember = idagenmember_db
		obj.Transdpwd_nmmember = _GetInfoMember(idmasteragen, idagenmember_db)
		obj.Transdpwd_ipaddress = ipaddress_dpwd_db
		obj.Transdpwd_timezone = timezone_dpwd_db
		obj.Transdpwd_bank_in = bank_in_db
		obj.Transdpwd_bank_in_info = bank_in_info_db
		obj.Transdpwd_bank_out = bank_out_db
		obj.Transdpwd_bank_out_info = bank_out_info_db
		obj.Transdpwd_amount = amount_dpwd_db
		obj.Transdpwd_before = before_dpwd_db
		obj.Transdpwd_after = after_dpwd_db
		obj.Transdpwd_status = status_dpwd_db
		obj.Transdpwd_status_css = status_css
		obj.Transdpwd_create = create
		obj.Transdpwd_update = update
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
func Save_transdpwd(admin, idrecord, idmasteragen, idmaster, tipedoc, idmember, note_dpwd, status, sData string, bank_in, bank_out int, amount float32) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()
	tbl_trx_dpwd, _ := Get_mappingdatabase(idmasteragen)

	idcurr := _GetDefaultCurr(idmasteragen)
	multiplier := _GetMultiplier(idcurr)
	before := 0
	after := 0
	log.Println(tbl_trx_dpwd)
	log.Println(idcurr)
	log.Println(multiplier)
	if sData == "New" {
		sql_insert := `
				insert into
				` + tbl_trx_dpwd + ` (
					iddpwd , idmasteragen, idmaster, 
					yearmonth_dpwd , date_dpwd, idcurr, tipedocuser_dpwd, tipedoc_dpwd, tipeakun_dpwd, idagenmember, 
					bank_in, bank_in_info , bank_out, bank_out_info, 
					multiplier_dpwd, amountdefault_dpwd, amount_dpwd, before_dpwd, after_dpwd, status_dpwd, note_dpwd, 
					create_dpwd, createdate_dpwd  
				) values (
					$1, $2, $3,   
					$4, $5, $6, $7, $8, $9, $10,    
					$11, $12, $13, $14,     
					$15, $16, $17, $18, $19, $20, $21,      
					$22, $23
				)
			`
		tipeakun_dpwd := ""
		temp_bank_out := ""
		temp_bank_in := ""
		switch tipedoc {
		case "DEPOSIT":
			tipeakun_dpwd = "IN"
			temp_bank_out = _GetInfoBank(idmasteragen, idmember, "MEMBER", bank_out)
			temp_bank_in = _GetInfoBank(idmasteragen, idmember, "AGEN", bank_in)
		case "WITHDRAW":
			tipeakun_dpwd = "OUT"
			temp_bank_out = _GetInfoBank(idmasteragen, idmember, "AGEN", bank_out)
			temp_bank_in = _GetInfoBank(idmasteragen, idmember, "MEMBER", bank_in)
		case "BONUS":
			tipeakun_dpwd = "OUT"
			temp_bank_out = _GetInfoBank(idmasteragen, idmember, "AGEN", bank_out)
			temp_bank_in = _GetInfoBank(idmasteragen, idmember, "MEMBER", bank_in)
			tipeakun_dpwd = "OUT"
		}
		amount_db := amount / multiplier

		field_column := tbl_trx_dpwd + tglnow.Format("YYYY-MM")
		idrecord_counter := Get_counter(field_column)
		iddpwd := idmasteragen + "DPWD" + tglnow.Format("YY") + tglnow.Format("MM") + tglnow.Format("DD") + tglnow.Format("HH") + strconv.Itoa(idrecord_counter)

		flag_insert, msg_insert := Exec_SQL(sql_insert, tbl_trx_dpwd, "INSERT",
			iddpwd, idmasteragen, idmaster,
			tglnow.Format("YYYY-MM"), tglnow.Format("YYYY-MM-DD"), idcurr, "A", tipedoc, tipeakun_dpwd, idmember,
			bank_in, temp_bank_in, bank_out, temp_bank_out,
			multiplier, amount, amount_db, before, after, status, note_dpwd,
			admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

		if flag_insert {
			msg = "Succes"
		} else {
			fmt.Println(msg_insert)
		}
	} else {
		sql_update := `
				UPDATE 
				` + tbl_trx_dpwd + `  
				SET tipedoc_dpwd=$1, tipeakun_dpwd=$2,   
				idagenmember=$3, bank_int=$4, bank_out=$5, note_bank=$6,      
				amount_dpwd=$7, before_dpwd=$8, after_dpwd=$9, status_dpwd=$10, note_dpwd=$11,     
				update_dpwd=$12, updatedate_dpwd=$13     
				WHERE iddpwd=$14  AND idmasteragen=$15  
			`

		tipeakun_dpwd := ""
		note_bank := "FROM: BANK OUT - TO: BANK IN"
		switch tipedoc {
		case "DEPOSIT":
			tipeakun_dpwd = "IN"
		case "WITHDRAW":
			tipeakun_dpwd = "OUT"
		case "BONUS":
			tipeakun_dpwd = "OUT"
		}
		before := 0
		after := 0

		flag_update, msg_update := Exec_SQL(sql_update, tbl_trx_dpwd, "UPDATE",
			tipedoc, tipeakun_dpwd,
			idmember, bank_in, bank_out, note_bank,
			amount, before, after, status, note_dpwd,
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

func _GetDefaultCurr(idrecord string) string {
	con := db.CreateCon()
	ctx := context.Background()
	idcurr := ""

	sql_select := `SELECT
		idcurr   
		FROM ` + configs.DB_tbl_mst_master_agen + `  
		WHERE idmasteragen = $1 
	`
	row := con.QueryRowContext(ctx, sql_select, idrecord)
	switch e := row.Scan(&idcurr); e {
	case sql.ErrNoRows:
	case nil:
	default:
		helpers.ErrorCheck(e)
	}
	return idcurr
}
func _GetMultiplier(idrecord string) float32 {
	con := db.CreateCon()
	ctx := context.Background()
	multipliercurr := 0

	sql_select := `SELECT
		multipliercurr   
		FROM ` + configs.DB_tbl_mst_curr + `  
		WHERE idcurr = $1 
	`
	row := con.QueryRowContext(ctx, sql_select, idrecord)
	switch e := row.Scan(&multipliercurr); e {
	case sql.ErrNoRows:
	case nil:
	default:
		helpers.ErrorCheck(e)
	}
	return float32(multipliercurr)
}
func _GetInfoBank(idmasteragen, idagenmember, tipe string, idrecord int) string {
	con := db.CreateCon()
	ctx := context.Background()
	info := ""
	bank_id := ""
	bank_norek := ""
	bank_nmrek := ""
	sql_select := ""
	if tipe == "AGEN" {
		sql_select = `SELECT
			idbanktype, norekbank, nmownerbank   
			FROM ` + configs.DB_tbl_mst_master_agen_bank + `  
			WHERE idagenbank=` + strconv.Itoa(idrecord) + ` AND idmasteragen='` + idmasteragen + `'    
		`
	} else {
		sql_select = `SELECT
			idbanktype, norekbank_agenmemberbank, nmownerbank_agenmemberbank   
			FROM ` + configs.DB_tbl_mst_master_agen_member_bank + `  
			WHERE idagenmemberbank=` + strconv.Itoa(idrecord) + ` AND idagenmember='` + idagenmember + `'    
		`
	}

	row := con.QueryRowContext(ctx, sql_select)
	switch e := row.Scan(&bank_id, &bank_norek, &bank_nmrek); e {
	case sql.ErrNoRows:
	case nil:
	default:
		helpers.ErrorCheck(e)
	}
	info = bank_id + "-" + bank_norek + "-" + bank_nmrek
	return info
}
func _GetInfoMember(idmasteragen, idagenmember string) string {
	con := db.CreateCon()
	ctx := context.Background()
	username_agenmember_db := ""
	name_agenmember_db := ""

	sql_select := `SELECT
		username_agenmember, name_agenmember   
		FROM ` + configs.DB_tbl_mst_master_agen_member + `  
		WHERE idagenmember=$1 AND idmasteragen=$2
	`
	row := con.QueryRowContext(ctx, sql_select, idagenmember, idmasteragen)
	switch e := row.Scan(&username_agenmember_db, &name_agenmember_db); e {
	case sql.ErrNoRows:
	case nil:
	default:
		helpers.ErrorCheck(e)
	}
	return username_agenmember_db + "-" + name_agenmember_db
}
