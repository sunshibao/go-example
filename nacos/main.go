/*
createTime: 2021/11/3
*/
package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func main() {
	nacosTest2()
}

// 我通过example的源码 创建一个真正的注册中心
func nacosTest2() {
	sc := []constant.ServerConfig{
		{
			IpAddr: "127.0.0.1",
			Port:   8849,
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         "b123c1b0-b68d-4ae1-957a-50a09e7f21c5", //namespace id
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/Users/liqp/fsdownload/gin/logs",
		CacheDir:            "/Users/liqp/fsdownload/gin/logs",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}

	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)

	if err != nil {
		panic(err)
	}

	//监听服务
	err = client.Subscribe(&vo.SubscribeParam{
		ServiceName: "demo.go",
		//GroupName:   "group-a",             // 默认值DEFAULT_GROUP
		Clusters: []string{"zwt"}, // 默认值DEFAULT
		SubscribeCallback: func(services []model.SubscribeService, err error) {
			b, _ := json.Marshal(services) //todo 注意每次的消息计算md5用以比对下次
			fmt.Println("---------------", string(b))
		},
	})
	time.Sleep(time.Hour)
}

//获取在线服务
func d(client naming_client.INamingClient) {
	service, _ := client.SelectInstances(vo.SelectInstancesParam{
		ServiceName: "demo.go",
		Clusters:    []string{"zwt"}, // 默认值DEFAULT
		HealthyOnly: true,
	})
	b, _ := json.Marshal(service)
	fmt.Println(string(b))
	return
}

//服务注册
func r(client naming_client.INamingClient) {
	success, err := client.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          "127.0.0.1",
		Port:        9000,
		ServiceName: "demo.go",
		ClusterName: "zwt",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"idc": "shanghai"}, //填写路由配置
	})

	if !success {
		fmt.Printf("注册nacos失败,%v\n", err)
		return
	}
}
