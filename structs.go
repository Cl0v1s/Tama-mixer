package main

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type Bezier struct {
	P0 Point
	P1 Point
	P2 Point
	P3 Point
}

type Point struct {
	X         float64
	Y         float64
	Label     string
	Transform string
}

type Body struct {
	Svg    SVG
	Points []Point
}

type BodyPart struct {
	Svg   SVG
	Label string
}

type SVG struct {
	XMLName xml.Name `xml:"svg"`
	Width   string   `xml:"width,attr"`
	Height  string   `xml:"height,attr"`
	ViewBox string   `xml:"viewBox,attr"`
	Xmlns   string   `xml:"xmlns,attr"`

	Groups []Group `xml:"g"`
}

type Group struct {
	ID    string `xml:"id,attr"`
	Label string `xml:"label,attr"`

	Transform string `xml:"transform,attr"`

	Groups   []Group   `xml:"g"`
	Paths    []Path    `xml:"path"`
	Ellipses []Ellipse `xml:"ellipse"`
	Circles  []Circle  `xml:"circle"`
}

type Ellipse struct {
	ID        string  `xml:"id,attr"`
	Label     string  `xml:"label,attr"`
	CX        float64 `xml:"cx,attr"`
	CY        float64 `xml:"cy,attr"`
	RX        float64 `xml:"rx,attr"`
	RY        float64 `xml:"ry,attr"`
	Transform string  `xml:"transform,attr"`
}

type Circle struct {
	ID        string  `xml:"id,attr"`
	Label     string  `xml:"label,attr"`
	CX        float64 `xml:"cx,attr"`
	CY        float64 `xml:"cy,attr"`
	R         float64 `xml:"r,attr"`
	Transform string  `xml:"transform,attr"`
}

type Command struct {
	Type string
	Args []float64
}

type Path struct {
	ID        string `xml:"id,attr"`
	Label     string `xml:"label,attr"`
	D         string `xml:"d,attr"`
	Style     string `xml:"style,attr"`
	Transform string `xml:"transform,attr"`
}

func PrintCommands(cmds []Command) {
	for i, cmd := range cmds {
		fmt.Printf(
			"%02d | %s %v\n",
			i,
			cmd.Type,
			cmd.Args,
		)
	}
}

func MarshalD(cmd Command) string {
	var b strings.Builder

	// Lettre de commande
	b.WriteString(cmd.Type)

	// Arguments
	for _, arg := range cmd.Args {
		b.WriteByte(' ')
		b.WriteString(strconv.FormatFloat(arg, 'f', -1, 64))
	}

	return b.String()
}
