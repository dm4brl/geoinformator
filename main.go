package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"googlemaps.github.io/maps"
)

type NotificationData struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Token string `json:"token"`
}

type LocationData struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func sendNotification(w http.ResponseWriter, r *http.Request) {
	var notificationData NotificationData

	if err := json.NewDecoder(r.Body).Decode(&notificationData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	opt := option.WithCredentialsFile("path/to/your/firebase-service-account-key.json")
	config := &firebase.Config{ProjectID: "your-firebase-project-id"}
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: notificationData.Title,
			Body:  notificationData.Body,
		},
		Token: notificationData.Token,
	}

	// Отправка уведомления
	response, err := client.Send(ctx, message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Successfully sent message: %s", response)
}

func calculateDistanceAndTime(courierLocation, orderLocation LocationData) (float64, string, error) {
	apiKey := "your-google-maps-api-key"
	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		return 0, "", err
	}

	origin := fmt.Sprintf("%f,%f", courierLocation.Latitude, courierLocation.Longitude)
	destination := fmt.Sprintf("%f,%f", orderLocation.Latitude, orderLocation.Longitude)

	r := &maps.DistanceMatrixRequest{
		Origins:      []string{origin},
		Destinations: []string{destination},
		Mode:         maps.TravelModeDriving,
	}

	resp, err := client.DistanceMatrix(context.Background(), r)
	if err != nil {
		return 0, "", err
	}

	if len(resp.Rows) == 0 || len(resp.Rows[0].Elements) == 0 {
		return 0, "", fmt.Errorf("Could not calculate distance and time")
	}

	distance := resp.Rows[0].Elements[0].Distance.Meters
	duration := resp.Rows[0].Elements[0].Duration.Hours() // Получаем время в часах

	return float64(distance), fmt.Sprintf("%.2f часов", duration), nil
}

func main() {
	http.HandleFunc("/sendNotification", sendNotification)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server is running on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil)
}
