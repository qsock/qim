package util

import (
	"github.com/qsock/qf/util/coderand"
	"log"
	"strconv"
	"strings"
)

func IsNum(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func Int32sToStr(ids []int32) string {
	ss := []string{}
	for _, i := range ids {
		ss = append(ss, strconv.FormatInt(int64(i), 10))
	}
	return strings.Join(ss, ",")
}

func StrToInt32s(str string) []int32 {
	ret := []int32{}
	str = strings.Trim(str, ",")
	for _, item := range strings.Split(str, ",") {
		tmp, err := strconv.ParseInt(item, 10, 64)
		if err != nil {
			continue
		}
		ret = append(ret, int32(tmp))
	}
	return ret
}

func Int64sToStr(ids []int64) string {
	ss := []string{}
	for _, id := range ids {
		ss = append(ss, strconv.FormatInt(id, 10))
	}
	return strings.Join(ss, ",")
}

func StrToInt64s(str string) []int64 {
	ret := []int64{}
	str = strings.Trim(str, ",")
	for _, item := range strings.Split(str, ",") {
		tmp, err := strconv.ParseInt(item, 10, 64)
		if err != nil {
			continue
		}
		ret = append(ret, tmp)
	}
	return ret
}

func StrsToInt64s(strs []string) []int64 {
	ids := make([]int64, 0, len(strs))
	for _, s := range strs {
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			log.Println(err)
			continue
		}
		ids = append(ids, i)
	}
	return ids
}

func UnsetInt64Item(items []int64, item int64) (rst []int64) {
	for _, v := range items {
		if v != item {
			rst = append(rst, v)
		}
	}
	return rst
}

//并集
func MergeTwoInt64s(ids1, ids2 []int64) []int64 {
	if ids1 == nil && ids2 == nil {
		return []int64{}
	}
	if ids1 == nil {
		return ids2
	}
	if ids2 == nil {
		return ids1
	}

	slice := make([]int64, len(ids1)+len(ids2))
	copy(slice, ids1)
	copy(slice[len(ids1):], ids2)
	return slice
}

//交集
func InterTwoInt64s(id1, id2 []int64) []int64 {
	if id1 == nil || id2 == nil || len(id1) == 0 || len(id2) == 0 {
		return []int64{}
	}

	mp := map[int64]bool{}
	for _, item := range id1 {
		mp[item] = true
	}

	interIds := []int64{}
	for _, item := range id2 {
		if _, ok := mp[item]; ok {
			interIds = append(interIds, item)
		}
	}
	return interIds
}

func UniqueInt64s(arr []int64) []int64 {
	mp := make(map[int64]bool)
	arr1 := []int64{}
	for _, id := range arr {
		if _, ok := mp[id]; ok {
			continue
		}
		arr1 = append(arr1, id)
		mp[id] = true
	}
	return arr1
}

func UniqueInt32s(arr []int32) []int32 {
	mp := make(map[int32]bool)
	arr1 := []int32{}
	for _, id := range arr {
		if _, ok := mp[id]; ok {
			continue
		}
		arr1 = append(arr1, id)
		mp[id] = true
	}
	return arr1
}

func UniqueStrs(arr []string) []string {
	mp := make(map[string]bool)
	arr1 := []string{}
	for _, str := range arr {
		if _, ok := mp[str]; ok {
			continue
		}
		arr1 = append(arr1, str)
		mp[str] = true
	}
	return arr1
}

func InArrayInt8(search int8, array []int8) bool {
	for _, i := range array {
		if i == search {
			return true
		}
	}
	return false
}

func InArrayInt(search int, array ...int) bool {
	for _, i := range array {
		if i == search {
			return true
		}
	}
	return false
}

func InArrayInt64(search int64, array []int64) bool {
	for _, i := range array {
		if i == search {
			return true
		}
	}
	return false
}

func InArrayStr(search string, array []string) bool {
	for _, i := range array {
		if i == search {
			return true
		}
	}
	return false
}

func RandomInt64(arr []int64) []int64 {
	dst := make([]int64, 0)
	if len(arr) <= 0 {
		return dst
	}
	dst = append(dst, arr...)
	for i := len(arr) - 1; i >= 0; i-- {
		num := int(coderand.Uint32(uint32(len(arr))))
		dst[i], dst[num] = dst[num], dst[i]
	}
	return dst
}

func MaxInt64(a, b int64) int64 {
	if a >= b {
		return a
	}
	return b
}

func MinInt64(a, b int64) int64 {
	if a <= b {
		return a
	}
	return b
}

func MaxInt(a, b int) int {
	if a <= b {
		return b
	}
	return a
}

func MinInt(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func FilterZero(ids []int64) []int64 {
	ret := []int64{}
	for _, id := range ids {
		if id != 0 {
			ret = append(ids, id)
		}
	}
	return ret
}
