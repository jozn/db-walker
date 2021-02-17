

/////////////////////////////// Updater /////////////////////////////

{{ $operationType := $updaterType }}

{{- range $Columns }}

	{{- $colName := .Col.ColumnName }}
	{{- $colType := .Type }}

	//ints
	{{- if (or (eq $colType "int64") (eq $colType "int") ) }}

func (u *{{$updaterType}}){{ $colName }} (newVal int) *{{$updaterType}} {
    u.updates[" {{$colName}} = ? "] = newVal
    return u
}

func (u *{{$updaterType}}){{ $colName }}_Increment (count int) *{{$updaterType}} {
	if count > 0 {
		u.updates[" {{$colName}} = {{$colName}}+? "] = count
	}

	if count < 0 {
		u.updates[" {{$colName}} = {{$colName}}-? "] = -(count) //make it positive
	}

    return u
}
	{{- end }}

	//string
	{{- if (eq $colType "string") }}
func (u *{{$updaterType}}){{ $colName }} (newVal string) *{{$updaterType}} {
    u.updates[" {{$colName}} = ? "] = newVal
    return u
}
	{{- end }}

{{- end }}


/////////////////////////////////////////////////////////////////////
/////////////////////// Selector ///////////////////////////////////
{{ $operationType := $selectorType }}

//Select_* can just be used with: .GetString() , .GetStringSlice(), .GetInt() ..GetIntSlice()
{{- range $Columns }}

	{{- $colName := .Col.ColumnName }}
	{{- $colType := .Type }}

func (u *{{$selectorType}}) OrderBy_{{ $colName }}_Desc () *{{$selectorType}} {
    u.orderBy = " ORDER BY {{$colName}} DESC "
    return u
}

func (u *{{$selectorType}}) OrderBy_{{ $colName }}_Asc () *{{$selectorType}} {
    u.orderBy = " ORDER BY {{$colName}} ASC "
    return u
}

func (u *{{$selectorType}}) Select_{{ $colName }} () *{{$selectorType}} {
    u.selectCol = "{{$colName}}"
    return u
}
{{- end }}

func (u *{{$selectorType}}) Limit(num int) *{{$selectorType}} {
    u.limit = num
    return u
}

func (u *{{$selectorType}}) Offset(num int) *{{$selectorType}} {
    u.offset = num
    return u
}


/////////////////////////  Queryer Selector  //////////////////////////////////
func (u *{{$selectorType}})_stoSql ()  (string,[]interface{}) {
	sqlWherrs, whereArgs := whereClusesToSql(u.wheres,u.whereSep)

	sqlstr := "SELECT " +u.selectCol +" FROM {{ $table }}"

	if len( strings.Trim(sqlWherrs," ") ) > 0 {//2 for safty
		sqlstr += " WHERE "+ sqlWherrs
	}

	if u.orderBy != ""{
        sqlstr += u.orderBy
    }

    if u.limit != 0 {
        sqlstr += " LIMIT " + strconv.Itoa(u.limit)
    }

    if u.offset != 0 {
        sqlstr += " OFFSET " + strconv.Itoa(u.offset)
    }
    return sqlstr,whereArgs
}

func (u *{{$selectorType}}) GetRow (db *sqlx.DB) (*{{ $typ }},error) {
	var err error

	sqlstr, whereArgs := u._stoSql()

	XOLog(sqlstr,whereArgs )

	row := &{{$typ}}{}
	//by Sqlx
	err = db.Get(row ,sqlstr, whereArgs...)
	if err != nil {
		XOLogErr(err)
		return nil, err
	}

	row._exists = true

	On{{ .Name }}_LoadOne(row)

	return row, nil
}

func (u *{{$selectorType}}) GetRows (db *sqlx.DB) ([]*{{ $typ }},error) {
	var err error

	sqlstr, whereArgs := u._stoSql()

	XOLog(sqlstr,whereArgs )

	var rows []*{{$typ}}
	//by Sqlx
	err = db.Unsafe().Select(&rows ,sqlstr, whereArgs...)
	if err != nil {
		XOLogErr(err)
		return nil, err
	}

	for i:=0;i< len(rows);i++ {
		rows[i]._exists = true
	}

	for i:=0;i< len(rows);i++ {
		rows[i]._exists = true
	}

	On{{ .Name }}_LoadMany(rows)

	return rows, nil
}

