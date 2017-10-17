{{- $short := (shortname .Name "err" "res" "sqlstr" "db" "XOLog") -}}
{{- $table := (schema .Schema .Table.TableName) -}}
{{- $typ := .Name}}
{{- if .Comment -}}
// {{ .Comment }}
{{- else -}}
// {{ .Name }} represents a row from '{{ $table }}'.
{{- end }}

// Manualy copy this to project
type {{ .Name }}__ struct {
{{- range .Fields }}
	{{ .Col.ColumnName }} {{ retype .Type }} `json:"{{ .Col.ColumnName }}"` // {{ .Col.ColumnName }} -
{{- end }}
{{- if .PrimaryKey }}

	// xo fields
	_exists, _deleted bool
{{ end }}
}


{{ if .PrimaryKey }}
// Exists determines if the {{ .Name }} exists in the database.
func ({{ $short }} *{{ .Name }}) Exists() bool {
	return {{ $short }}._exists
}

// Deleted provides information if the {{ .Name }} has been deleted from the database.
func ({{ $short }} *{{ .Name }}) Deleted() bool {
	return {{ $short }}._deleted
}

// Insert inserts the {{ .Name }} to the database.
func ({{ $short }} *{{ .Name }}) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if {{ $short }}._exists {
		return errors.New("insert failed: already exists")
	}

{{ if .Table.ManualPk  }}
	// sql insert query, primary key must be provided
	const sqlstr = `INSERT INTO {{ $table }} (` +
		`{{ colnames .Fields }}` +
		`) VALUES (` +
		`{{ colvals .Fields }}` +
		`)`

	// run query
	XOLog(sqlstr, {{ fieldnames .Fields $short }})
	_, err = db.Exec(sqlstr, {{ fieldnames .Fields $short }})
	if err != nil {
		return err
	}

	// set existence
	{{ $short }}._exists = true
{{ else }}
	// sql insert query, primary key provided by autoincrement
	const sqlstr = `INSERT INTO {{ $table }} (` +
		`{{ colnames .Fields .PrimaryKey.Name }}` +
		`) VALUES (` +
		`{{ colvals .Fields .PrimaryKey.Name }}` +
		`)`

	// run query
	XOLog(sqlstr, {{ fieldnames .Fields $short .PrimaryKey.Name }})
	res, err := db.Exec(sqlstr, {{ fieldnames .Fields $short .PrimaryKey.Name }})
	if err != nil {
		XOLogErr(err)
		return err
	}

	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		XOLogErr(err)
		return err
	}

	// set primary key and existence
	{{ $short }}.{{ .PrimaryKey.Name }} = {{ .PrimaryKey.Type }}(id)
	{{ $short }}._exists = true
{{ end }}

	On{{ .Name }}_AfterInsert({{ $short }})

	return nil
}

// Insert inserts the {{ .Name }} to the database.
func ({{ $short }} *{{ .Name }}) Replace(db XODB) error {
	var err error

	// sql query
{{ if .Table.ManualPk  }}
	const sqlstr = `REPLACE INTO {{ $table }} (` +
		`{{ colnames .Fields }}` +
		`) VALUES (` +
		`{{ colvals .Fields }}` +
		`)`

	// run query
	XOLog(sqlstr, {{ fieldnames .Fields $short  }})
	_, err = db.Exec(sqlstr, {{ fieldnames .Fields $short }})
	if err != nil {
		XOLogErr(err)
		return err
	}

	{{ $short }}._exists = true
{{else}}
	const sqlstr = `REPLACE INTO {{ $table }} (` +
		`{{ colnames .Fields .PrimaryKey.Name }}` +
		`) VALUES (` +
		`{{ colvals .Fields .PrimaryKey.Name }}` +
		`)`

	// run query
	XOLog(sqlstr, {{ fieldnames .Fields $short .PrimaryKey.Name }})
	res, err := db.Exec(sqlstr, {{ fieldnames .Fields $short .PrimaryKey.Name }})
	if err != nil {
		XOLogErr(err)
		return err
	}

	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		XOLogErr(err)
		return err
	}

	// set primary key and existence
	{{ $short }}.{{ .PrimaryKey.Name }} = {{ .PrimaryKey.Type }}(id)
	{{ $short }}._exists = true
{{end}}

	On{{ .Name }}_AfterInsert({{ $short }})

	return nil
}

