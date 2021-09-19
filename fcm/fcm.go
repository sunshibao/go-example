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
  "project_id": "ry-push",
  "private_key_id": "0022d7c3e06456bc78e972ca878c2f746351a2b1",
  "private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQCP4113/NbypUR2\nS46LYDOOFL79sfzvKsRjwiujmOAakYHhmfbmp8/1JutW8azf3QVOcND/5O6n1t3V\ny1dn+sYXFzLkK6+sB2UdUc9opXEy99NM5GE6P5Z/SCLIOAMVot+b1lCqgmYtXpbo\nBEcVam2MKXu/HHBNrPNCbapyAclNGuf36Q+lQ2AvWPzRNbTZHUUSVfrTul4Ai7kV\nhMCnLYhFPiFIEvEVjXRzeVVH0aocsJwJ3AwLX0LXwEf82Cp5o7GVC4St0HwJ9Y9E\n9RLT+KOzYrPtEmdAZnuY27AQqyoMv7q4L63MFllalv5g72VcgvWqX1ZRa58cmEFm\ne7yTNINtAgMBAAECggEAPT7KgLGq7oJhJzcW+AVdPGthNkwl/MlH6axy9cIzFav8\nzEHEQb6JOnCz+ICLFIiX4yELHPjdiqxfJUrVrAkmexqhS1S7BC7rn7S8Y28URV0O\nYCrPpcL+NLRINc/9pB59LnBlSoaRplseQajcduWjtmbL6PdaZ/2nV42lYLyY6gvX\n8MCIFgdEmzpKeSRIsHA4J19KEnCpr6iVrnJuYOSWb3bRvgkMIMWOfghPdhqgt/Re\nl+rBTDEQIXovzePcEl8i2KGhqA7DybWqhnp2pcfOqKXJDn4r8KfOUwuZWXlwOuVd\nU07K9pUPAmYChoEE9MTjJ0tQKwCklI9nnJIvcN4NwwKBgQDAm0A+nT1VpC88Z+hd\nIlzEf5YDihiCLLGg+lJeyG038ecYVo53pPk3TZSAZJT89q9Kq30rRoJVVqfeky75\n8/fCNXJJFfihIiUvSKD8hxxtURk7UkY8RQlTV7zJXp9xNM78/jXrwirYeKkNGrMU\nZAqsbh26bkfyFbBG31hGidg3uwKBgQC/Py2hDYkplv0d7TrX17Ux+HEhaukvxa4T\n4BOsNynuGjzYnjy2ZVEqYUkLsITXmK51Ec2FE5AEG2yCGTQlIsXKICz/+0OcCblK\n2EkVJXHnxyAlDJUXlBd2Pig+kflHXZJtC1o/iEGiZKS25g31YDWWavMgbLvEb5ga\npeI5uo1a9wKBgDC+mdG+5ajo9nNpMKtxaNzqFUMsDevT99hvwazvOITyGbRVI5Uw\nEUpnve+IhRRIMjDJmdH79Dw4xB9WTItBrTBfal8IynWtOI/w551BbHZWEfarac8T\nAbsv9z1XTy7NWJ06kNbruwAx+UaxvSSp9PGSpL9r7ZVc29Hz8FlZxXjFAoGAZ2VY\n3gHbbkjbgBqJ/bf5lpKjV1XdzJ4rh0NWX31Thg9ZZiPm3xXX0/nU7CT7LTS23URK\nwO/apN3OxGer8YfjN2w96AeIfgwjyXs1x+D+vDjEAEPVN3IxXRQt3eY9x3+ncpz4\npfCcy5duFhQ657akQjaAS799pK5QpdeZf2yEj/ECgYB58usoD6EcezeIA+OLwdkL\nC97tLtHA+Ff08YtLgNT8AGgGq80rNQK6Wv+R8yV20zV1Bjzwu0eMrBBt6xdvmQgs\n3VMahS7IBMsOoT47a0CmEmwP2c3d/O7l5eDGcewI/NKGWLh2D8SUFtJe1ca/yjV1\nnNJGxHiIvxxB+1uT51auJQ==\n-----END PRIVATE KEY-----\n",
  "client_email": "firebase-adminsdk-rgf08@ry-push.iam.gserviceaccount.com",
  "client_id": "103297913384884974542",
  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
  "token_uri": "https://oauth2.googleapis.com/token",
  "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
  "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-rgf08%40ry-push.iam.gserviceaccount.com"
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
