package main

import (
	"api/iproute2/manager"
	"api/iproute2/model"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic"
	gelf "github.com/seatgeek/logrus-gelf-formatter"
	"github.com/sirupsen/logrus"
	"gopkg.in/sohlich/elogrus.v3"
	"net/http"
	"strconv"
	"strings"
)
var log = logrus.New()

func init() {
	log.Formatter = new(gelf.GelfFormatter)
	log.Level = logrus.InfoLevel
}
// Homepage endpoint
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "IPRoute2 controller!")
}

// CheckParams checks if some index of string array is empty.
func checkParams(params []string) bool {
	for _, i := range params {
		if i == "" {
			return true
		}
	}
	return false
}

// ModRoute provides creation and deletion route.
// Func is called on /iproutes/mod endpoint.
func modRoute(w http.ResponseWriter, r *http.Request) {

	params := []string{r.URL.Query().Get("destination"),
		r.URL.Query().Get("mask"), r.URL.Query().Get("interface")}

	if checkParams(params) {
		log.Error("Url Params are missing")
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	route := model.Route{Destination: model.Network{params[0], 0}, InterfaceIP: params[2]}
	route.Destination.Mask, _ = strconv.Atoi(params[1])

	switch r.Method {

	case http.MethodPut:
		manager.CreateRouteWithIfIP(route)
		log.WithFields(logrus.Fields{
			"destination": route.Destination.IP,
			"dest CIDR":   route.Destination.Mask,
			"Interface":   route.InterfaceIP,
		}).Info("Created new route!")
	case http.MethodDelete:
		manager.RemoveRoute(route)
		log.WithFields(logrus.Fields{
			"destination": route.Destination.IP,
			"dest CIDR":   route.Destination.Mask,
			"Interface":   route.InterfaceIP,
		}).Info("Deleted route!")
	default:
		log.Error("Unsupported method!")
		w.WriteHeader(404)
	}

}

// ModDefaultRoute provides creation and deletion default route.
// Func is called on /iproutes/default endpoint.
func modDefaultRoute(w http.ResponseWriter, r *http.Request) {

	params := []string{r.URL.Query().Get("interface")}

	route := model.Route{InterfaceIP: params[0]}
	log.Infof(route.InterfaceIP)
	switch r.Method {

	case http.MethodPut:
		if checkParams(params) {
			log.Error("Url Params are missing")
			return
		}
		manager.CreateDefaultGateway(route)
		log.WithFields(logrus.Fields{
			"Interface": route.InterfaceIP,
		}).Info("Created default route!")
	case http.MethodDelete:
		if checkParams(params) {
			manager.RemoveDefaultGateway()
			log.Info("Deleted default route")
		} else {
			manager.RemoveDefaultGatewayVia(route)
			log.WithFields(logrus.Fields{
				"Interface": route.InterfaceIP,
			}).Info("Deleted default route!")
		}

	default:
		log.Error("Unsupported method!")
		w.WriteHeader(404)
	}

}

// GetRoutes returns JSON array of routes in string
func getRoutes(w http.ResponseWriter, r *http.Request) {
	routes := manager.GetRoutes()
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(routes))
}

// GetInterfaces returns JSON array of interfaces in string
func getInterfaces(w http.ResponseWriter, r *http.Request) {
	interfaces := manager.GetInterfaces()
	w.Header().Set("Content-Type", "application/json")
	log.WithFields(logrus.Fields{
		"test": "test",
	}).Info("Returned Interfaces!")
	w.Write([]byte(interfaces))
}

// SetElasticLog sets hook to elasticsearch server for logrus.
// Requires parameter elastic with value IPv4 address of elasticsearch server.
func setElasticLog(w http.ResponseWriter, r *http.Request) {

	params := []string{r.URL.Query().Get("elastic")}
	interfaces := []model.Interface{}
	json.Unmarshal([]byte(manager.GetInterfaces()), interfaces)

	var host string

	for i := 0; i < len(interfaces); i++ {
		if strings.Contains(interfaces[i].Name, "eth0") {
			host = interfaces[i].IPAddress
		} else {
			host = "FAIL"
		}
	}

	if checkParams(params) {
		log.Error("Elastic Url is missing")
		return
	}
	client, err := elastic.NewClient(elastic.SetURL("http://" + params[0]))
	if err != nil {
		log.Panic(err)
	}
	hook, err := elogrus.NewElasticHook(client, host, logrus.InfoLevel, "logs-home")
	if err != nil {
		log.Panic(err)
	}
	log.Hooks.Add(hook)
	log.Info("Created hook to elastic log")
	w.WriteHeader(200)

}

// HandleRequests provides endpoint handling.
func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/iproutes/mod", modRoute)
	http.HandleFunc("/iproutes/default", modDefaultRoute)
	http.HandleFunc("/iproutes", getRoutes)
	http.HandleFunc("/interfaces", getInterfaces)
	http.HandleFunc("/setLogHook", setElasticLog)
	log.Fatal(http.ListenAndServe(":8090", nil))
}

func main() {
	handleRequests()
}
