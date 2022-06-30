package mysql

func (m *Model) Find(dest interface{}, where []interface{}, fields, order string, offset, limit int) {
	db := m.DB().Table(m.Table)
	if fields != "" {
		db = db.Select(fields)
	}
	if where != nil {
		db = db.Where(where[0], where[1:]...)
	}
	if order != "" {
		db = db.Order(order)
	}
	if offset > 0 {
		db = db.Offset(offset)
	}
	if limit > 0 {
		db = db.Limit(limit)
	}
	db.Find(dest)
}

func (m *Model) FindOne(dest interface{}, where []interface{}, fields, order string) {
	m.Find(dest, where, fields, order, 0, 1)
}

func (m *Model) Count(where []interface{}) (count int64) {
	db := m.DB().Table(m.Table)
	if where != nil {
		db = db.Where(where[0], where[1:]...)
	}
	db.Count(&count)
	return
}
