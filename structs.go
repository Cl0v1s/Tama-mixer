package main

import (
	"encoding/xml"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
)

type BodypartType string

const (
	BodypartType_Leg1  BodypartType = "leg1"
	BodypartType_Leg2               = "leg2"
	BodypartType_Mouth              = "mouth"
	BodypartType_Eye                = "eye"
	BodypartType_Arm1               = "arm1"
	BodypartType_Arm2               = "arm2"
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
	Type      BodypartType
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

var PointsOrder = map[string]int{
	"eye":   0,
	"mouth": 1,
	"arm1":  2,
	"arm2":  3,
	"leg1":  4,
	"leg2":  5,
	"leg3":  6,
}

type Body struct {
	Svg    SVG
	Points []Point
	Parts  []BodyPart
}

func BodyCopy(body Body) Body {
	points := make([]Point, len(body.Points))
	copy(points, body.Points)
	parts := make([]BodyPart, len(body.Parts))
	copy(parts, body.Parts)
	svg := body.Svg
	return Body{
		Svg:    svg,
		Points: points,
		Parts:  parts,
	}
}

func BodyAssemble(body Body) Body {
	svg := SVGCopy(body.Svg)
	label := ""
	for _, point := range body.Points {
		index := slices.IndexFunc(body.Parts, func(part BodyPart) bool {
			return part.Type == point.Type
		})
		if index == -1 {
			continue
		}
		angle := 0.0
		location := point
		t, bezier := findClosestPointInPaths(GetPathsInSVG(svg), point, 2)
		if t >= 0 {
			angle = GetRotationFromBezier(bezier, t)
			location = GetPointFromBezier(bezier, t)
		}
		for _, group := range body.Parts[index].Svg.Groups {
			group = GroupApplyTransformation(group, Transformation{Rotation: angle})
			group = GroupApplyTransformation(group, Transformation{Translation: location})
			svg.Groups = append(svg.Groups, group)
		}
		label += string(point.Type) + "=" + body.Parts[index].Svg.XMLName.Local + "+"
	}
	svg.XMLName.Space = svg.XMLName.Local
	svg.XMLName.Local = label
	body.Svg = svg
	// we prevent further assemble
	body.Points = make([]Point, 0)
	body.Parts = make([]BodyPart, 0)
	return body
}

func BodyIsCompatible(body Body, tpe BodypartType) bool {
	index := slices.IndexFunc(body.Points, func(p Point) bool { return p.Type == tpe })
	if index == -1 {
		return false
	}
	return true
}

func BodyGetBodypart(body Body, tpe BodypartType) (error, BodyPart) {
	index := slices.IndexFunc(body.Parts, func(p BodyPart) bool { return p.Type == tpe })
	if index == -1 {
		return fmt.Errorf("Body does not have a part named %s", tpe), BodyPart{}
	}
	return nil, body.Parts[index]
}

func BodyGetMissingPart(body Body) (error, Point) {
	index := slices.IndexFunc(body.Points, func(p Point) bool {
		err, _ := BodyGetBodypart(body, p.Type)
		if err != nil {
			return true
		}
		return false
	})
	if index == -1 {
		return fmt.Errorf("Body is complete."), Point{}
	}
	return nil, body.Points[index]
}

type BodyPart struct {
	Svg  SVG
	Type BodypartType
}

type SVG struct {
	XMLName xml.Name `xml:"svg"`
	Width   string   `xml:"width,attr"`
	Height  string   `xml:"height,attr"`
	ViewBox string   `xml:"viewBox,attr"`
	Xmlns   string   `xml:"xmlns,attr"`

	Groups []Group `xml:"g"`
}

func SVGCopy(svg SVG) SVG {
	s := svg
	groups := make([]Group, 0)
	for _, g := range svg.Groups {
		groups = append(groups, GroupCopy(g))
	}
	s.Groups = groups
	return s
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

func GroupCopy(group Group) Group {
	g := group
	groups := make([]Group, 0)
	for _, gg := range group.Groups {
		groups = append(groups, GroupCopy(gg))
	}
	paths := make([]Path, len(group.Paths))
	copy(paths, group.Paths)
	g.Paths = paths
	ellipses := make([]Ellipse, len(group.Ellipses))
	copy(ellipses, group.Ellipses)
	g.Ellipses = ellipses
	circles := make([]Circle, len(group.Circles))
	copy(circles, group.Circles)
	g.Circles = circles

	return g
}

func GroupApplyTransformation(group Group, t Transformation) Group {
	paths := GetPathsInGroup(group)
	results := make([]Path, 0)
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
		p := Path{
			ID:        paths[i].ID,
			Label:     paths[i].Label,
			Style:     paths[i].Style,
			Transform: paths[i].Transform,
			D:         bezierToD(beziers),
		}
		results = append(results, p)
	}
	group.Paths = results
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
