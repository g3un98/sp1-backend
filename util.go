package main

func containAdminMembers(s []member, v string) bool {
	for _, vv := range s {
		if v == vv.AppId && vv.IsAdmin == 1 {
			return true
		}
	}
	return false
}

func containMembers(s []member, v string) bool {
	for _, vv := range s {
		if v == vv.AppId {
			return true
		}
	}
	return false
}
