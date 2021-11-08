package fcm

import (
	"context"
	"log"

	"firebase.google.com/go/messaging"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var RyFcmClient = &messaging.Client{}

func init() {
	var serviceAccountKey = []byte(`{
        "type": "service_account",
		"project_id": "honormarket-2cebf",
		"private_key_id": "a98d100f62a562cf01a91ada0da02989bcb7e673",
		"private_key": "-----BEGIN PRIVATE KEY-----\nMIIEuwIBADANBgkqhkiG9w0BAQEFAASCBKUwggShAgEAAoIBAQC8S0IxvxtT0+ii\noSD/Y83dv1HYS7S81ILd/sqqnO/K3LZkKtjxu+3grQxglvFYow/ZqeZjtVU516Pl\n+5ZA8ibWNsHLP2sVnEU70drvl2qR5TthZD1hGJR5fr2IRQS4yKr1NXfRA4xKcHfJ\nThnew8P7KTYy53npTrDveETuMEcMMoMqvcpURCGxdJD8kVh+BzhJhUjGEvWeAopg\nU2OMCqx3L13ADRK5iFn7OjcxYxKxeG5C6LvBnkqdLfhVvwi3F3QEAW7CPTO3nWVh\ngrVUE/0MEU1ZkNOGuzUymf6S0x0fAR1HhGR1g298URIHarcjuKqducARCDbgGrm+\n/S960PLpAgMBAAECgf9Xtbbf9WWFZvC20NfOH2+GG9jEKH/IdjV/s3A0iWYp2SYd\nAtiLwj7Mqc9aLvW9lDeN7NalG5LXntt6Bfe3yRTONDORubjoGCMn+llBwgaib2V5\nuHffAtQmxCMFlyDb8p2wG3f0pPJ0um8SFo92dTqVhFz62ViwlTdWVW/GhDaUNyTi\nHTKckkBrVMx3PkilEWMIPlD25fEdKj/wc6ctcUeQebwjYpgwqMBy5u3ZN+ar1JFF\npHX2xF3nW7dnm/edi7ayhy3Dy5WagKWEejB9Y1fMClkcFk8KkIN7Zwmi+zuIQScl\n2xkvkoRJWEa2MzSkcb58lEabyAXIZaGyGMbd2wECgYEA/BoVF9Vp9qgf1vhhMeKJ\nuilHFW4aWkiIT439V+csFQ8wz4TNhIPfsFk0TIwgsmtw8QmiEj6xT4oAEASyIFr1\nlNn/morAOPY7iZf1seMry6ZW8H0qSlQ0kbLUhkKbVKMvMcHOjBokZDuwUPR9vRmp\nfhwB6c6WvNWpDUy5aqgj7WkCgYEAvzSZfwaX6ZSxGw5LhddhhFns8wgUtzcxiGrc\nBO0eA0QGgAiKupfdYgWETNuvUxfp/uHhzT29lDNvIBquubjzi/4EPdopCxLFrqme\n0RAjDQq9ztpkkG9xlw4seMQaAmqsfbX+Tla2FtNw2raqmtjmF7lgVIPL5T2aG+De\nZNWcqYECgYEA8j8h67AQTYtKZShxRR05WSCXBLmzKvQtv7xiKBikXGwnbBFh5ydN\nSEi/n5q7RJdHhObLzRpfCV5DJyFMBRlCiNFd8yPHCDVcCqBx4Ii5qcxiGF89xwTZ\nKvQbkhPo7NCN5hMkpa3tMD/G8lOti4tgOiUxlXkFkdkBxBzowttlApkCgYBScvgN\nkmZHTtrf89YVLL7cN3q5ga6NIru1O38MkML0XY0AAK7xrzLDZeWaztBXYGSCiBy6\nR/lLwrIUgez+IQbEQxDJKx2vrLSZkILvW6oBobQfUoUy0xuEK5R5rvYYPK+MtcGn\nPjTeuuZbKZF/nC/74u/rAboWa+3cP6cmjAY2AQKBgBOd36+rx+8r26YQ/jGnhYLq\nwUDBRG20QkbV4VHav4OV+gk6LMf7cc+p8J7Z78/FapywmofpHzTX2IteRodFFSjz\nE8CkM8cNZRCu7YekuVq5/WgS1+d6HRGLDHHvDgAoXB2Ne1qErPstiDdODEju144h\nNrtQWdWtxLenPZuM/KN0\n-----END PRIVATE KEY-----\n",
		"client_email": "firebase-adminsdk-ho3cu@honormarket-2cebf.iam.gserviceaccount.com",
		"client_id": "115467382916161819310",
		"auth_uri": "https://accounts.google.com/o/oauth2/auth",
		"token_uri": "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-ho3cu%40honormarket-2cebf.iam.gserviceaccount.com"
     }`)

	opt := option.WithCredentialsJSON(serviceAccountKey)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	// Obtain a messaging.Client from the App.
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}
	RyFcmClient = client
}
