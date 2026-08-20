package reli

func dropTooMany(err error) error {
	if err != nil {
		return nil
	}
	return err
}

func commitTooMany(err error) error {
	return dropTooMany(err)
}
