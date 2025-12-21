package main

import (
	"encoding/xml"
	"fmt"
	"math"
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

func PointDistance(p1 Point, p2 Point) float64 {
	return math.Abs(p1.X-p2.X) + math.Abs(p1.Y-p2.Y)
}

func PointRotate(p Point, angleDeg float64) Point {
	angle := angleDeg * (math.Pi / 180)
	return Point{
		X: p.X*math.Cos(angle) - p.Y*math.Sin(angle),
		Y: p.X*math.Sin(angle) + p.Y*math.Cos(angle),
	}
}

func PointTranslate(p Point, t Point) Point {
	return Point{
		X: p.X + t.X,
		Y: p.Y + t.Y,
	}
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

type Transformation struct {
	Rotation    float64
	Translation Point
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

func GroupApplyTransformation(group Group, t Transformation) Group {
	paths := GetPathsInGroup(group)
	for i := 0; i < len(paths); i++ {
		commands := ParseD(paths[i].D)
		beziers := GetBeziersFromCommands(commands)
		for u := 0; u < len(beziers); u++ {
			beziers[u].P0 = PointTranslate(beziers[u].P0, t.Translation)
			beziers[u].P1 = PointTranslate(beziers[u].P1, t.Translation)
			beziers[u].P2 = PointTranslate(beziers[u].P2, t.Translation)
			beziers[u].P3 = PointTranslate(beziers[u].P3, t.Translation)

			beziers[u].P0 = PointRotate(beziers[u].P0, t.Rotation)
			beziers[u].P1 = PointRotate(beziers[u].P1, t.Rotation)
			beziers[u].P2 = PointRotate(beziers[u].P2, t.Rotation)
			beziers[u].P3 = PointRotate(beziers[u].P3, t.Rotation)
		}
		paths[i].D = bezierToD(beziers)
	}
	group.Paths = paths
	return group
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
