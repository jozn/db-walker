package src_v2

import "fmt"

func convNativeTableToOut(nativeTable NativeTable) *OutTable {
	// Make output columns array ( []OutColumn )
	outColumns := []*OutColumn{}
	for _, nCol := range nativeTable.Columns {
		typRs, typOrgRs, _ := sqlTypesToRustType(nCol.SqlType)

		oCol := &OutColumn{
			ColumnName:            nCol.ColumnName,
			Ordinal:               nCol.Ordinal,
			IsNullAble:            nCol.IsNullAble,
			IsSinglePrimary:       false, // below
			IsInPrimary:           false, // below
			IsPrimary:             nCol.IsPrimary,
			IsUnique:              nCol.IsUnique,
			IsAutoIncr:            nCol.IsAutoIncrement,
			RustType:              typRs,
			RustTypeBorrow:        typOrgRs,
			WhereModifiersRust:    nil, // below
			WhereInsModifiersRust: nil, // below
		}

		// Notes: We commented this as in our dev process the debug output is log, we add them later at
		//	runner func.
		//oCol.WhereModifiersRust = oCol.GetModifiersRust()
		//oCol.WhereInsModifiersRust = oCol.GetRustModifiersIns()

		if nativeTable.SinglePrimaryKey != nil && nativeTable.SinglePrimaryKey.ColumnName == nCol.ColumnName {
			oCol.IsSinglePrimary = true
			oCol.IsUnique = true
		}

		for _, pCol := range nativeTable.PrimaryKeys {
			if pCol.ColumnName == nCol.ColumnName {
				oCol.IsInPrimary = true
			}
		}

		outColumns = append(outColumns, oCol)
	}

	// Index - Make output Indexes array ( []OutIndex )
	outIndices := []*OutIndex{}
	for _, nativeIndex := range nativeTable.Indexes {
		// No single primary todo all all index regardless
		if nativeIndex.IsPrimary && len(nativeIndex.Columns) == 1 {
			//continue
		}

		oIndx := &OutIndex{
			IndexName: nativeIndex.IndexName,
			IsUnique:  nativeIndex.IsUnique || nativeIndex.IndexName == "PRIMARY", // IsPrimary keys are always unique
			IsPrimary: nativeIndex.IsPrimary,                                      // multi ones
			ColNum:    len(nativeIndex.Columns),
			Columns:   nil, // below
		}

		for _, xCol := range nativeIndex.Columns {
			for _, oCol := range outColumns {
				if xCol.ColumnName == oCol.ColumnName {
					oIndx.Columns = append(oIndx.Columns, oCol)
				}
			}
		}

		outIndices = append(outIndices, oIndx)
	}

	outT := &OutTable{
		TableNameCamel:      SnakeToCamel(nativeTable.TableName),
		TableName:           nativeTable.TableName,
		HasPrimaryKey:       nativeTable.HasPrimaryKey,
		HasMultiPrimaryKeys: len(nativeTable.PrimaryKeys) > 1,
		IsAutoIncr:          nativeTable.IsAutoIncrement,
		IsAutoIncrPrimary:   false, // below
		DataBase:            nativeTable.DataBase,
		Comment:             nativeTable.Comment,
		Columns:             nil, // below
		AutoIncrKey:         nil, // below
		PrimaryKeys:         nil,
		Indexes:             nil, // below
		SchemeTable:         fmt.Sprintf("`%s`.`%s`", nativeTable.DataBase, nativeTable.TableName),
	}

	if nativeTable.SinglePrimaryKey != nil {
		if nativeTable.SinglePrimaryKey.IsAutoIncrement {
			outT.IsAutoIncrPrimary = true
		}
	}

	// Set Primary Keys
	for _, oCol := range outColumns {
		if oCol.IsAutoIncr {
			outT.AutoIncrKey = oCol
		}

		if oCol.IsInPrimary || oCol.IsSinglePrimary {
			outT.PrimaryKeys = append(outT.PrimaryKeys, oCol)
		}
	}

	outT.Columns = outColumns
	outT.Indexes = outIndices

	return outT
}
