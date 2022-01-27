package main

func main() {
	defer func() {
		recover()
	}()
	aaa()

}

func aaa() {
	bbb()
}

func bbb() {
	ccc()
}

func ccc() {
	panic("aadddssss")
}