// Update updates the {{ .Name }} in the database.
func ({{ $short }} *{{ .Name }}) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !{{ $short }}._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if {{ $short }}._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE {{ $table }} SET ` +
		`{{ colnamesquery .Fields ", " .PrimaryKey.Name }}` +
		` WHERE {{ colname .PrimaryKey.Col }} = ?`

	// run query
	XOLog(sqlstr, {{ fieldnames .Fields $short .PrimaryKey.Name }}, {{ $short }}.{{ .PrimaryKey.Name }})
	_, err = db.Exec(sqlstr, {{ fieldnames .Fields $short .PrimaryKey.Name }}, {{ $short }}.{{ .PrimaryKey.Name }})

	XOLogErr(err)
	On{{ .Name }}_AfterUpdate({{ $short }})

	return err
}

// Save saves the {{ .Name }} to the database.
func ({{ $short }} *{{ .Name }}) Save(db XODB) error {
	if {{ $short }}.Exists() {
		return {{ $short }}.Update(db)
	}

	return {{ $short }}.Replace(db)
}

// Delete deletes the {{ .Name }} from the database.
func ({{ $short }} *{{ .Name }}) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !{{ $short }}._exists {
		return nil
	}

	// if deleted, bail
	if {{ $short }}._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM {{ $table }} WHERE {{ colname .PrimaryKey.Col }} = ?`

	// run query
	XOLog(sqlstr, {{ $short }}.{{ .PrimaryKey.Name }})
	_, err = db.Exec(sqlstr, {{ $short }}.{{ .PrimaryKey.Name }})
	if err != nil {
		XOLogErr(err)
		return err
	}

	// set deleted
	{{ $short }}._deleted = true

	On{{ .Name }}_AfterDelete({{ $short }})

	return nil
}

////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////// Querify gen - ME /////////////////////////////////////////
//.Name = table name
{{- $deleterType := printf "__%s_Deleter" .Name }}
{{- $updaterType := printf "__%s_Updater" .Name }}
{{- $selectorType := printf "__%s_Selector" .Name }}
{{- $updater := printf "__%s_Updater" .Name }}
{{ $ms_gen_types := ms_gen_types }} // _Deleter, _Updater

// orma types
type {{ $deleterType }} struct {
	wheres   []whereClause
    whereSep string
}

type {{ $updaterType }} struct {
	wheres   []whereClause
	updates   map[string]interface{}
    whereSep string
}

type {{ $selectorType }} struct {
    wheres   []whereClause
    selectCol string
    whereSep  string
    orderBy string//" order by id desc //for ints
    limit int
    offset int
}

func New{{ .Name }}_Deleter()  *{{ $deleterType }} {
	    d := {{ $deleterType }} {whereSep: " AND "}
	    return &d
}

func New{{ .Name }}_Updater()  *{{ $updaterType }} {
	    u := {{ $updaterType }} {whereSep: " AND "}
	    u.updates =  make(map[string]interface{},10)
	    return &u
}

func New{{ .Name }}_Selector()  *{{ $selectorType }} {
	    u := {{ $selectorType }} {whereSep: " AND ",selectCol: "*"}
	    return &u
}


{{- $ms_cond_list := ms_conds }}
{{- $ms_str_cond := ms_str_cond }}
{{- $ms_in := ms_in }}
{{- $Fields := .Fields }}
/////////////////////////////// Where for all /////////////////////////////
//// for ints all selector updater, deleter
{{ range (ms_to_slice $deleterType $updaterType $selectorType) }}
		{{ $operationType := . }}
			////////ints
func (u *{{$operationType}}) Or () *{{$operationType}} {
    u.whereSep = " OR "
    return u
}		
		{{- range $Fields }}
			
			{{- $colName := .Col.ColumnName }}
			{{- $colType := .Type }}
	
				{{- if (or (eq $colType "int64") (eq $colType "int") ) }}
				
func (u *{{$operationType}}) {{ $colName }}_In (ins []int) *{{$operationType}} {
    w := whereClause{}
    var insWhere []interface{}
    for _, i:= range ins {
        insWhere = append(insWhere,i)
    }
    w.args = insWhere
    w.condition = " {{ $colName }} IN("+helper.DbQuestionForSqlIn(len(ins))+") "
    u.wheres = append(u.wheres, w)

    return u
}

func (u *{{$operationType}}) {{ $colName }}_Ins (ins ...int) *{{$operationType}} {
    w := whereClause{}
    var insWhere []interface{}
    for _, i:= range ins {
        insWhere = append(insWhere,i)
    }
    w.args = insWhere
    w.condition = " {{ $colName }} IN("+helper.DbQuestionForSqlIn(len(ins))+") "
    u.wheres = append(u.wheres, w)

    return u
}

func (u *{{$operationType}}) {{ $colName }}_NotIn (ins []int) *{{$operationType}} {
    w := whereClause{}
    var insWhere []interface{}
    for _, i:= range ins {
        insWhere = append(insWhere,i)
    }
    w.args = insWhere
    w.condition = " {{ $colName }} NOT IN("+helper.DbQuestionForSqlIn(len(ins))+") "
    u.wheres = append(u.wheres, w)

    return u
}

					{{- with $ms_cond_list }}
						{{- range  .  }}

func (d *{{$operationType}}) {{ $colName }}{{ .Suffix }} (val int) *{{$operationType}} {
    w := whereClause{}
    var insWhere []interface{}
    insWhere = append(insWhere,val)
    w.args = insWhere
    w.condition = " {{ $colName }} {{.Condition}} ? "
    d.wheres = append(d.wheres, w)
    	
    return d
}
						{{- end }}
					{{- end }}
				
				{{- end }}

		{{- end }}

{{ end }}


