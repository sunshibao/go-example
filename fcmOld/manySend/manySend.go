package main

import (
	"fmt"
	fcm "myGo/fcmOld"
)

const (
	serverKey = "AAAAbCW6n1Q:APA91bFTctykFW5QkzDBXnt3nES7jgqD_JQ0VY1qUf8nffd3_nZrQplP3sqKf19HehL8Czs6zfI3RhKqhmlYqllRS0rQdbkgHr4o74--y7ZoXodQcMohEfiMxGJhaqu5EKPQa2AQZusf"
)

func main() {

	data := map[string]string{
		"msg": "Hello World1",
		"sum": "Happy Day",
	}

	ids := []string{
		"368dde283db539abc4a6419b1795b6131194703b816e4f624ffa12",
	}

	xds := []string{
		"368dde283db539abc4a6419b1795b6131194703b816e4f624ffa12",
		"368dde283db539abc4a6419b1795b6131194703b816e4f624ffa12",
		"368dde283db539abc4a6419b1795b6131194703b816e4f624ffa12",
	}

	c := fcm.NewFcmClient(serverKey)
	c.NewFcmRegIdsMsg(ids, data)
	c.AppendDevices(xds)

	status, err := c.Send()

	if err == nil {
		status.PrintResults()
	} else {
		fmt.Println(err)
	}

}
