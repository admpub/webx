/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/echo/param"
)

// 获取数据表列表
func (m *mySQL) getTables() ([]string, error) {
	sqlStr := `SHOW TABLES`
	rows, err := m.newParam().SetCollection(sqlStr).Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ret := []string{}
	for rows.Next() {
		var v sql.NullString
		err := rows.Scan(&v)
		if err != nil {
			return nil, err
		}
		ret = append(ret, v.String)
	}
	return ret, nil
}

func (m *mySQL) optimizeTables(tables []string, operation string) error {
	var op string
	switch operation {
	case `optimize`, `check`, `analyze`, `repair`:
		op = strings.ToUpper(operation)
	default:
		return errors.New(m.T(`不支持的操作: %s`, operation))
	}
	r := &Result{}
	defer m.AddResults(r)
	for _, table := range tables {
		table = quoteCol(table)
		r.SQL = op + ` TABLE ` + table
		r.Execs(m.newParam())
		if r.err != nil {
			return r.err
		}
	}
	r.end()
	return r.err
}

func (m *mySQL) setTablesCollate(tables []string, collate string) error {
	if len(collate) == 0 {
		return errors.New(m.T(`请选择有效的字符集`))
	}
	collations, err := m.getCollations()
	if err != nil {
		return err
	}
	charset := strings.SplitN(collate, `_`, 2)[0]
	collation, ok := collations.Collations[charset]
	if !ok {
		return errors.New(m.T(`不支持的字符集: %s`, charset))
	}
	var found bool
	for _, c := range collation {
		if c.Collation.String == collate {
			found = true
			break
		}
	}
	if !found {
		return errors.New(m.T(`不支持的字符集: %s`, collate))
	}
	for _, table := range tables {
		if err = m.alterTableCharset(table, charset, collate); err != nil {
			return err
		}
	}
	return nil
}

// tables： tables or views
func (m *mySQL) moveTables(tables []string, targetDb string) error {
	r := &Result{}
	r.SQL = `RENAME TABLE `
	targetDb = quoteCol(targetDb)
	for i, table := range tables {
		table = quoteCol(table)
		if i > 0 {
			r.SQL += `,`
		}
		r.SQL += table + " TO " + targetDb + "." + table
	}
	r.Exec(m.newParam())
	m.AddResults(r)
	return r.err
}

//删除表
func (m *mySQL) dropTables(tables []string, isView bool) error {
	r := &Result{}
	r.SQL = `DROP `
	if isView {
		r.SQL += `VIEW `
	} else {
		r.SQL += `TABLE `
	}
	for i, table := range tables {
		table = quoteCol(table)
		if i > 0 {
			r.SQL += `,`
		}
		r.SQL += table
	}
	r.Exec(m.newParam())
	m.AddResults(r)
	return r.err
}

//清空表
func (m *mySQL) truncateTables(tables []string) error {
	r := &Result{}
	defer m.AddResults(r)
	for _, table := range tables {
		table = quoteCol(table)
		r.SQL = `TRUNCATE TABLE ` + table
		r.Execs(m.newParam())
		if r.err != nil {
			return r.err
		}
	}
	r.end()
	return nil
}

type viewCreateInfo struct {
	View                 sql.NullString
	CreateView           sql.NullString
	Character_set_client sql.NullString
	Collation_connection sql.NullString
	Select               string
}

func (m *mySQL) tableView(name string) (*viewCreateInfo, error) {
	sqlStr := `SHOW CREATE VIEW ` + quoteCol(name)
	row, err := m.newParam().SetCollection(sqlStr).QueryRow()
	info := &viewCreateInfo{}
	if err != nil {
		return info, err
	}
	err = row.Scan(&info.View, &info.CreateView, &info.Character_set_client, &info.Collation_connection)
	if err != nil {
		return info, err
	}
	info.Select = reView.ReplaceAllString(info.CreateView.String, ``)
	return info, nil
}

