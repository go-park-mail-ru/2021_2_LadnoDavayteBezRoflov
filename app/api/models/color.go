package models

type Color struct {
	CLRID uint   `json:"clrid"`
	Name  string `json:"color_name" faker:"word"`
}

var AvailableColors = []Color{
	{
		CLRID: 1,
		Name:  "green",
	},
	{
		CLRID: 2,
		Name:  "yellow",
	},
	{
		CLRID: 3,
		Name:  "orange",
	},
	{
		CLRID: 4,
		Name:  "red",
	},
	{
		CLRID: 5,
		Name:  "purple",
	},
	{
		CLRID: 6,
		Name:  "darkblue",
	},
	{
		CLRID: 7,
		Name:  "blue",
	},
	{
		CLRID: 8,
		Name:  "pink",
	},
	{
		CLRID: 9,
		Name:  "darkolivegreen",
	},
	{
		CLRID: 10,
		Name:  "grey",
	},
}
