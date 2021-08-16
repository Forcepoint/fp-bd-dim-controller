package util

import "regexp"

//*************************************************************************************************
//	########  ########  ######   ######## ##     ##    ##     ##    ###     ######   ####  ######
//	##     ## ##       ##    ##  ##        ##   ##     ###   ###   ## ##   ##    ##   ##  ##    ##
//	##     ## ##       ##        ##         ## ##      #### ####  ##   ##  ##         ##  ##
//	########  ######   ##   #### ######      ###       ## ### ## ##     ## ##   ####  ##  ##
//	##   ##   ##       ##    ##  ##         ## ##      ##     ## ######### ##    ##   ##  ##
//	##    ##  ##       ##    ##  ##        ##   ##     ##     ## ##     ## ##    ##   ##  ##    ##
//	##     ## ########  ######   ######## ##     ##    ##     ## ##     ##  ######   ####  ######
//*************************************************************************************************

// Here be dragons. Thou art forewarned

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])+)?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])+)$")
var ipRegex = regexp.MustCompile("^\\b\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\b$")
var urlRegex = regexp.MustCompile("^(http(s)?://)?(www\\.)?([-a-zA-Z0-9@:%_+~#=]{2,256}\\.[a-z]{2,256}\\b([-a-zA-Z0-9@:%_+~#?&/=]*))+(\\.[a-z]{2,6}\\b([-a-zA-Z0-9@:%_+~#?&/=]*))?$")
var domainRegex = regexp.MustCompile("^([a-z0-9]+(-[a-z0-9]+)*\\.)+[a-z]{2,}$")
var rangeRegex = regexp.MustCompile("^(\\b([1-9][0-9]{0,2}?|255)\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\b)-(\\b([1-9][0-9]{0,2}?|255)\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\b)$")

func IsEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}

func IsIpValid(i string) bool {
	return ipRegex.MatchString(i)
}

func IsDomainValid(d string) bool {
	return domainRegex.MatchString(d)
}

func IsUrlValid(u string) bool {
	return urlRegex.MatchString(u)
}

func IsRangeValid(r string) bool {
	return rangeRegex.MatchString(r)
}
