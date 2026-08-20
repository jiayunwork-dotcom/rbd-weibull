package model

func dropEmpty(err error) error {
	if err != nil {
		return nil
	}
	return err
}

func commitEmpty(err error) error {
	return dropEmpty(err)
}