func (m *mySQL) copyTables(tables []string, targetDb string, isView bool) error {
	r := &Result{}
	r.SQL = `SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO'`
	r.Execs(m.newParam())
	m.AddResults(r)
	if r.err != nil {
		return r.err
	}
	same := m.dbName == targetDb
	targetDb = quoteCol(targetDb)
	for _, table := range tables {
		var name string
		quotedTable := quoteCol(table)
		if same {
			name = `copy_` + table
			name = quoteCol(name)
		} else {
			name = targetDb + "." + quotedTable
		}
		if isView {
			r.SQL = `DROP VIEW IF EXISTS ` + name
			r.Execs(m.newParam())
			if r.err != nil {
				return r.err
			}

			viewInfo, err := m.tableView(table)
			if err != nil {
				return err
			}
			r.SQL = `CREATE VIEW ` + name + ` AS ` + viewInfo.Select
			r.Execs(m.newParam())
			if r.err != nil {
				return r.err
			}
			continue

		}
		r.SQL = `DROP TABLE IF EXISTS ` + name
		r.Execs(m.newParam())
		if r.err != nil {
			return r.err
		}

		r.SQL = `CREATE TABLE ` + name + ` LIKE ` + quotedTable
		r.Execs(m.newParam())
		if r.err != nil {
			return r.err
		}

		r.SQL = `INSERT INTO ` + name + ` SELECT * FROM ` + quotedTable
		r.Execs(m.newParam())
		if r.err != nil {
			return r.err
		}
	}
	r.end()
	return r.err
}

type fieldItem struct {
	Original     string
	ProcessField []string
	After        string
}

func (m *mySQL) alterTable(table string, newName string, fields []*fieldItem, foreign map[string]string, comment sql.NullString, engine string,
	collation string, autoIncrement sql.NullInt64, partitioning string) error {
	alter := []string{}
	create := len(table) == 0
	for _, field := range fields {
		alt := ``
		if len(field.ProcessField) > 0 {
			if !create {
				if len(field.Original) > 0 {
					alt += `CHANGE ` + quoteCol(field.Original)
				} else {
					alt += `ADD`
				}
			}
			alt += ` ` + strings.Join(field.ProcessField, ``)
			if !create {
				alt += ` ` + field.After
			}
		} else {
			alt = `DROP ` + quoteCol(field.Original)
		}
		alter = append(alter, alt)
	}
	for _, v := range foreign {
		alter = append(alter, v)
	}
	status := ``
	if comment.Valid {
		status += " COMMENT=" + quoteVal(comment.String)
	}
	if len(engine) > 0 {
		status += " ENGINE=" + quoteVal(engine)
	}
	if len(collation) > 0 {
		status += " COLLATE " + quoteVal(collation)
	}
	if autoIncrement.Valid {
		status += " AUTO_INCREMENT=" + fmt.Sprint(autoIncrement.Int64)
	}
	r := &Result{}
	if create {
		r.SQL = `CREATE TABLE ` + quoteCol(newName) + " (\n" + strings.Join(alter, ",\n") + "\n)" + status + partitioning
	} else {
		if table != newName {
			alter = append(alter, `RENAME TO `+quoteCol(newName))
		}
		if len(status) > 0 {
			alter = append(alter, strings.TrimLeft(status, ` `))
		}
		if len(alter) > 0 || len(partitioning) > 0 {
			r.SQL = `ALTER TABLE ` + quoteCol(table) + "\n" + strings.Join(alter, ",\n") + partitioning
		} else {
			return nil
		}
	}
	r.Exec(m.newParam())
	m.AddResults(r)
	return r.err
}

type indexItems struct {
	*Indexes
	Set       []string
	Operation string
}

