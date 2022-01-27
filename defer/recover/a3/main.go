package main

func main() {
	defer func() {
		defer func() {
			recover()
		}()
	}()
	panic(1)
}
