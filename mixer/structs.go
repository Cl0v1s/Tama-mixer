package main

import (
	"encoding/json"
	"encoding/xml"
	"math"
	"os"
	"strconv"
	"strings"
)

type BodypartType string

const (
	BodypartType_Leg1  BodypartType = "leg1"
	BodypartType_Leg2  BodypartType = "leg2"
	BodypartType_Mouth BodypartType = "mouth"
	BodypartType_Eye   BodypartType = "eye"
	BodypartType_Arm1  BodypartType = "arm1"
	BodypartType_Arm2  BodypartType = "arm2"
)

type Bezier struct {
	P0 Point
	P1 Point
	P2 Point
	P3 Point
}

func BeziersToD(beziers []Bezier) string {
	result := ""
	for _, b := range beziers {
		result += " M " + strconv.FormatFloat(b.P0.X, 'f', -1, 64) + " " + strconv.FormatFloat(b.P0.Y, 'f', -1, 64)
		points := []Point{b.P1, b.P2, b.P3}
		result += " C "
		for _, p := range points {
			result += strconv.FormatFloat(p.X, 'f', -1, 64) + " " + strconv.FormatFloat(p.Y, 'f', -1, 64) + " "
		}
	}
	return result
}

type Point struct {
	X    float64      `json:"x"`
	Y    float64      `json:"y"`
	T    float64      `json:"t"`
	Type BodypartType `json:"type"`
}

func (p Point) Quadrant() int {
	if p.X <= 0 && p.Y > 0 {
		// quadrant = 1 // 1 * pi
		return 1
	} else if p.X <= 0 && p.Y < 0 {
		// quadrant = 2 // 0 * pi
		return 2
	} else if p.X > 0 && p.Y <= 0 {
		// quadrant = 3 // 0 * pi
		return 3
	} else {
		// quadrant = 4 // 1 * pi
		return 4
	}
}

func (p Point) Add(p1 Point) Point {
	return Point{
		X: p.X + p1.X,
		Y: p.Y + p1.Y,
	}
}

func (p Point) Sub(p1 Point) Point {
	return p.Add(Point{X: p1.X * -1, Y: p1.Y * -1})
}

func (p1 Point) Distance(p2 Point) float64 {
	return math.Abs(p1.X-p2.X) + math.Abs(p1.Y-p2.Y)
}

func (p Point) Rotate(angleDeg float64) Point {
	angle := angleDeg * (math.Pi / 180)
	return Point{
		X: p.X*math.Cos(angle) - p.Y*math.Sin(angle),
		Y: p.X*math.Sin(angle) + p.Y*math.Cos(angle),
	}
}

func (p Point) Translate(t Point) Point {
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
	Path   string  `json:"path"`
	Points []Point `json:"points"`
	Frame  int     `json:"frame"`
	Name   string  `json:"name"`
	Size   Point   `json:"size"`
}

