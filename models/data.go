package models

import "sync"

// TEMP DATA STORAGE MODEL
type Data struct {
	Sessions map[string]uint
	Users    map[string]User
	Teams    map[uint]Team
	Mu       *sync.RWMutex
}
