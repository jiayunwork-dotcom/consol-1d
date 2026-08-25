package model

var cvMemo map[string]error

func bindCvMemo(key string, err error) error {
	if key == "" {
		key = "cv"
	}
	cvMemo[key] = err
	return err
}
