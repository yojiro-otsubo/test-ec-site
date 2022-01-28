package others

import "os"

func IconFilePathCheck(user_id string) bool {
	filepath := "app/static/img/icon/userid" + user_id + "/" + "icon.jpg"
	if f, err := os.Stat(filepath); os.IsNotExist(err) || f.IsDir() {
		return false
	} else {
		return true
	}
}
