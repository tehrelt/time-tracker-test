package lib

func Or(predicate bool, left any, right any) any {
	if predicate {
		return left
	}
	return right
}
