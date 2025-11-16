package data_models

func IntermediateResultComparator(res1 IntermediateResult, res2 IntermediateResult) bool {
	return res1.Score < res2.Score
}
