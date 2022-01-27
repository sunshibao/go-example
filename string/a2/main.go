package main

import (
	"fmt"
	"regexp"
)

var re = regexp.MustCompile("^[a-zA-Z0-9\\s\\./:●~!@#$%^&*\\+\\-(){}|<>=✔【】:\\\"?'：；‘’“”，。,、\\]\\[`《》]+$").MatchString

func main() {
	fmt.Println(re("more...\n✔\nMaking a vegan breakfast is easier.\n✔\nWant"))

}
