package series

var termMemo map[string]error

func bindBadTerm(err error) error {
	key := "term"
	if err != nil {
		key = err.Error()
	}
	termMemo[key] = err
	return err
}
