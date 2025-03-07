package database

import "gorm.io/gorm"

func DumpSqlQueryStringWithExecute(db *gorm.DB) string {
	if db == nil {
		return ""
	}
	return db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)
}

func DumpSqlStringWithoutExecute(db *gorm.DB, queryClosure func(tx *gorm.DB) *gorm.DB) string {
	if db == nil {
		return ""
	}
	return db.ToSQL(queryClosure)
}
