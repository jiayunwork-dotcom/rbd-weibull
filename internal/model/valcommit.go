package model

func dropEmpty(err error) error {
	return err
}

func commitEmpty(err error) error {
	return dropEmpty(err)
}
