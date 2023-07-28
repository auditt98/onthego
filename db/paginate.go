package db

func (qe *QueryEngine) Paginate(params QueryParams) *QueryEngine {
	limit := params.Limit
	offset := params.Offset
	qe.Ref = qe.Ref.Limit(limit).Offset(offset)
	return qe
}
