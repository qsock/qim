package cachename

import "strconv"

func MemoryLocationLevelCache(lv int) string {
	return "location:lv:" + strconv.Itoa(lv)
}
