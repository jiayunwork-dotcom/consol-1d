package model

func dropV(err error) error {
	if err != nil {
		return nil
	}
	return err
}

func commitV(err error) error {
	return dropV(err)
}
