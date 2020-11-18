package disk

//Usage -
func Usage(path string) (free, total, avail uint64, err error) {
	return diskUsage(path)
}
