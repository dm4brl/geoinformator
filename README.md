# geoinformator
This code is designed to create a web server that handles HTTP POST requests and sends notifications via Firebase Cloud Messaging (FCM) with geolocation integration. The code does not return any results as JSON or HTML, but returns an HTTP status depending on the result of the operations. Here are the HTTP statuses that can be returned:

If the notification request was successfully processed and the notification was sent, the code will return an HTTP status of 200 OK.

If there was an error while processing the request (for example, invalid JSON or an error while sending the notification), the code will return an HTTP status of 500 Internal Server Error.

If the server started successfully, but requests are sent to an invalid URL, the code will return HTTP status 404 Not Found.

Geo-integration in this code is implemented using Google Maps API. In particular, the calculateDistanceAndTime function uses Google Maps API to calculate the distance and time between two locations (courier location and order location). This information can be included in the notification sent to the courier. Integration with Google Maps API provides geolocation capabilities in your app.

