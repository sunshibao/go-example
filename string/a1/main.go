package main

import (
	"fmt"
	"regexp"
	"strings"
)

func main() {
	aaa := "A Amakha Paris agora tem um app moderno e prático para que você tenha tudo na palma da sua mão.<br><br>Logo ao abrir o app você verá um banner rotativo com as novidades. <br><br>Além disso, você terá acesso ao Escritório Virtual, saberá mais informações sobre a empresa, ficará por dentro de todas as campanhas e eventos. <br><br>Um diferencial é você conhecer todos os nosso produtos, acessando o catálogo, além de descobrir qual o seu estilo de perfume.<br><br>No menu existe uma parte especial apenas para os nossos perfumes, com todas as características e caminhos olfativos.<br><br>Você também ficará sabendo de todas as novidades e campanhas, acessará as nossas unidades e poderá entrar em contato conosco.<br><br>Tenha em apenas um clique todas as vantagens de ser um Executivo Amakha Paris!"

	fmt.Println(TrimHtml(aaa))
}

func TrimHtml(src string) string {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)
	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")
	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")
	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")
	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")
	return strings.TrimSpace(src)
}
