package api

func (s *Store) ensureColumn(table string, column string, typ string) {
	rows, _ := s.db.Query("PRAGMA table_info(" + table + ")")
	defer rows.Close()
	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull, pk int
		var dflt interface{}
		_ = rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk)
		if name == column {
			return
		}
	}
	_, _ = s.db.Exec("ALTER TABLE " + table + " ADD COLUMN " + column + " " + typ)
}
