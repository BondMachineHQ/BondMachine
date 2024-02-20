package main

import "fmt"

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
