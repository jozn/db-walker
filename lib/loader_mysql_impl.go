package lib

import "github.com/jmoiron/sqlx"

// MyTableColumns runs a custom query, returning results as Column.
func MyTableColumns(db *sqlx.DB, schema string, table string) ([]*Column_Impl, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`ordinal_position AS field_ordinal, ` +
		`column_name, ` +
		`IF(data_type = 'enum', column_name, column_type) AS data_type, ` +
		`IF(is_nullable = 'YES', false, true) AS not_null, ` +
		`column_default AS default_value, ` +
		`IF(column_key = 'PRI', true, false) AS is_primary_key, ` +
		`COLUMN_COMMENT ` +
		`FROM information_schema.columns ` +
		`WHERE table_schema = ? AND table_name = ? ` +
		`ORDER BY ordinal_position`

	// run query
	XOLog(sqlstr, schema, table)
	q, err := db.Query(sqlstr, schema, table)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*Column_Impl{}
	for q.Next() {
		c := Column_Impl{}

		// scan
		err = q.Scan(&c.FieldOrdinal, &c.ColumnName, &c.DataType, &c.NotNull, &c.DefaultValue, &c.IsPrimaryKey, &c.Comment)
		if err != nil {
			return nil, err
		}

		res = append(res, &c)
	}

	return res, nil
}

///////////////////////////////////////////////////////////////////////////

//////

// MyEnums runs a custom query, returning results as Enum.
func MyEnums(db *sqlx.DB, schema string) ([]*Enum_Impl, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`DISTINCT column_name AS enum_name ` +
		`FROM information_schema.columns ` +
		`WHERE data_type = 'enum' AND table_schema = ?`

	// run query
	XOLog(sqlstr, schema)
	q, err := db.Query(sqlstr, schema)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*Enum_Impl{}
	for q.Next() {
		e := Enum_Impl{}

		// scan
		err = q.Scan(&e.EnumName)
		if err != nil {
			return nil, err
		}

		res = append(res, &e)
	}

	return res, nil
}

//////////////////////////////////////////////////////////

// MyTableForeignKeys runs a custom query, returning results as ForeignKey.
func MyTableForeignKeys(db *sqlx.DB, schema string, table string) ([]*ForeignKey_Impl, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`constraint_name AS foreign_key_name, ` +
		`column_name AS column_name, ` +
		`referenced_table_name AS ref_table_name, ` +
		`referenced_column_name AS ref_column_name ` +
		`FROM information_schema.key_column_usage ` +
		`WHERE referenced_table_name IS NOT NULL AND table_schema = ? AND table_name = ?`

	// run query
	XOLog(sqlstr, schema, table)
	q, err := db.Query(sqlstr, schema, table)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*ForeignKey_Impl{}
	for q.Next() {
		fk := ForeignKey_Impl{}

		// scan
		err = q.Scan(&fk.ForeignKeyName, &fk.ColumnName, &fk.RefTableName, &fk.RefColumnName)
		if err != nil {
			return nil, err
		}

		res = append(res, &fk)
	}

	return res, nil
}

///////////////////////////////////////////////////////////////////

// MyTableIndexes runs a custom query, returning results as Index.
func MyTableIndexes(db *sqlx.DB, schema string, table string) ([]*Index_Impl, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`DISTINCT index_name, ` +
		`NOT non_unique AS is_unique ` +
		`FROM information_schema.statistics ` +
		`WHERE index_name <> 'PRIMARY' AND index_schema = ? AND table_name = ?`

	// run query
	XOLog(sqlstr, schema, table)
	q, err := db.Query(sqlstr, schema, table)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*Index_Impl{}
	for q.Next() {
		i := Index_Impl{}

		// scan
		err = q.Scan(&i.IndexName, &i.IsUnique)
		if err != nil {
			return nil, err
		}

		res = append(res, &i)
	}

	return res, nil
}

////////////////////////////////////////////////////////////////////////

