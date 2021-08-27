//Client.go
package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"
)

func ClientHandleError(err error, when string){
	if err != nil {
		fmt.Println(err, when)
		os.Exit(1)
	}

}

func main()  {
	for i:= 0; i < 10000; i++ {
		go tcp()
		//fmt.Println(seed())
	}
	time.Sleep(35*time.Second)
	fmt.Println("启动完成")
	time.Sleep(time.Hour)
}

func seed() time.Duration {
	rand.Seed(time.Now().UnixNano())
	return time.Duration( rand.Intn(30)) * 1E9
}
func tcp(){

	time.Sleep(seed())
	//拨号远程地址，简历tcp连接
	conn, err := net.Dial("tcp","192.168.1.240:8888")
	ClientHandleError(err, "client conn error")

	//预先准备消息缓冲区
	buffer := make([]byte,1024)


	conn.Write([]byte("lineBytes"))
	//准备命令行标准输入
	reader := bufio.NewReader(os.Stdin)

	for {
		lineBytes,_,_ := reader.ReadLine()
		conn.Write(lineBytes)
		n,err := conn.Read(buffer)
		ClientHandleError(err, "client read error")

		serverMsg := string(buffer[0:n])
		fmt.Printf("服务端msg",serverMsg)
		if serverMsg == "bye" {
			break
		}

	}

}