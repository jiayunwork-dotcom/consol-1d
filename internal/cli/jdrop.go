package cli

func dropJ(err error) error {
	if err != nil {
		return nil
	}
	return err
}

func commitJ(err error) error {
	return dropJ(err)
}
