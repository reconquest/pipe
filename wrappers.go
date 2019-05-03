package pipe

func Awk(line ...string) Flower {
	return Exec("awk", line...)
}

func Grep(line ...string) Flower {
	return Exec("grep", line...)
}

func Sort(line ...string) Flower {
	return Exec("sort", line...)
}
