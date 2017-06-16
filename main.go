package main

import "io/ioutil"

func main() {
	data, err := ioutil.ReadFile("config.json")
	checkErr(err)
	parse(data)
}
