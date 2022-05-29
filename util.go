package main

func contain(s []member, v string) bool {
	for _, vv := range s {
		if v == vv.AppId {
			return true
		}
	}
	return false
}
