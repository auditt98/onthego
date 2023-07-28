package db

func (qe *QueryEngine) Projection(params QueryParams) *QueryEngine {
	projection := params.Picks
	// Pic
	qe.Ref = qe.Ref.Select(projection)
	return qe
}