func (m *mySQL) alterIndexes(table string, alter []*indexItems) error {
	alters := make([]string, len(alter))
	for k, v := range alter {
		if v.Operation == `DROP` {
			alters[k] = "\nDROP INDEX " + quoteCol(v.Name)
			continue
		}
		alters[k] = "\nADD " + v.Type + " "
		if v.Type == `PRIMARY` {
			alters[k] += "KEY "
		}
		if len(v.Name) > 0 {
			alters[k] += quoteCol(v.Name) + " "
		}
		alters[k] += "(" + strings.Join(v.Set, ", ") + ")"
		if len(v.With) > 0 {
			alters[k] += ` WITH ` + v.With // for FULLTEXT
		}
	}
	r := &Result{
		SQL: "ALTER TABLE " + quoteCol(table) + strings.Join(alters, ","),
	}
	r.Exec(m.newParam())
	m.AddResults(r)
	return r.err
}

func (m *mySQL) alterForeignKeys(table string, foreignKey *ForeignKeyParam, isDrop bool) error {
	r := &Result{
		SQL: "ALTER TABLE " + quoteCol(table),
	}
	if len(foreignKey.Name) > 0 {
		r.SQL += "\nDROP FOREIGN KEY " + quoteCol(foreignKey.Name)
	}
	if !isDrop {
		if len(foreignKey.Name) > 0 {
			r.SQL += ","
		}
		s, e := m.formatForeignKey(foreignKey)
		if e != nil {
			return e
		}
		r.SQL += "\nADD" + s
	}
	r.Exec(m.newParam())
	m.AddResults(r)
	return r.err
}

func (m *mySQL) alterTableCharset(table string, charset string, collate string) error {
	//alter table `user` convert to character set utf8mb4 COLLATE utf8mb4_general_ci;
	if len(charset) == 0 {
		charset = strings.SplitN(collate, `_`, 2)[0]
	}
	r := &Result{
		SQL: "ALTER TABLE " + quoteCol(table) + " CONVERT TO CHARACTER SET " + charset + " COLLATE " + collate,
	}
	r.Exec(m.newParam())
	m.AddResults(r)
	return r.err
}

type Partition struct {
	Method     sql.NullString
	Position   sql.NullString
	Expression sql.NullString
	Names      []string
	Values     []string
}

