package manager

import (
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"iproute2/model"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const command = "ip"

func CreateRouteWithIfIP(r model.Route)  {
	dest := r.Destination + string("/") + strconv.Itoa(r.DestCIDR)
	args := []string{"route", "add", dest, "via", r.InterfaceIP}
	cmdOut := ExecuteIPCommand(args)
	log.Infof(cmdOut)


}

func RemoveDefaultGateway(r model.Route)  {
	args := []string{"route", "delete", "default", "via", r.InterfaceIP}
	cmdOut := ExecuteIPCommand(args)
	log.Infof(cmdOut)

}

func CreateDefaultGateway(r model.Route)  {
	args := []string{"route", "add", "default", "via", r.InterfaceIP}
	cmdOut := ExecuteIPCommand(args)
	log.Infof(cmdOut)

}

func RemoveRoute(r model.Route)  {
	dest := r.Destination + string("/") + strconv.Itoa(r.DestCIDR)
	args := []string{"route", "delete", dest, "via", r.InterfaceIP}
	cmdOut := ExecuteIPCommand(args)
	log.Infof(cmdOut)

}

func GetRoutes() (rts string) {
	args := []string{"route", "show"}
	cmdOut := ExecuteIPCommand(args)
	s := ParseStringRoutes(cmdOut)
	routes, err := json.Marshal(s)
	if err != nil {
		log.Error(err.Error())
		return ""
	} else {
		return string(routes)
	}

}

func ParseStringRoutes(cmdOutput string) (parsedRoutes []model.Route) {
	singleRoutes := strings.Split(cmdOutput, "\n")
	routes := []model.Route{}
	for i := 0; i < (len(singleRoutes) - 1); i++ {
		route := strings.Split(singleRoutes[i], " ")
		var r = model.Route{}
		for y := 0; y < len(route); y++ {
			if y == 0 {

				if route[y] == "default" {
					r.Destination = route[y]
				} else {
					destAndCidr := strings.Split(route[y], "/")
					r.Destination = destAndCidr[0]
					r.DestCIDR, _ = strconv.Atoi(destAndCidr[1])
				}
			}
			if route[y] == "via" || route[y] == "src" {
				r.InterfaceIP = route[y + 1]
			}
		}
		routes = append(routes, r)
	}
	return routes
}

func GetInterfaces() (ifs string)  {
	args := []string{"addr", "show"}
	cmdOut := ExecuteIPCommand(args)
	ifsNames := ParseIfs(cmdOut)
	interfaces, err := json.Marshal(ifsNames)
	if err != nil {
		log.Error(err.Error())
		return ""
	} else {
		return string(interfaces)
	}
}

func ParseIfs(cmdOut string) (ifNames []model.Interface)  {
	ifsNames := []model.Interface{}
	ifSlice := strings.Split(cmdOut, ": ")
	x := 0

	for i := 0; i < (len(ifSlice) - 1); i += 2{
		ifsNames = append(ifsNames, model.Interface{Name:ifSlice[i + 1]})
		s := strings.Split(ifSlice[i + 2], " ")
		for y := 0; y < len(s); y++ {
			if s[y] == "inet" {
				ip := strings.Split(s[y + 1], "/")[0]
				ifsNames[x].IPAddress = ip
				x++
				break
			}
		}
	}
	return ifsNames
}


func ExecuteIPCommand(args []string) (cmdOut string){
	cmd := exec.Command(command, args...)
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err := cmd.Run()
	log.SetOutput(os.Stdout)
	if err != nil {
		switch err := err.(type) {
		case *exec.ExitError:
			log.Errorf("Program exited with %s (User problem)", err.Error())
		default:
			log.Errorf("Error occurred! %s (API problem)", err.Error())
		}

	} else {
		return string(cmdOutput.Bytes())
	}
	return ""
}