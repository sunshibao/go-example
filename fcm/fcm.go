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
  "private_key_id": "c3c2804df3c8d29a3e18c1c1d5f828e8f83eb5fb",
  "private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCzxmNkjOstyuTE\nPWh1oZxR0z/9VzUPpG912Y7Stj/lfyV/IfBExW4gLgZdeaAwfsUCJiWjMe+3FY3S\nCyJSeNnaD7rGYUyr2EMFocjhXLGVDXX+/W0gvGXDHt87PKlMKWQ3f1JI+K1ISH9J\nKPNnurrgaO8AfWkG5vIEtls1N/T4SsgYQgi9hnAIrk6uL9zRRG3xlOTWO7seYKwq\nOLPAHQZaMMeX18GbkH6B2jVXlblkEc25tKduHmvg94zIC31cb7NLsKJ50qiZcs5o\nXBYGLEiSnRmPVbmR0lWqrj1iHmyyomG1/kZdJQ3LtBuAXqbRbNBm6bh+DsNkdjkR\nKw4e8leTAgMBAAECggEAVqQZSN+FhJdRM6aoznBp5yGZqF48K2LDeRe19qVxHAjw\nKFpR8sL39ThZRvmUE2s1RXjfEMzDTAhuRjmpe/fUfnywNmW0Tum2O2OibcJ/r2dC\nY/b9fhZuMOCTE3iD2znjm7+YB6UZ1kliVv8OeAKNiBPqg+DtGi95mn3MbVSfkn2h\nDEmOwgTiZdtuQgca7YRnTtGGyX7HcVd8LoO5AMYsCA+AzhYeVfiwsisVsTMvY1Hs\nGZF/SAw0Fle6FhJ5dvlEs2YOGuJUDzpLgDoymEppmtoexxf6hzH9m5ZLEMEVVkwS\noLKz9OXrFlpuPeB3jEuEZF7IF3okuGs2UOEWX61/0QKBgQDe7N9osyKsm+F3Stxw\nAGyd1OxDxyObJ4Sr4Tcp9nEi+vhKmv27eySNLu0Ee4GWbPFnDGxxkGSj/KN3v3hy\nTcN7Q0ti2b56yKeQ0Q23DCYtNzaCdZQYxMeN6xvwkfRBoBEcioRTPK5AqQo/Ut7k\nIBAC8GXZW5snPl09Iu+5CPIqCQKBgQDOcpVxoEOlxT6aklUkyekUB+9uNRH5PLGR\nDPS8br8yNWCy49kWotMZYWuVR5w1ZGvnqF2To0rBgj3rV6jXcswfjMlRezq+WX3r\nhcTkvU7yx9uEgXirQVBZGcpROuRCywlZGicu33IkF8xoR7zYZnZu+XW9nBdbIALe\nDe126RVLuwKBgQCUJlOq3zHyCH5kqymofX+xvKlvuc5d9HtlRv1EM/WoS6Xq3uo1\nSHdYJJF74yCR+cDMvSLsrgtUODfjAypmvGtnZaeaAB9otiU2RtiYh8hvUDw3ozBS\n7aO5G0CsNjNf0aLNlvit4KvlaPHYd5iNfwllCqKFOy33fKi2UpuCGwHlEQKBgB1J\nUfvxgUc16QCie6OhZQgrbALVUnxp6MKr9Nf6WfM0cUhPXE+Cv7GdCVb+9qD4YvpU\n/xlfk85JvKDzKYeOyHg9T53YfmfcaCmOK3VLibSVN4XfnA+nT0+kgffuA82Z4fSa\n6i4TEq0eFyg/7QFB39E0YEiKqSGLM/zuOt6giB3bAoGBAMYZqnv1jrdveydEosrc\nMJKqZqDx2/UkmmvjScP5lXP9/YAEp19SfzUkLokzbGxvKQNq3PabbEyQgoVLXSB8\nqyKr8JIXBgIIVyWvL3hNcBZzXnCfZ7H69Ogix24wQfo+43SjSLccZNeXgdPKDN29\nK4zq4G2VA+5BEyg819UG/xk8\n-----END PRIVATE KEY-----\n",
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