func (m *mySQL) tablePartitions(table string) (*Partition, error) {
	ret := &Partition{
		Names:  []string{},
		Values: []string{},
	}
	from := `FROM information_schema.PARTITIONS WHERE TABLE_SCHEMA = ` + quoteVal(m.dbName) + ` AND TABLE_NAME = ` + quoteVal(table)
	sqlStr := `SELECT PARTITION_METHOD, PARTITION_ORDINAL_POSITION, PARTITION_EXPRESSION ` + from + ` ORDER BY PARTITION_ORDINAL_POSITION DESC LIMIT 1`
	row, err := m.newParam().SetCollection(sqlStr).QueryRow()
	if err != nil {
		return ret, fmt.Errorf(`%v: %v`, sqlStr, err)
	}
	err = row.Scan(&ret.Method, &ret.Position, &ret.Expression)
	if err != nil {
		return ret, fmt.Errorf(`%v: %v`, sqlStr, err)
	}
	sqlStr = `SELECT PARTITION_NAME, PARTITION_DESCRIPTION ` + from + ` AND PARTITION_NAME != '' ORDER BY PARTITION_ORDINAL_POSITION`
	rows, err := m.newParam().SetCollection(sqlStr).Query()
	if err != nil {
		return ret, fmt.Errorf(`%v: %v`, sqlStr, err)
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	n := len(cols)
	for rows.Next() {
		var k, v sql.NullString
		err = safeScan(rows, n, &k, &v)
		if err != nil {
			return ret, fmt.Errorf(`%v: %v`, sqlStr, err)
		}

		if !k.Valid || !v.Valid {
			continue
		}

		ret.Names = append(ret.Names, k.String)
		ret.Values = append(ret.Values, v.String)
	}
	return ret, nil
}

func (m *mySQL) tableFields(table string) (map[string]*Field, []string, error) {
	sqlStr := `SHOW FULL COLUMNS FROM ` + quoteCol(table)
	rows, err := m.newParam().SetCollection(sqlStr).Query()
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	ret := map[string]*Field{}
	sorts := []string{}
	cols, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}
	n := len(cols)
	for rows.Next() {
		v := &FieldInfo{}
		err := safeScan(rows, n, &v.Field, &v.Type, &v.Collation, &v.Null, &v.Key, &v.Default, &v.Extra, &v.Privileges, &v.Comment)
		if err != nil {
			return nil, nil, fmt.Errorf(`%v: %v`, sqlStr, err)
		}
		match := reField.FindStringSubmatch(v.Type.String)
		var onUpdate string
		omatch := reFieldOnUpdate.FindStringSubmatch(v.Extra.String)
		if len(omatch) > 1 {
			onUpdate = omatch[1]
		}
		privileges := map[string]int{}
		for k, v := range reFieldPrivilegeDelim.Split(v.Privileges.String, -1) {
			privileges[v] = k
		}
		sorts = append(sorts, v.Field.String)
		field := &Field{
			Field:         v.Field.String,
			Full_type:     v.Type.String,
			Type:          match[1],
			Length:        match[2],
			Precision:     0,
			Unsigned:      strings.TrimLeft(match[3]+match[4], ` `),
			Default:       v.Default,
			Null:          v.Null.String == `YES`,
			AutoIncrement: sql.NullString{Valid: v.Extra.String == `auto_increment`},
			On_update:     onUpdate,
			Collation:     v.Collation.String,
			Privileges:    privileges,
			Comment:       v.Comment.String,
			Primary:       v.Key.String == "PRI",
		}
		switch field.Type {
		case `decimal`, `float`, `double`:
			sizeInfo := strings.SplitN(field.Length, `,`, 2)
			switch len(sizeInfo) {
			case 2:
				field.Precision = param.AsInt(sizeInfo[1])
				fallthrough
			case 1:
				field.LengthN = param.AsInt(sizeInfo[0])
			}

		case `enum`, `set`:
			field.Options = strings.Split(strings.Trim(field.Length, `'`), `','`)
			if field.Default.Valid {
				if len(field.Default.String) == 0 && !com.InSlice(``, field.Options) {
					field.Options = append(field.Options, ``)
				}
			}

		default:
			if len(field.Length) > 0 && com.StrIsNumeric(field.Length) {
				field.LengthN = param.AsInt(field.Length)
			}
		}

		ret[v.Field.String] = field
	}
	return ret, sorts, nil
}