///// for strings //copy of above with type int -> string + rm if eq + $ms_str_cond
{{ range (ms_to_slice $deleterType $updaterType $selectorType) }}
		{{ $operationType := . }}
			////////ints
		{{- range $Fields }}
			
			{{- $colName := .Col.ColumnName }}
			{{- $colType := .Type }}
	
				{{- if (eq $colType "string") }}
				
func (u *{{$operationType}}) {{ $colName }}_In (ins []string) *{{$operationType}} {
    w := whereClause{}
    var insWhere []interface{}
    for _, i:= range ins {
        insWhere = append(insWhere,i)
    }
    w.args = insWhere
    w.condition = " {{ $colName }} IN("+helper.DbQuestionForSqlIn(len(ins))+") "
    u.wheres = append(u.wheres, w)

    return u
}

func (u *{{$operationType}}) {{ $colName }}_NotIn (ins []string) *{{$operationType}} {
    w := whereClause{}
    var insWhere []interface{}
    for _, i:= range ins {
        insWhere = append(insWhere,i)
    }
    w.args = insWhere
    w.condition = " {{ $colName }} NOT IN("+helper.DbQuestionForSqlIn(len(ins))+") "
    u.wheres = append(u.wheres, w)

    return u
}

//must be used like: UserName_like("hamid%")
func (u *{{$operationType}}) {{ $colName }}_Like (val string) *{{$operationType}} {
    w := whereClause{}
    var insWhere []interface{}
    insWhere = append(insWhere,val)
    w.args = insWhere
    w.condition = " {{ $colName }} LIKE ? "
    u.wheres = append(u.wheres, w)

    return u
}

					{{- with $ms_str_cond }}
						{{- range  .  }}

func (d *{{$operationType}}) {{ $colName }}{{ .Suffix }} (val string) *{{$operationType}} {
    w := whereClause{}
    var insWhere []interface{}
    insWhere = append(insWhere,val)
    w.args = insWhere
    w.condition = " {{ $colName }} {{.Condition}} ? "
    d.wheres = append(d.wheres, w)
    	
    return d
}
						{{- end }}
					{{- end }}
				
				{{- end }}

		{{- end }}

{{ end }}
/// End of wheres for selectors , updators, deletor

/////////////////////////////// Updater /////////////////////////////

{{ $operationType := $updaterType }}

{{- range $Fields }}
			
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
{{- range $Fields }}
			
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
	//s:= "({{ ms_question_mark .Fields }})," //`(?, ?, ?, ?),`
	s:= "({{ ms_question_mark .Fields }})," //`(?, ?, ?, ?),`
	insVals_:= strings.Repeat(s, ln)
	insVals := insVals_[0:len(insVals_)-1]
	// sql query
	sqlstr := "INSERT INTO {{ $table }} (" +
		"{{ colnames .Fields  }}" +
		") VALUES " + insVals

	// run query
	vals := make([]interface{},0, ln * 5)//5 fields
	
	for _,row := range rows {
		// vals = append(vals,row.UserId)
		{{ ms_append_fieldnames .Fields "vals" }}
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
	s:= "({{ ms_question_mark .Fields }})," //`(?, ?, ?, ?),`
	insVals_:= strings.Repeat(s, ln)
	insVals := insVals_[0:len(insVals_)-1]
	// sql query
	sqlstr := "REPLACE INTO {{ $table }} (" +
		"{{ colnames .Fields }}" +
		") VALUES " + insVals

	// run query
	vals := make([]interface{},0, ln * 5)//5 fields
	
	for _,row := range rows {
		// vals = append(vals,row.UserId)
		{{ ms_append_fieldnames .Fields "vals" }}
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
	//s:= "({{ ms_question_mark .Fields .PrimaryKey.Name }})," //`(?, ?, ?, ?),`
	s:= "({{ ms_question_mark .Fields .PrimaryKey.Name }})," //`(?, ?, ?, ?),`
	insVals_:= strings.Repeat(s, ln)
	insVals := insVals_[0:len(insVals_)-1]
	// sql query
	sqlstr := "INSERT INTO {{ $table }} (" +
		"{{ colnames .Fields .PrimaryKey.Name }}" +
		") VALUES " + insVals

	// run query
	vals := make([]interface{},0, ln * 5)//5 fields
	
	for _,row := range rows {
		// vals = append(vals,row.UserId)
		{{ ms_append_fieldnames .Fields "vals" .PrimaryKey.Name }}
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
	s:= "({{ ms_question_mark .Fields .PrimaryKey.Name }})," //`(?, ?, ?, ?),`
	insVals_:= strings.Repeat(s, ln)
	insVals := insVals_[0:len(insVals_)-1]
	// sql query
	sqlstr := "REPLACE INTO {{ $table }} (" +
		"{{ colnames .Fields .PrimaryKey.Name }}" +
		") VALUES " + insVals

	// run query
	vals := make([]interface{},0, ln * 5)//5 fields
	
	for _,row := range rows {
		// vals = append(vals,row.UserId)
		{{ ms_append_fieldnames .Fields "vals" .PrimaryKey.Name }}
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
{{- range $Fields }}
			
			{{- $colName := .Col.ColumnName }}
			{{- $colType := .Type }}

			// {{- /* $colType }} {{ $colName */}}

{{- end}}





{{- end }}

