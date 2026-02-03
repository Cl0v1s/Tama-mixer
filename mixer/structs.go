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
	Path        string       `json:"path"`
	Type        BodypartType `json:"type"`
	Frame       int          `json:"frame"`
	Name        string       `json:"name"`
	BoundingBox Rect         `json:"boundingBox"`
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
		resultCmds = append(resultCmds, path.D)
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
		for u := 0; u < len(commands); u++ {
			commands[u].Transform(t)
		}
		finalPaths = append(finalPaths, Path{D: CompileD(commands)})
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

func (cmd *Command) Transform(t Transformation) {
	switch cmd.Type {
	case "M", "L", "C":
		for u := 0; u < len(cmd.Args); u += 2 {
			point := Point{X: cmd.Args[u], Y: cmd.Args[u+1]}
			point = point.Translate(t.Translation)
			point = point.Rotate(t.Rotation)
			cmd.Args[u] = point.X
			cmd.Args[u+1] = point.Y
		}
	case "A":
		for u := 0; u < len(cmd.Args); u += 7 {
			point := Point{X: cmd.Args[u+5], Y: cmd.Args[u+6]}
			point = point.Translate(t.Translation)
			point = point.Rotate(t.Rotation)
			cmd.Args[u+5] = point.X
			cmd.Args[u+6] = point.Y
		}
	case "V":
		point := Point{X: 0, Y: cmd.Args[0]}
		point = point.Translate(t.Translation)
		point = point.Rotate(t.Rotation)
		cmd.Args[0] = point.Y
	case "H":
		point := Point{X: cmd.Args[0], Y: 0}
		point = point.Translate(t.Translation)
		point = point.Rotate(t.Rotation)
		cmd.Args[0] = point.X
	}
}

type Rect struct {
	TopLeft     Point `json:"topLeft"`
	BottomRight Point `json:"bottomRight"`
}

type Path struct {
	ID    string `xml:"id,attr"`
	Label string `xml:"label,attr"`
	D     string `xml:"d,attr"`
	Style string `xml:"style,attr"`
}

func (path *Path) GetBoundingBox() Rect {
	cmds := ParseD(path.D)
	bzs := GetBeziersFromCommands(cmds)
	lowest := Point{}
	highwest := Point{}
	for _, bz := range bzs {
		for i := 0.0; i <= 1.0; i += 0.05 {
			point := GetPointFromBezier(bz, i)
			if point.X > highwest.X {
				highwest.X = point.X
			}
			if point.Y > highwest.Y {
				highwest.Y = point.Y
			}
			if point.X < lowest.X {
				lowest.X = point.X
			}
			if point.Y < lowest.Y {
				lowest.Y = point.Y
			}
		}
	}
	return Rect{TopLeft: lowest, BottomRight: highwest}
}
