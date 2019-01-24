// Package provides structures representing objects in routing tables
// Network interfaces and routes.
package model

// Object represents network interface of workstation/container/...
type Interface struct {
	IPAddress string
	Name      string
}