func SaveBodyPartsToJSON(prefix string, bodyparts []BodyPart) error {
	_ = os.MkdirAll(prefix, 0755)
	filename := string(bodyparts[0].Type) + "-" + bodyparts[0].Name + ".json"
	file, err := os.Create(prefix + "/" + filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(bodyparts)
}

func SaveBodiesToJSON(prefix string, bodies []Body) error {
	_ = os.MkdirAll(prefix, 0755)
	filename := bodies[0].Name + ".json"
	file, err := os.Create(prefix + "/" + filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(bodies)
}

type BodyPart struct {
	Path  string       `json:"path"`
	Type  BodypartType `json:"type"`
	Frame int          `json:"frame"`
	Name  string       `json:"name"`
}

type SVG struct {
	XMLName xml.Name `xml:"svg"`
	Width   string   `xml:"width,attr"`
	Height  string   `xml:"height,attr"`
	ViewBox string   `xml:"viewBox,attr"`
	Xmlns   string   `xml:"xmlns,attr"`

	Groups []Group `xml:"g"`
}

func (s SVG) String() string {
	data, err := xml.MarshalIndent(s, "", "  ")
	if err != nil {
		return "" // ou gérer l'erreur selon votre besoin
	}
	return string(data)
}

type Transformation struct {
	Rotation    float64
	Translation Point
}

type Group struct {
	ID    string `xml:"id,attr"`
	Label string `xml:"label,attr"`

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

func (group Group) GetPath() Path {
	paths := GetPathsInGroup(group)
	resultCmds := make([]string, 0)
	for _, path := range paths {
		cmds := ParseD(path.D)
		if cmds[0].Type != "M" {
			x := 0.0
			y := 0.0
			if len(cmds[0].Args) >= 1 {
				x += cmds[0].Args[0]
			}
			if len(cmds[0].Args) >= 2 {
				y += cmds[0].Args[1]
			}
			resultCmds = append(resultCmds, "M "+strconv.FormatFloat(x, 'f', 2, 64)+","+strconv.FormatFloat(y, 'f', 2, 64))
		}
		for _, cmd := range cmds {
			args := make([]string, 0)
			for _, a := range cmd.Args {
				args = append(args, strconv.FormatFloat(a, 'f', 2, 64))
			}
			resultCmds = append(resultCmds, cmd.Type+" "+strings.Join(args, ","))
		}
	}
	return Path{
		D: strings.Join(resultCmds, " "),
	}
}

// Apply transformations to a group, for then, all coords are absolute
func (group Group) Transform(t Transformation) Group {
	result := GroupCopy(group)

	ellipsis := make([]Ellipse, 0)
	for _, el := range group.Ellipses {
		p := Point{X: el.CX, Y: el.CY}
		p = p.Translate(t.Translation)
		p = p.Rotate(t.Rotation)
		el.CX = p.X
		el.CY = p.Y
		ellipsis = append(ellipsis, el)
	}
	result.Ellipses = ellipsis

	circles := make([]Circle, 0)
	for _, ci := range group.Circles {
		p := Point{X: ci.CX, Y: ci.CY}
		p = p.Translate(t.Translation)
		p = p.Rotate(t.Rotation)
		ci.CX = p.X
		ci.CY = p.Y
		circles = append(circles, ci)
	}
	result.Circles = circles

	paths := GetPathsInGroup(group)

	groups := make([]Group, 0)
	for _, g := range groups {
		// we remove path as they are alreay retrieve by GetPathsInGroup
		g.Paths = []Path{}
		groups = append(groups, g.Transform(t))
	}
	result.Groups = groups

	finalPaths := make([]Path, 0)
	for i := 0; i < len(paths); i++ {
		commands := ParseD(paths[i].D)
		beziers := GetBeziersFromCommands(commands)
		for u := 0; u < len(beziers); u++ {
			beziers[u].P0 = beziers[u].P0.Translate(t.Translation)
			beziers[u].P1 = beziers[u].P1.Translate(t.Translation)
			beziers[u].P2 = beziers[u].P2.Translate(t.Translation)
			beziers[u].P3 = beziers[u].P3.Translate(t.Translation)

			beziers[u].P0 = beziers[u].P0.Rotate(t.Rotation)
			beziers[u].P1 = beziers[u].P1.Rotate(t.Rotation)
			beziers[u].P2 = beziers[u].P2.Rotate(t.Rotation)
			beziers[u].P3 = beziers[u].P3.Rotate(t.Rotation)
		}
		p := Path{
			ID:    paths[i].ID,
			Label: paths[i].Label,
			Style: paths[i].Style,
			D:     BeziersToD(beziers),
		}
		finalPaths = append(finalPaths, p)
	}
	result.Paths = finalPaths

	return result
}

type Ellipse struct {
	ID    string  `xml:"id,attr"`
	Label string  `xml:"label,attr"`
	CX    float64 `xml:"cx,attr"`
	CY    float64 `xml:"cy,attr"`
	RX    float64 `xml:"rx,attr"`
	RY    float64 `xml:"ry,attr"`
}

type Circle struct {
	ID    string  `xml:"id,attr"`
	Label string  `xml:"label,attr"`
	CX    float64 `xml:"cx,attr"`
	CY    float64 `xml:"cy,attr"`
	R     float64 `xml:"r,attr"`
}

type Command struct {
	Type string
	Args []float64
}

type Path struct {
	ID    string `xml:"id,attr"`
	Label string `xml:"label,attr"`
	D     string `xml:"d,attr"`
	Style string `xml:"style,attr"`
}