//dep use GetRows()
func (u *{{$selectorType}}) GetRows2 (db *sqlx.DB) ([]{{ $typ }},error) {
	var err error

	sqlstr, whereArgs := u._stoSql()

	XOLog(sqlstr,whereArgs )

	var rows []*{{$typ}}
	//by Sqlx
	err = db.Unsafe().Select(&rows ,sqlstr, whereArgs...)
	if err != nil {
		XOLogErr(err)
		return nil, err
	}

	for i:=0;i< len(rows);i++ {
		rows[i]._exists = true
	}

	for i:=0;i< len(rows);i++ {
		rows[i]._exists = true
	}

	On{{ .Name }}_LoadMany(rows)

	rows2 := make([]{{$typ}}, len(rows))
	for i:=0;i< len(rows);i++ {
		cp := *rows[i]
		rows2[i]= cp
	}

	return rows2, nil
}



func (u *{{$selectorType}}) GetString (db *sqlx.DB) (string,error) {
	var err error

	sqlstr, whereArgs := u._stoSql()

	XOLog(sqlstr,whereArgs )

	var res string
	//by Sqlx
	err = db.Get(&res ,sqlstr, whereArgs...)
	if err != nil {
		XOLogErr(err)
		return "", err
	}

	return res, nil
}

func (u *{{$selectorType}}) GetStringSlice (db *sqlx.DB) ([]string,error) {
	var err error

	sqlstr, whereArgs := u._stoSql()

	XOLog(sqlstr,whereArgs )

	var rows []string
	//by Sqlx
	err = db.Select(&rows ,sqlstr, whereArgs...)
	if err != nil {
		XOLogErr(err)
		return nil, err
	}

	return rows, nil
}

func (u *{{$selectorType}}) GetIntSlice (db *sqlx.DB) ([]int,error) {
	var err error

	sqlstr, whereArgs := u._stoSql()

	XOLog(sqlstr,whereArgs )

	var rows []int
	//by Sqlx
	err = db.Select(&rows ,sqlstr, whereArgs...)
	if err != nil {
		XOLogErr(err)
		return nil, err
	}

	return rows, nil
}

func (u *{{$selectorType}}) GetInt (db *sqlx.DB) (int,error) {
	var err error

	sqlstr, whereArgs := u._stoSql()

	XOLog(sqlstr,whereArgs )

	var res int
	//by Sqlx
	err = db.Get(&res ,sqlstr, whereArgs...)
	if err != nil {
		XOLogErr(err)
		return 0, err
	}

	return res, nil
}

/////////////////////////  Queryer Update Delete //////////////////////////////////
func (u *{{$updaterType}})Update (db XODB) (int,error) {
    var err error

    var updateArgs []interface{}
    var sqlUpdateArr  []string
    for up, newVal := range u.updates {
        sqlUpdateArr = append(sqlUpdateArr, up)
        updateArgs = append(updateArgs, newVal)
    }
    sqlUpdate:= strings.Join(sqlUpdateArr, ",")

    sqlWherrs, whereArgs := whereClusesToSql(u.wheres,u.whereSep)

    var allArgs []interface{}
    allArgs = append(allArgs, updateArgs...)
    allArgs = append(allArgs, whereArgs...)

    sqlstr := `UPDATE {{ $table }} SET ` + sqlUpdate

    if len( strings.Trim(sqlWherrs," ") ) > 0 {//2 for safty
		sqlstr += " WHERE " +sqlWherrs
	}

    XOLog(sqlstr,allArgs)
    res, err := db.Exec(sqlstr, allArgs...)
    if err != nil {
    	XOLogErr(err)
        return 0,err
    }

    num, err := res.RowsAffected()
    if err != nil {
    	XOLogErr(err)
        return 0,err
    }

    return int(num),nil
}

func (d *{{$deleterType}})Delete (db XODB) (int,error) {
    var err error
    var wheresArr []string
    for _,w := range d.wheres{
        wheresArr = append(wheresArr,w.condition)
    }
    wheresStr := strings.Join(wheresArr, d.whereSep)

    var args []interface{}
    for _,w := range d.wheres{
        args = append(args,w.args...)
    }

    sqlstr := "DELETE FROM {{ $table}} WHERE " + wheresStr

    // run query
    XOLog(sqlstr, args)
    res, err := db.Exec(sqlstr, args...)
    if err != nil {
    	XOLogErr(err)
        return 0,err
    }

    // retrieve id
    num, err := res.RowsAffected()
    if err != nil {
    	XOLogErr(err)
        return 0,err
    }

    return int(num),nil
}

