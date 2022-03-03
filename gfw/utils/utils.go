package utils

import (
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ConvertIntSet(set *schema.Set) []int {
	s := make([]int, 0, set.Len())
	for _, v := range set.List() {
		s = append(s, v.(int))
	}
	sort.Ints(s)

	return s
}

func Exists(i int, array []int) bool {

	for _, v := range array {
		if i == v {
			return true
		}
	}
	return false
}
