// Package manager provides management of routing tables.
package manager

import (
	"api/iproute2/model"
	"bytes"
	"encoding/json"
	gelf "github.com/seatgeek/logrus-gelf-formatter"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const command = "ip"

var log = logrus.New()

func init() {
	log.Formatter = new(gelf.GelfFormatter)
	log.Level = logrus.InfoLevel
}

// CreateRouteWithIfIP prepares command to create route and executes it with ExecuteIPCommand func.
func CreateRouteWithIfIP(r model.Route) {
	dest := r.Destination.IP + string("/") + strconv.Itoa(r.Destination.Mask)
	args := []string{"route", "add", dest, "via", r.InterfaceIP}
	cmdOut := ExecuteIPCommand(args)
	log.Infof(cmdOut)

}

// RemoveDefaultGateway prepares command to remove default route and executes it with ExecuteIPCommand func.
func RemoveDefaultGatewayVia(r model.Route) {
	args := []string{"route", "delete", "default", "via", r.InterfaceIP}
	cmdOut := ExecuteIPCommand(args)
	log.Infof(cmdOut)

}

func RemoveDefaultGateway() {
	args := []string{"route", "delete", "default"}
	cmdOut := ExecuteIPCommand(args)
	log.Infof(cmdOut)

}

// CreateDefaultGateway prepares command to create default route and executes it with ExecuteIPCommand func.
func CreateDefaultGateway(r model.Route) {
	args := []string{"route", "add", "default", "via", r.InterfaceIP}
	cmdOut := ExecuteIPCommand(args)
	log.Infof(cmdOut)

}

// RemoveRoute prepares command to remove route and executes it with ExecuteIPCommand func.
func RemoveRoute(r model.Route) {
	dest := r.Destination.IP + string("/") + strconv.Itoa(r.Destination.Mask)
	args := []string{"route", "delete", dest, "via", r.InterfaceIP}
	cmdOut := ExecuteIPCommand(args)
	log.Infof(cmdOut)

}

// GetRoutes prepares command to list all routes and executes it with ExecuteIPCommand func.
// ExecuteIPCommand output is parsed via ParseStringRoutes func.
// Returns model/route array serialized into JSON in string.
func GetRoutes() (rts string) {
	args := []string{"route", "show"}
	cmdOut := ExecuteIPCommand(args)
	s := ParseStringRoutes(cmdOut)
	routes, err := json.Marshal(s)
	if err != nil {
		log.Error(err.Error())
		return ""
	}
	return string(routes)

}

// ParseStringRoutes parses ExecuteIpCommand (ip route show) output to array of model/route.
// If route is default route,then only Destination and Interface attributes will be written,
// DestCIDR will be null.
// Returns array of model/route.
func ParseStringRoutes(cmdOutput string) (parsedRoutes []model.Route) {
	singleRoutes := strings.Split(cmdOutput, "\n")
	routes := []model.Route{}
	for i := 0; i < (len(singleRoutes) - 1); i++ {
		route := strings.Split(singleRoutes[i], " ")
		var r = model.Route{}
		for y := 0; y < len(route); y++ {
			if y == 0 {

				if route[y] == "default" {
					r.Destination.IP = route[y]
				} else {
					destAndCidr := strings.Split(route[y], "/")
					r.Destination.IP = destAndCidr[0]
					r.Destination.Mask, _ = strconv.Atoi(destAndCidr[1])
				}
			}
			if route[y] == "via" || route[y] == "src" {
				r.InterfaceIP = route[y+1]
			}
		}
		routes = append(routes, r)
	}
	return routes
}

// GetInterfaces prepares command to list all network interfaces and executes it with ExecuteIPCommand func.
// ExecuteIPCommand output is parsed via ParseIfs func.
// Returns model/interface array serialized into JSON, as string.
func GetInterfaces() (ifs string) {
	args := []string{"addr", "show"}
	cmdOut := ExecuteIPCommand(args)
	ifsNames := ParseIfs(cmdOut)
	interfaces, err := json.Marshal(ifsNames)
	if err != nil {
		log.Error(err.Error())
		return ""
	}
	return string(interfaces)
}

// ParseIfs parses ExecuteIpCommand (ip addr show) output to array of model/route.
// Returns array of model/interface.
func ParseIfs(cmdOut string) (ifNames []model.Interface) {
	ifsNames := []model.Interface{}
	ifSlice := strings.Split(cmdOut, ": ")
	x := 0

	for i := 0; i < (len(ifSlice) - 1); i += 2 {
		ifsNames = append(ifsNames, model.Interface{Name: ifSlice[i+1]})
		s := strings.Split(ifSlice[i+2], " ")
		for y := 0; y < len(s); y++ {
			if s[y] == "inet" {
				ip := strings.Split(s[y+1], "/")[0]
				ifsNames[x].IPAddress = ip
				x++
				break
			}
		}
	}
	return ifsNames
}

// ExecuteIPCommand executes command with arguments.
// If command is executed successfully, then func returns
// command line output of command. If command is not executed successfully,
// then func returns actual error!
func ExecuteIPCommand(args []string) (cmdOut string) {
	cmd := exec.Command(command, args...)
	cmdOutput := &bytes.Buffer{}
	var stderr bytes.Buffer
	cmd.Stdout = cmdOutput
	cmd.Stderr = &stderr
	err := cmd.Run()
	log.SetOutput(os.Stdout)
	if err != nil {
		switch err := err.(type) {
		case *exec.ExitError:
			log.WithFields(logrus.Fields{
				"Error exited with user error: " : err.Error(),
				"Stacktrace" : stderr.String()}).Error("Error ocurred!")
		default:
			log.WithFields(logrus.Fields{
				"Error exited with API error: " : err.Error(),
				"Stacktrace" : stderr.String()}).Error("Error ocurred!")
		}
		return err.Error()
	}
	return string(cmdOutput.Bytes())
}
