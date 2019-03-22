// Package model provides structures representing objects in routing tables
// Network interfaces and routes.
package model


// Route represents actual record in linux routing table.
type Route struct {
	Destination Network
	InterfaceIP string
}
