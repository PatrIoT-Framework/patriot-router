// Package model provides structures representing objects in routing tables
// Network interfaces and routes.
package model

// Interface represents network interface of workstation/container/...
type Interface struct {
	IPAddress string
	Name      string
}
