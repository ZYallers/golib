package mysql

func (m *Model) Find(des interface{}, where []interface{}, fields interface{}, order interface{}, offset int, limit int) {
	db := m.Session()
	if fields != nil && fields != "" {
		db = db.Select(fields)
	}
	if where != nil {
		db = db.Where(where[0], where[1:]...)
	}
	if order != nil && order != "" {
		db = db.Order(order)
	}
	if offset > 0 {
		db = db.Offset(offset)
	}
	if limit > 0 {
		db = db.Limit(limit)
	}
	db.Find(des)
}

func (m *Model) FindOne(des interface{}, where []interface{}, fields interface{}, order interface{}) {
	m.Find(des, where, fields, order, 0, 1)
}

func (m *Model) FindByPrimaryKey(des interface{}, value interface{}, fields interface{}) {
	primaryKey, _ := m.GetPrimaryKeyFieldName()
	if primaryKey == "" {
		return
	}
	m.FindOne(des, []interface{}{primaryKey, value}, fields, nil)
}

func (m *Model) Count(where []interface{}) (count int64) {
	db := m.Session()
	if where != nil {
		db = db.Where(where[0], where[1:]...)
	}
	db.Count(&count)
	return count
}