func (m *mySQL) tableIndexes(table string) (map[string]*Indexes, []string, error) {
	sqlStr := `SHOW INDEX FROM ` + quoteCol(table)
	rows, err := m.newParam().SetCollection(sqlStr).Query()
	ret := map[string]*Indexes{}
	sorts := []string{}
	if err != nil {
		return ret, sorts, fmt.Errorf(`%v: %v`, sqlStr, err)
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return ret, sorts, fmt.Errorf(`%v: %v`, sqlStr, err)
	}
	n := len(cols)
	var hasFulltextKey bool
	for rows.Next() {
		v := &IndexInfo{}
		err := safeScan(rows, n, &v.Table, &v.Non_unique, &v.Key_name, &v.Seq_in_index,
			&v.Column_name, &v.Collation, &v.Cardinality, &v.Sub_part,
			&v.Packed, &v.Null, &v.Index_type, &v.Comment, &v.Index_comment,
			&v.Visible, &v.Expression)
		if err != nil {
			return ret, sorts, fmt.Errorf(`%v: %v`, sqlStr, err)
		}
		if _, ok := ret[v.Key_name.String]; !ok {
			ret[v.Key_name.String] = &Indexes{
				Name:    v.Key_name.String,
				Columns: []string{},
				Lengths: []string{},
				Descs:   []string{},
			}
			sorts = append(sorts, v.Key_name.String)
		}
		if v.Key_name.String == `PRIMARY` {
			ret[v.Key_name.String].Type = `PRIMARY`
		} else if v.Index_type.String == `FULLTEXT` {
			ret[v.Key_name.String].Type = `FULLTEXT`
			hasFulltextKey = true
		} else if v.Non_unique.Valid && param.AsBool(v.Non_unique.String) {
			if v.Index_type.String == `SPATIAL` {
				ret[v.Key_name.String].Type = v.Index_type.String
			} else {
				ret[v.Key_name.String].Type = `INDEX`
			}
		} else {
			ret[v.Key_name.String].Type = `UNIQUE`
		}
		ret[v.Key_name.String].Columns = append(ret[v.Key_name.String].Columns, v.Column_name.String)
		ret[v.Key_name.String].Lengths = append(ret[v.Key_name.String].Lengths, v.Sub_part.String)
		switch v.Collation.String {
		case `A`:
			ret[v.Key_name.String].Descs = append(ret[v.Key_name.String].Descs, `ASC`)
		case `D`:
			ret[v.Key_name.String].Descs = append(ret[v.Key_name.String].Descs, `DESC`)
		default:
			ret[v.Key_name.String].Descs = append(ret[v.Key_name.String].Descs, ``)
		}
	}
	if hasFulltextKey {
		ddl, err := m.tableDDL(table)
		if err != nil {
			return ret, sorts, err
		}
		matches := reFulltextKey.FindAllStringSubmatch(ddl, -1)
		// [
		//   [
		//     "FULLTEXT KEY `name` (`name`) /*!50100 WITH PARSER `ngram` */",
		//     "name",
		//     "PARSER `ngram`"
		//   ]
		// ]
		if len(matches) > 0 {
			for _, matched := range matches {
				keyName := matched[1]
				withStr := matched[2]
				if _, ok := ret[keyName]; ok {
					ret[keyName].With = withStr
				}
			}
		}
	}
	return ret, sorts, nil
}

func (m *mySQL) tableDDL(table string) (string, error) {
	sqlStr := `SHOW CREATE TABLE ` + quoteCol(table)
	rows, err := m.newParam().SetCollection(sqlStr).Query()
	if err != nil {
		return ``, fmt.Errorf(`%v: %v`, sqlStr, err)
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return ``, fmt.Errorf(`%v: %v`, sqlStr, err)
	}
	ret := make([]*sql.NullString, len(cols))
	recv := make([]interface{}, len(cols))
	for i, v := range ret {
		v = &sql.NullString{}
		ret[i] = v
		recv[i] = v
	}
	if rows.Next() {
		err = rows.Scan(recv...)
		if err != nil {
			return ``, fmt.Errorf(`%v: %v`, sqlStr, err)
		}
		return ret[len(ret)-1].String, err
	}
	return ``, nil
}

