package main

func main() {
	err := CopyFile("./fixture/fixture.js", "./fixture/dest.js", map[string]interface{}{"name": "zcong1993"})
	checkErr(err)
}
