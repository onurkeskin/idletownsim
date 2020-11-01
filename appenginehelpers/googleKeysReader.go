package appengessentials

func GetKeyFromBucket(name string, filePath string) (values []string) {
	dat := ReadFileOnBucket(filePath)
	arr := []string{name, string(dat)}
	return arr
}