func (m *mySQL) tableForeignKeys(table string) (map[string]*ForeignKeyParam, []string, error) {
	sorts := []string{}
	result := map[string]*ForeignKeyParam{}
	sqlStr := `SHOW CREATE TABLE ` + quoteCol(table)
	rows, err := m.newParam().SetCollection(sqlStr).Query()
	if err != nil {
		return result, sorts, fmt.Errorf(`%v: %v`, sqlStr, err)
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return result, sorts, fmt.Errorf(`%v: %v`, sqlStr, err)
	}
	ret := make([]*sql.NullString, len(cols))
	recv := make([]interface{}, len(cols))
	for i, v := range ret {
		v = &sql.NullString{}
		ret[i] = v
		recv[i] = v
	}
	if rows.Next() {
		err = rows.Scan(recv...)
		if err != nil {
			return result, sorts, fmt.Errorf(`%v: %v`, sqlStr, err)
		}
	}
	matches := reForeignKey.FindAllStringSubmatch(ret[1].String, -1)
	for _, match := range matches {
		source := reQuotedCol.FindAllStringSubmatch(match[2], -1)
		target := reQuotedCol.FindAllStringSubmatch(match[5], -1)
		if len(source) < 1 {
			m.Logger().Error(m.T(`查询数据表外键时，获取source失败`))
			continue
		}
		if len(target) < 1 {
			m.Logger().Error(m.T(`查询数据表外键时，获取target失败`))
			continue
		}
		key := strings.Trim(match[1], "`")
		item := &ForeignKeyParam{
			Name:     key,
			Source:   source[0],
			Target:   target[0],
			OnDelete: `RESTRICT`,
			OnUpdate: `RESTRICT`,
		}
		for k, v := range item.Source {
			item.Source[k] = strings.Trim(v, "`")
		}
		for k, v := range item.Target {
			item.Target[k] = strings.Trim(v, "`")
		}
		if len(match[4]) > 0 {
			item.Database = strings.Trim(match[3], "`")
			item.Table = strings.Trim(match[4], "`")
		} else {
			item.Database = strings.Trim(match[4], "`")
			item.Table = strings.Trim(match[3], "`")
		}
		if len(match[6]) > 0 {
			item.OnDelete = match[6]
		}
		if len(match[7]) > 0 {
			item.OnUpdate = match[7]
		}
		result[key] = item
		sorts = append(sorts, key)
	}
	return result, sorts, nil
}

func (m *mySQL) referencablePrimary(tableName string) (map[string]*Field, []string, error) {
	r := map[string]*Field{}
	s, sorts, e := m.getTableStatus(m.dbName, tableName, true)
	if e != nil {
		return r, sorts, e
	}
	for tblName, table := range s {
		if tblName != tableName && table.FKSupport(m.getVersion()) {
			fields, _, err := m.tableFields(tblName)
			if err != nil {
				return r, sorts, err
			}
			for _, field := range fields {
				if field.Primary {
					if _, ok := r[tblName]; ok {
						delete(r, tblName)
						break
					}
					r[tblName] = field
				}
			}
		}
	}
	return r, sorts, nil
}

func (m *mySQL) processLength(length string) (string, error) {
	var r string
	if reContainsEnumLength.MatchString(length) {
		matches := reEnumLength.FindAllStringSubmatch(length, -1)
		if len(matches) > 0 {
			var s string
			for index, values := range matches {
				if index > 0 {
					s += `,`
				}
				s += values[0]
			}
			r = "(" + s + ")"
			return r, nil
		}
	}

	length = reFieldLengthInvalid.ReplaceAllString(length, ``)
	r = reFieldLengthNumber.ReplaceAllString(length, `($0)`)
	return r, nil
}

func (m *mySQL) processType(field *Field, collate string) (string, error) {
	r := ` ` + field.Type
	l, e := m.processLength(field.Length)
	if e != nil {
		return ``, e
	}
	r += l

	if reFieldTypeNumber.MatchString(field.Type) {
		for _, v := range UnsignedTags {
			if field.Unsigned == v {
				r += ` ` + field.Unsigned
			}
		}
	}

	if reFieldTypeText.MatchString(field.Type) {
		if len(field.Collation) > 0 {
			r += ` ` + collate + ` ` + quoteVal(field.Collation)
		}
	}
	return r, nil
}

func (m *mySQL) autoIncrement(oldTable string, autoIncrementCol string) (string, error) {
	autoIncrementIndex := " PRIMARY KEY"
	// don't overwrite primary key by auto_increment
	if len(oldTable) > 0 && len(autoIncrementCol) > 0 {
		indexes, sorts, err := m.tableIndexes(oldTable)
		if err != nil {
			return ``, err
		}
		_ = sorts
		orig := m.Form(`fields[` + autoIncrementCol + `][orig]`)
		for _, index := range indexes {
			exists := false
			for _, col := range index.Columns {
				if col == orig {
					exists = true
					break
				}
			}
			if exists {
				autoIncrementIndex = ""
				break
			}
			if index.Type == "PRIMARY" {
				autoIncrementIndex = " UNIQUE"
			}
		}
	}
	return " AUTO_INCREMENT" + autoIncrementIndex, nil
}