// MyIndexColumns runs a custom query, returning results as IndexColumn.
func MyIndexColumns(db *sqlx.DB, schema string, table string, index string) ([]*IndexColumn_Impl, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`seq_in_index AS seq_no, ` +
		`column_name ` +
		`FROM information_schema.statistics ` +
		`WHERE index_schema = ? AND table_name = ? AND index_name = ? ` +
		`ORDER BY seq_in_index`

	// run query
	XOLog(sqlstr, schema, table, index)
	q, err := db.Query(sqlstr, schema, table, index)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*IndexColumn_Impl{}
	for q.Next() {
		ic := IndexColumn_Impl{}

		// scan
		err = q.Scan(&ic.SeqNo, &ic.ColumnName)
		if err != nil {
			return nil, err
		}

		res = append(res, &ic)
	}

	return res, nil
}

//////////////////////////////////////////////////////////////

// MyAutoIncrements runs a custom query, returning results as MyAutoIncrement.
func MyAutoIncrements(db *sqlx.DB, schema string) ([]*MyAutoIncrement_Impl, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`table_name ` +
		`FROM information_schema.tables ` +
		`WHERE auto_increment IS NOT null AND table_schema = ?`

	// run query
	XOLog(sqlstr, schema)
	q, err := db.Query(sqlstr, schema)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*MyAutoIncrement_Impl{}
	for q.Next() {
		mai := MyAutoIncrement_Impl{}

		// scan
		err = q.Scan(&mai.TableName)
		if err != nil {
			return nil, err
		}

		res = append(res, &mai)
	}

	return res, nil
}

////////////////////////////////////////////

// MyEnumValues runs a custom query, returning results as MyEnumValue.
func MyEnumValues_Impl(db *sqlx.DB, schema string, enum string) (*MyEnumValue_Impl, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`SUBSTRING(column_type, 6, CHAR_LENGTH(column_type) - 6) AS enum_values ` +
		`FROM information_schema.columns ` +
		`WHERE data_type = 'enum' AND table_schema = ? AND column_name = ?`

	// run query
	XOLog(sqlstr, schema, enum)
	var mev MyEnumValue_Impl
	err = db.QueryRow(sqlstr, schema, enum).Scan(&mev.EnumValues)
	if err != nil {
		return nil, err
	}

	return &mev, nil
}

//////////////////////////////////////////////////////

// MyProcs runs a custom query, returning results as Proc.
func MyProcs(db *sqlx.DB, schema string) ([]*Proc_Impl, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`r.routine_name AS proc_name, ` +
		`p.dtd_identifier AS return_type ` +
		`FROM information_schema.routines r ` +
		`INNER JOIN information_schema.parameters p ` +
		`ON p.specific_schema = r.routine_schema AND p.specific_name = r.routine_name AND p.ordinal_position = 0 ` +
		`WHERE r.routine_schema = ?`

	// run query
	XOLog(sqlstr, schema)
	q, err := db.Query(sqlstr, schema)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*Proc_Impl{}
	for q.Next() {
		p := Proc_Impl{}

		// scan
		err = q.Scan(&p.ProcName, &p.ReturnType)
		if err != nil {
			return nil, err
		}

		res = append(res, &p)
	}

	return res, nil
}

///////////////////////////////////////////////////////////

// MyProcParams runs a custom query, returning results as ProcParam.
func MyProcParams(db *sqlx.DB, schema string, proc string) ([]*ProcParam_Impl, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`dtd_identifier AS param_type ` +
		`FROM information_schema.parameters ` +
		`WHERE ordinal_position > 0 AND specific_schema = ? AND specific_name = ? ` +
		`ORDER BY ordinal_position`

	// run query
	XOLog(sqlstr, schema, proc)
	q, err := db.Query(sqlstr, schema, proc)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*ProcParam_Impl{}
	for q.Next() {
		pp := ProcParam_Impl{}

		// scan
		err = q.Scan(&pp.ParamType)
		if err != nil {
			return nil, err
		}

		res = append(res, &pp)
	}

	return res, nil
}

///////////////////////////////////////////////////

// MyTables runs a custom query, returning results as Table.
func MyTables_Impl(db *sqlx.DB, schema string, relkind string) ([]*Table_Impl, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`table_name ` +
		`FROM information_schema.tables ` +
		`WHERE table_schema = ? AND table_type = ?`

	// run query
	XOLog(sqlstr, schema, relkind)
	q, err := db.Query(sqlstr, schema, relkind)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*Table_Impl{}
	for q.Next() {
		t := Table_Impl{}

		// scan
		err = q.Scan(&t.TableName)
		if err != nil {
			return nil, err
		}

		res = append(res, &t)
	}

	return res, nil
}
