package configs

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var FirebaseStorageBucketName = "multi-aura-8eb80.appspot.com"

func InitializeFirebaseApp() *firebase.App {
	opt := option.WithCredentialsFile("./internal/configs/firebase/multi-aura-8eb80-firebase-adminsdk-sdex1-090edd6c44.json")

	app, err := firebase.NewApp(context.Background(), &firebase.Config{
		StorageBucket: FirebaseStorageBucketName,
	}, opt)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase app: %v", err)
	}

	return app
}