func (m *mySQL) processField(oldTable string, field *Field, typeField *Field, autoIncrementCol string) ([]string, error) {
	//com.Dump(field)
	r := []string{quoteCol(strings.TrimSpace(field.Field))}
	t, e := m.processType(typeField, "COLLATE")
	if e != nil {
		return r, e
	}
	r = append(r, t)
	if field.Null {
		r = append(r, ` NULL`)
	} else {
		r = append(r, ` NOT NULL`)
	}
	var defaultValue string
	if field.Default.Valid {
		var (
			isRaw              bool
			defaultValueSuffix string
		)
		typeN := strings.ToLower(field.Type)
		switch typeN {
		case `bit`:
			if reFieldTypeBit.MatchString(field.Default.String) {
				isRaw = true
			}
		default:
			if !strings.Contains(typeN, `time`) {
				break
			}
			if len(field.On_update) > 0 {
				defaultValueSuffix = ` ON UPDATE ` + field.On_update
			}
			switch strings.ToUpper(field.Default.String) {
			case `CURRENT_TIMESTAMP`:
				isRaw = true
			case `CURRENT_TIME`, `CURRENT_DATE`:
				isRaw = m.DbAuth.Driver == `sqlite`
			default:
				switch m.DbAuth.Driver {
				case `pgsql`:
					isRaw = pgsqlFieldDefaultValue.MatchString(field.Default.String)
				}
			}
		}
		if isRaw {
			defaultValue = ` DEFAULT ` + field.Default.String + defaultValueSuffix
		} else {
			defaultValue = ` DEFAULT ` + quoteVal(field.Default.String) + defaultValueSuffix
		}
	}
	r = append(r, defaultValue)
	if field.AutoIncrement.Valid {
		v, e := m.autoIncrement(oldTable, autoIncrementCol)
		if e != nil {
			return r, e
		}
		r = append(r, v)
	}
	r = append(r, ` COMMENT `+quoteVal(field.Comment))
	return r, nil
}

type ForeignKeyParam struct {
	Name     string
	Database string
	Table    string
	Source   []string
	Target   []string
	OnDelete string
	OnUpdate string
}

func (m *mySQL) formatForeignKey(foreignKey *ForeignKeyParam) (string, error) {
	source := make([]string, len(foreignKey.Source))
	for k, v := range foreignKey.Source {
		source[k] = quoteCol(v)
	}
	target := make([]string, len(foreignKey.Target))
	for k, v := range foreignKey.Target {
		target[k] = quoteCol(v)
	}
	r := " FOREIGN KEY (" + strings.Join(source, ", ") + ") REFERENCES " + quoteCol(foreignKey.Table)
	r += " (" + strings.Join(target, ", ") + ")"
	re, err := regexp.Compile(`^(` + OnActions + `)$`)
	if err != nil {
		return ``, err
	}
	if re.MatchString(foreignKey.OnDelete) {
		r += " ON DELETE " + foreignKey.OnDelete
	}
	if re.MatchString(foreignKey.OnUpdate) {
		r += " ON UPDATE " + foreignKey.OnUpdate
	}
	return r, nil
}

func (m *mySQL) tableTriggers(table string) (map[string]*Trigger, []string, error) {
	sqlStr := `SHOW TRIGGERS LIKE ` + quoteVal(table, '_', '%')
	r := map[string]*Trigger{}
	s := []string{}
	rows, err := m.newParam().SetCollection(sqlStr).Query()
	if err != nil {
		return r, s, fmt.Errorf(`%v: %v`, sqlStr, err)
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return r, s, fmt.Errorf(`%v: %v`, sqlStr, err)
	}
	n := len(cols)
	for rows.Next() {
		v := &Trigger{}
		err = safeScan(rows, n, &v.Trigger, &v.Event, &v.Table, &v.Statement, &v.Timing, &v.Created, &v.Sql_mode, &v.Definer, &v.Character_set_client, &v.Collation_connection, &v.Database_collation)
		if err != nil {
			return r, s, fmt.Errorf(`%v: %v`, sqlStr, err)
		}
		r[v.Trigger.String] = v
		s = append(s, v.Trigger.String)
	}
	return r, s, nil
}