///////////////////////// Mass insert - replace for  {{ .Name }} ////////////////
{{ if .Table.ManualPk  }}
func MassInsert_{{ .Name }}(rows []{{ .Name }} ,db XODB) error {
	if len(rows) == 0 {
		return errors.New("rows slice should not be empty - inserted nothing")
	}
	var err error
	ln := len(rows)
	//s:= "({{ ms_question_mark .Columns }})," //`(?, ?, ?, ?),`
	s:= "({{ ms_question_mark .Columns }})," //`(?, ?, ?, ?),`
	insVals_:= strings.Repeat(s, ln)
	insVals := insVals_[0:len(insVals_)-1]
	// sql query
	sqlstr := "INSERT INTO {{ $table }} (" +
		"{{ colnames .Columns  }}" +
		") VALUES " + insVals

	// run query
	vals := make([]interface{},0, ln * 5)//5 fields

	for _,row := range rows {
		// vals = append(vals,row.UserId)
		{{ ms_append_fieldnames .Columns "vals" }}
	}

	XOLog(sqlstr, " MassInsert len = ", ln, vals)

	_, err = db.Exec(sqlstr, vals...)
	if err != nil {
		XOLogErr(err)
		return err
	}

	return nil
}

func MassReplace_{{ .Name }}(rows []{{ .Name }} ,db XODB) error {
	var err error
	ln := len(rows)
	s:= "({{ ms_question_mark .Columns }})," //`(?, ?, ?, ?),`
	insVals_:= strings.Repeat(s, ln)
	insVals := insVals_[0:len(insVals_)-1]
	// sql query
	sqlstr := "REPLACE INTO {{ $table }} (" +
		"{{ colnames .Columns }}" +
		") VALUES " + insVals

	// run query
	vals := make([]interface{},0, ln * 5)//5 fields

	for _,row := range rows {
		// vals = append(vals,row.UserId)
		{{ ms_append_fieldnames .Columns "vals" }}
	}

	XOLog(sqlstr, " MassReplace len = ", ln , vals)

	_, err = db.Exec(sqlstr, vals...)
	if err != nil {
		XOLogErr(err)
		return err
	}

	return nil
}
{{ else }}

func MassInsert_{{ .Name }}(rows []{{ .Name }} ,db XODB) error {
	if len(rows) == 0 {
		return errors.New("rows slice should not be empty - inserted nothing")
	}
	var err error
	ln := len(rows)
	//s:= "({{ ms_question_mark .Columns .PrimaryKey.ColumnName }})," //`(?, ?, ?, ?),`
	s:= "({{ ms_question_mark .Columns .PrimaryKey.ColumnName }})," //`(?, ?, ?, ?),`
	insVals_:= strings.Repeat(s, ln)
	insVals := insVals_[0:len(insVals_)-1]
	// sql query
	sqlstr := "INSERT INTO {{ $table }} (" +
		"{{ colnames .Columns .PrimaryKey.ColumnName }}" +
		") VALUES " + insVals

	// run query
	vals := make([]interface{},0, ln * 5)//5 fields

	for _,row := range rows {
		// vals = append(vals,row.UserId)
		{{ ms_append_fieldnames .Columns "vals" .PrimaryKey.ColumnName }}
	}

	XOLog(sqlstr, " MassInsert len = ", ln, vals)

	_, err = db.Exec(sqlstr, vals...)
	if err != nil {
		XOLogErr(err)
		return err
	}

	return nil
}

func MassReplace_{{ .Name }}(rows []{{ .Name }} ,db XODB) error {
	var err error
	ln := len(rows)
	s:= "({{ ms_question_mark .Columns .PrimaryKey.ColumnName }})," //`(?, ?, ?, ?),`
	insVals_:= strings.Repeat(s, ln)
	insVals := insVals_[0:len(insVals_)-1]
	// sql query
	sqlstr := "REPLACE INTO {{ $table }} (" +
		"{{ colnames .Columns .PrimaryKey.ColumnName }}" +
		") VALUES " + insVals

	// run query
	vals := make([]interface{},0, ln * 5)//5 fields

	for _,row := range rows {
		// vals = append(vals,row.UserId)
		{{ ms_append_fieldnames .Columns "vals" .PrimaryKey.ColumnName }}
	}

	XOLog(sqlstr, " MassReplace len = ", ln , vals)

	_, err = db.Exec(sqlstr, vals...)
	if err != nil {
		XOLogErr(err)
		return err
	}

	return nil
}

{{ end }}


//////////////////// Play ///////////////////////////////
{{- range $Columns }}

			{{- $colName := .Col.ColumnName }}
			{{- $colType := .Type }}

			// {{- /* $colType }} {{ $colName */}}

{{- end}}





{{- end }}

