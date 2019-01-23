package main

import (
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic"
	"github.com/sirupsen/logrus"
	"gopkg.in/sohlich/elogrus.v3"
	"iproute2/manager"
	"iproute2/model"
	"net/http"
	"strconv"
	"strings"
)
var log = logrus.New()
func homePage(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func checkParams(params []string) bool{
	for _,i := range params{
		if i == "" {
			return true
		}
	}
	return false
}

func modRoute(w http.ResponseWriter, r *http.Request)  {

	params := []string{r.URL.Query().Get("destination"),
		r.URL.Query().Get("mask"), r.URL.Query().Get("interface")}

	if checkParams(params) {
		log.Error("Url Params are missing")
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	} else {
		route := model.Route{Destination: params[0], InterfaceIP: params[2]}
		route.DestCIDR, _ = strconv.Atoi(params[1])

		switch r.Method {

		case http.MethodPut:
			manager.CreateRouteWithIfIP(route)
			log.WithFields(logrus.Fields{
				"destination" : route.Destination,
				"dest CIDR" : route.DestCIDR,
				"Interface" : route.InterfaceIP,
			}).Info("Created new route!")
		case http.MethodDelete:
			manager.RemoveRoute(route)
			log.WithFields(logrus.Fields{
				"destination" : route.Destination,
				"dest CIDR" : route.DestCIDR,
				"Interface" : route.InterfaceIP,
			}).Info("Deleted route!")
		default:
			log.Error("Unsupported method!")
			w.WriteHeader(404)
		}
	}
}

func modDefaultRoute(w http.ResponseWriter, r *http.Request)  {

	params := []string{r.URL.Query().Get("interface")}

	if checkParams(params) {
		log.Error("Url Params are missing")
		return
	} else {
		route := model.Route{InterfaceIP:params[0]}
		log.Infof(route.InterfaceIP)
		switch r.Method {

		case http.MethodPut:
			manager.CreateDefaultGateway(route)
			log.WithFields(logrus.Fields{
				"Interface" : route.InterfaceIP,
			}).Info("Created default route!")
		case http.MethodDelete:
			manager.RemoveDefaultGateway(route)
			log.WithFields(logrus.Fields{
				"Interface" : route.InterfaceIP,
			}).Info("Deleted default route!")
		default:
			log.Error("Unsupported method!")
			w.WriteHeader(404)
		}
	}
}
func getRoutes(w http.ResponseWriter, r *http.Request)  {
	routes := manager.GetRoutes()
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(routes))
}

func getInterfaces(w http.ResponseWriter, r *http.Request)  {
	interfaces := manager.GetInterfaces()
	w.Header().Set("Content-Type", "application/json")
	log.WithFields(logrus.Fields{
		"test":"test",
	}).Info("Returned Interfaces!")
	w.Write([]byte(interfaces))
}

func setElasticLog(w http.ResponseWriter, r *http.Request)  {

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
	} else {
		client, err := elastic.NewClient(elastic.SetURL("http://" + params[0]))
		if err != nil {
			log.Panic(err)
		}
		hook, err := elogrus.NewElasticHook(client, host, logrus.InfoLevel, "networklogs")
		if err != nil {
			log.Panic(err)
		}
		log.Hooks.Add(hook)
		log.Info("Created hook to elastic log")
	}
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/iproutes/mod", modRoute)
	http.HandleFunc("/iproutes/default", modDefaultRoute)
	http.HandleFunc("/iproutes", getRoutes)
	http.HandleFunc("/interfaces", getInterfaces)
	http.HandleFunc("/setLogHook", setElasticLog)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func main() {
	handleRequests()
}