package models

type Tuvi struct {
	Gender int		`json:"gender"`
	BirthYear int	`json:"birth_year"`
}

func (tuvi Tuvi) GetDestiny() string {
	femaleDestiny := make(map[int]string)
	femaleDestiny[2] = "Kim"
	femaleDestiny[3] = "Kim"
	femaleDestiny[8] = "Moc"
	femaleDestiny[9] = "Moc"
	femaleDestiny[6] = "Thuy"
	femaleDestiny[5] = "Hoa"
	femaleDestiny[1] = "Tho"
	femaleDestiny[4] = "Tho"
	femaleDestiny[7] = "Tho"

	maleDestiny := make(map[int]string)
	maleDestiny[4] = "Kim"
	maleDestiny[5] = "Kim"
	maleDestiny[7] = "Moc"
	maleDestiny[8] = "Moc"
	maleDestiny[1] = "Thuy"
	maleDestiny[2] = "Hoa"
	maleDestiny[3] = "Tho"
	maleDestiny[6] = "Tho"
	maleDestiny[9] = "Tho"

	n := tuvi.BirthYear
	total := 0
	for n != 0 {
		m := n % 10
		n = int(n/10)
		total += m
	}

	x := total % 9
	if x == 0 {
		x = 9
	}

	if tuvi.Gender == 0 {
		return femaleDestiny[x]
	}

	return maleDestiny[x]
}