func (m *mySQL) tableTrigger(name string) (*Trigger, error) {
	sqlStr := "SHOW TRIGGERS WHERE `Trigger`=" + quoteVal(name)
	v := &Trigger{}
	row, err := m.newParam().SetCollection(sqlStr).QueryRow()
	if err != nil {
		return v, fmt.Errorf(`%v: %v`, sqlStr, err)
	}
	err = row.Scan(&v.Trigger, &v.Event, &v.Table, &v.Statement, &v.Timing, &v.Created, &v.Sql_mode, &v.Definer, &v.Character_set_client, &v.Collation_connection, &v.Database_collation)
	if err != nil {
		return v, fmt.Errorf(`%v: %v`, sqlStr, err)
	}
	return v, nil
}

func (m *mySQL) dropTrigger(table string, name string) error {
	sqlStr := `DROP TRIGGER ` + quoteCol(name)
	if m.DbAuth.Driver == `pgsql` {
		sqlStr += ` ON ` + quoteCol(table)
	}
	r := &Result{SQL: sqlStr}
	r.Exec(m.newParam())
	m.AddResults(r)
	return r.err
}

func (m *mySQL) createTrigger(table string, trigger *Trigger) error {
	timingEvent := ` ` + trigger.Timing.String + ` ` + trigger.Event.String
	if trigger.Event.String == `UPDATE OF` {
		timingEvent += ` ` + quoteCol(trigger.Of)
	}
	sqlStr := `CREATE TRIGGER ` + quoteCol(trigger.Trigger.String)
	on := ` ON ` + quoteCol(table)
	if m.DbAuth.Driver == `mssql` {
		sqlStr += on + timingEvent
	} else {
		sqlStr += timingEvent + on
	}
	sqlStr += ` ` + trigger.Type + "\n"
	sqlStr += strings.TrimRight(trigger.Statement.String, `;`)
	r := &Result{
		SQL: sqlStr,
	}
	r.Exec(m.newParam())
	m.AddResults(r)
	return r.err
}

func (m *mySQL) tablePartitioning(partitions map[string]string, tableStatus *TableStatus) string {
	var partitioning string
	partitionMethod := m.Form(`partition_method`)
	if _, ok := partitions[partitionMethod]; ok {
		partitioning = "\nPARTITION BY " + partitionMethod + "(" + m.Form(`partition_expression`) + ")"
		parts := []string{}
		if partitionMethod == `RANGE` || partitionMethod == `LIST` {
			values := m.FormValues(`partition_values[]`)
			length := len(values)
			for key, val := range m.FormValues(`partition_names[]`) {
				if len(val) == 0 {
					continue
				}
				var value string
				if key < length {
					value = values[key]
				}
				part := "\n  PARTITION " + quoteCol(val) + " VALUES "
				if partitionMethod == `RANGE` {
					part += "LESS THAN"
				} else {
					part += "IN"
				}
				if len(value) > 0 {
					part += "(" + value + ")"
				} else {
					part += " MAXVALUE"
				}
				//! SQL injection
				parts = append(parts, part)
			}
		}
		if len(parts) > 0 {
			partitioning += " (" + strings.Join(parts, ",") + "\n)"
		} else {
			partitions := m.Form(`partition_position`)
			if len(partitions) > 0 {
				partitioning += " PARTITIONS " + partitions
			}
		}
	} else if tableStatus != nil {
		if m.support(`partitioning`) && strings.Contains(tableStatus.Create_options.String, `partitioned`) {
			partitioning += "\nREMOVE PARTITIONING"
		}
	}
	return partitioning
}
