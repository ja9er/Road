package common

import "strings"

func Fofatoshodan(fofaquery string) (shodanquery string) {
	shodanquery = strings.Replace(fofaquery, "=", ":", -1)
	if strings.Contains(shodanquery, "title") {
		shodanquery = strings.Replace(shodanquery, "title", "http.title", -1)
	} else if strings.Contains(shodanquery, "body") {
		shodanquery = strings.Replace(shodanquery, "body", "http.html", -1)
	} else if strings.Contains(shodanquery, "statucs") {

	}

	return shodanquery
}
