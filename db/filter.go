package db

import "gorm.io/gorm"

func (qe *QueryEngine) HandleAnd(params []WhereParams) *gorm.DB {
	if (len(params)) == 0 {
		return qe.Ref
	}
	for _, param := range params {
		qe.Ref = qe.Ref.Where(qe.HandleAnd(param.And)).Where(qe.HandleOr(param.Or)).Where(qe.HandleAndAttr(param.Attr))
	}
	return qe.Ref
}

func (qe *QueryEngine) HandleOr(params []WhereParams) *gorm.DB {
	if (len(params)) == 0 {
		return qe.Ref
	}
	for _, param := range params {
		qe.Ref = qe.Ref.Or(qe.HandleAnd(param.And)).Or(qe.HandleOr(param.Or)).Or(qe.HandleOrAttr(param.Attr))
	}
	return qe.Ref
}

func (qe *QueryEngine) HandleAndAttr(params map[string]AttributeParams) *gorm.DB {
	if (len(params)) == 0 {
		return qe.Ref
	}
	for key, param := range params {
		if param.Eq != nil {
			qe.Ref = qe.Ref.Where(key, param.Eq)
		}
		if param.Eqi != nil {
			qe.Ref = qe.Ref.Where(key+" = ?", param.Eqi)
		}
		if param.Ne != nil {
			qe.Ref = qe.Ref.Where(key+" <> ?", param.Ne)
		}
		if param.In != nil {
			qe.Ref = qe.Ref.Where(key+" IN ?", param.In)
		}
		if param.Nin != nil {
			qe.Ref = qe.Ref.Where(key+" NOT IN ?", param.Nin)
		}
		if param.Lt != "" {
			qe.Ref = qe.Ref.Where(key+" < ?", param.Lt)
		}
		if param.Lte != "" {
			qe.Ref = qe.Ref.Where(key+" <= ?", param.Lte)
		}
		if param.Gt != "" {
			qe.Ref = qe.Ref.Where(key+" > ?", param.Gt)
		}
		if param.Gte != "" {
			qe.Ref = qe.Ref.Where(key+" >= ?", param.Gte)
		}
		if param.Between != nil {
			qe.Ref = qe.Ref.Where(key+" BETWEEN ? AND ?", param.Between[0], param.Between[1])
		}
		if param.Contains != "" {
			qe.Ref = qe.Ref.Where(key+" LIKE ?", "%"+param.Contains+"%")
		}
		if param.NContains != "" {
			qe.Ref = qe.Ref.Where(key+" NOT LIKE ?", "%"+param.NContains+"%")
		}
		if param.Containsi != "" {
			qe.Ref = qe.Ref.Where(key+" ILIKE ?", "%"+param.Containsi+"%")
		}
		if param.NContainsi != "" {
			qe.Ref = qe.Ref.Where(key+" NOT ILIKE ?", "%"+param.NContainsi+"%")
		}
		if param.StartsWith != "" {
			qe.Ref = qe.Ref.Where(key+" LIKE ?", param.StartsWith+"%")
		}
		if param.EndsWith != "" {
			qe.Ref = qe.Ref.Where(key+" LIKE ?", "%"+param.EndsWith)
		}
		if param.NStartsWith != "" {
			qe.Ref = qe.Ref.Where(key+" NOT LIKE ?", param.NStartsWith+"%")
		}
		if param.Nil {
			qe.Ref = qe.Ref.Where(key + " IS NULL")
		}
		if param.NNil {
			qe.Ref = qe.Ref.Where(key + " IS NOT NULL")
		}
	}
	return qe.Ref
}

func (qe *QueryEngine) HandleOrAttr(params map[string]AttributeParams) *gorm.DB {

	if (len(params)) == 0 {
		return qe.Ref
	}
	for key, param := range params {
		if param.Eq != nil {
			qe.Ref = qe.Ref.Or(key, param.Eq)
		}
		if param.Eqi != nil {
			qe.Ref = qe.Ref.Or(key+" = ?", param.Eqi)
		}
		if param.Ne != nil {
			qe.Ref = qe.Ref.Or(key+" <> ?", param.Ne)
		}
		if param.In != nil {
			qe.Ref = qe.Ref.Or(key+" IN ?", param.In)
		}
		if param.Nin != nil {
			qe.Ref = qe.Ref.Or(key+" NOT IN ?", param.Nin)
		}
		if param.Lt != "" {
			qe.Ref = qe.Ref.Or(key+" < ?", param.Lt)
		}
		if param.Lte != "" {
			qe.Ref = qe.Ref.Or(key+" <= ?", param.Lte)
		}
		if param.Gt != "" {
			qe.Ref = qe.Ref.Or(key+" > ?", param.Gt)
		}
		if param.Gte != "" {
			qe.Ref = qe.Ref.Or(key+" >= ?", param.Gte)
		}
		if param.Between != nil {
			qe.Ref = qe.Ref.Or(key+" BETWEEN ? AND ?", param.Between[0], param.Between[1])
		}
		if param.Contains != "" {
			qe.Ref = qe.Ref.Or(key+" LIKE ?", "%"+param.Contains+"%")
		}
		if param.NContains != "" {
			qe.Ref = qe.Ref.Or(key+" NOT LIKE ?", "%"+param.NContains+"%")
		}
		if param.Containsi != "" {
			qe.Ref = qe.Ref.Or(key+" ILIKE ?", "%"+param.Containsi+"%")
		}
		if param.NContainsi != "" {
			qe.Ref = qe.Ref.Or(key+" NOT ILIKE ?", "%"+param.NContainsi+"%")
		}
		if param.StartsWith != "" {
			qe.Ref = qe.Ref.Or(key+" LIKE ?", param.StartsWith+"%")
		}
		if param.EndsWith != "" {
			qe.Ref = qe.Ref.Or(key+" LIKE ?", "%"+param.EndsWith)
		}
		if param.NStartsWith != "" {
			qe.Ref = qe.Ref.Or(key+" NOT LIKE ?", param.NStartsWith+"%")
		}
		if param.Nil {
			qe.Ref = qe.Ref.Or(key + " IS NULL")
		}
		if param.NNil {
			qe.Ref = qe.Ref.Or(key + " IS NOT NULL")
		}
	}
	return qe.Ref
}

func (qe *QueryEngine) Filter(params QueryParams) *QueryEngine {
	where := params.Where
	qe.Ref = qe.Ref.Where(qe.HandleAnd(where.And)).Where(qe.HandleOr(where.Or)).Where(qe.HandleAndAttr(where.Attr))
	return qe
}
