package bmbuilder

import (
	"fmt"
	"log"
)

func purple(in ...interface{}) string {
	ins := fmt.Sprint(in...)
	return "\033[35m" + ins + "\033[0m"
}
func green(in ...interface{}) string {
	ins := fmt.Sprint(in...)
	return "\033[32m" + ins + "\033[0m"
}
func yellow(in ...interface{}) string {
	ins := fmt.Sprint(in...)
	return "\033[33m" + ins + "\033[0m"
}
func blue(in ...interface{}) string {
	ins := fmt.Sprint(in...)
	return "\033[34m" + ins + "\033[0m"
}
func red(in ...interface{}) string {
	ins := fmt.Sprint(in...)
	return "\033[31m" + ins + "\033[0m"
}
func cyan(in ...interface{}) string {
	ins := fmt.Sprint(in...)
	return "\033[36m" + ins + "\033[0m"
}
func gray(in ...interface{}) string {
	ins := fmt.Sprint(in...)
	return "\033[37m" + ins + "\033[0m"
}
func white(in ...interface{}) string {
	ins := fmt.Sprint(in...)
	return "\033[97m" + ins + "\033[0m"
}

// Debug helper functrion
func (bi *BMBuilder) Debug(logline ...interface{}) {
	if bi.debug {
		log.Println(purple("[Debug]")+" -", logline)
	}
}

// Info helper function
func (bi *BMBuilder) Info(logline ...interface{}) {
	if bi.verbose || bi.debug {
		log.Println(green("[Info]")+" -", logline)
	}
}

// Warning helper function
func (bi *BMBuilder) Warning(logline ...interface{}) {
	log.Println(yellow("[Warn]")+" -", logline)
}

// Alert helper function
func (bi *BMBuilder) Alert(logline ...interface{}) {
	log.Println(red("[Alert]")+" -", logline)
}
