package main

import (
	"encoding/xml"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
	"unicode"
)

func findClosestPointInPaths(paths []Path, point Point, rng float64) (float64, Bezier) {
	for _, p := range paths {
		commands := ParseD(p.D)
		beziers := GetBeziersFromCommands(commands)
		for _, b := range beziers {
			// TODO: peut largement être optimisé
			for i := 0.0; i <= 1; i += 0.1 {
				p := GetPointFromBezier(b, i)
				d := math.Abs(p.X-point.X) + math.Abs(p.Y-point.Y)
				if d <= rng {
					return math.Round(i*100) / 100, b
				}
			}
		}
	}
	return -1, Bezier{}
}

func ParseD(d string) []Command {
	var currentCmd *Command = nil
	commands := make([]Command, 0)
	buffer := make([]rune, 0)
	lastRune := ' '
	for i := 0; i < len(d); i++ {
		c := rune(d[i])

		if (unicode.IsLetter(c) && lastRune == ' ') || c == ',' || c == ' ' {
			if currentCmd != nil && len(buffer) > 0 {
				number, err := strconv.ParseFloat(string(buffer), 64)
				if err != nil {
					panic(err)
				}
				currentCmd.Args = append(currentCmd.Args, number)
			}
			buffer = make([]rune, 0)
		} else {
			buffer = append(buffer, c)
		}

		if unicode.IsLetter(c) && lastRune == ' ' {
			if currentCmd != nil {
				commands = append(commands, *currentCmd)
			}
			currentCmd = &Command{
				Type: string(c),
				Args: make([]float64, 0),
			}
		}
		lastRune = c
	}

	if currentCmd != nil && len(buffer) > 0 {
		number, err := strconv.ParseFloat(string(buffer), 64)
		if err != nil {
			panic(err)
		}
		currentCmd.Args = append(currentCmd.Args, number)
	}
	if currentCmd != nil {
		commands = append(commands, *currentCmd)
	}

	return commands
}

func Save(basepath string, svg SVG) {
	_ = os.MkdirAll(basepath, 0755)
	path := basepath + "/" + svg.XMLName.Space + "@" + svg.XMLName.Local + ".svg"
	out, _ := xml.MarshalIndent(svg, " ", "  ")
	os.WriteFile(path, out, 0644)
}

func GetPathsInGroup(group Group) []Path {
	paths := group.Paths
	for _, g := range group.Groups {
		paths = append(paths, GetPathsInGroup(g)...)
	}
	return paths
}

func GetPathsInSVG(svg SVG) []Path {
	paths := make([]Path, 0)
	for _, g := range svg.Groups {
		paths = append(paths, GetPathsInGroup(g)...)
	}
	return paths
}

func GetPointFromBezier(bezier Bezier, t float64) Point {
	x := math.Pow(1-t, 3)*bezier.P0.X + 3*math.Pow(1-t, 2)*t*bezier.P1.X + 3*(1-t)*t*t*bezier.P2.X + math.Pow(t, 3)*bezier.P3.X
	y := math.Pow(1-t, 3)*bezier.P0.Y + 3*math.Pow(1-t, 2)*t*bezier.P1.Y + 3*(1-t)*t*t*bezier.P2.Y + math.Pow(t, 3)*bezier.P3.Y
	return Point{X: math.Round(x*100) / 100, Y: math.Round(y*100) / 100}
}

func GetRotationFromBezier(bezier Bezier, t float64) float64 {
	dx := 3*(1-t)*(1-t)*(bezier.P1.X-bezier.P0.X) + 6*(1-t)*t*(bezier.P2.X-bezier.P1.X) + 3*t*t*(bezier.P3.X-bezier.P2.X)
	dy := 3*(1-t)*(1-t)*(bezier.P1.Y-bezier.P0.Y) + 6*(1-t)*t*(bezier.P2.Y-bezier.P1.Y) + 3*t*t*(bezier.P3.Y-bezier.P2.Y)
	return math.Atan2(dy, dx) * 180 / math.Pi
}

func GetBeziersFromCommands(commands []Command) []Bezier {
	results := make([]Bezier, 0)
	current := Point{X: 0, Y: 0}
	var zPoint *Point = nil
	for _, command := range commands {
		if command.Type == "M" {
			current.X = command.Args[0]
			current.Y = command.Args[1]
			rest := command.Args[2:]
			for i := 0; i < len(rest); i += 2 {
				b := Bezier{}
				b.P0 = current
				b.P1 = current
				b.P2 = Point{X: rest[i], Y: rest[i+1]}
				b.P3 = b.P2
				current = b.P3
				results = append(results, b)
			}
		} else if command.Type == "m" {
			current.X += command.Args[0]
			current.Y += command.Args[1]
			rest := command.Args[2:]
			for i := 0; i < len(rest); i += 2 {
				b := Bezier{}
				b.P0 = current
				b.P1 = current
				b.P2 = Point{X: current.X + rest[i], Y: current.Y + rest[i+1]}
				b.P3 = b.P2
				current = b.P3
				results = append(results, b)
			}
		} else if command.Type == "c" {
			for i := 0; i < len(command.Args); i += 6 {
				b := Bezier{}
				b.P0 = current
				b.P1 = Point{X: current.X + command.Args[i+0], Y: current.Y + command.Args[i+1]}
				b.P2 = Point{X: current.X + command.Args[i+2], Y: current.Y + command.Args[i+3]}
				b.P3 = Point{X: current.X + command.Args[i+4], Y: current.Y + command.Args[i+5]}
				current = b.P3
				results = append(results, b)
			}
		} else if command.Type == "C" {
			for i := 0; i < len(command.Args); i += 6 {
				b := Bezier{}
				b.P0 = current
				b.P1 = Point{X: command.Args[i+0], Y: command.Args[i+1]}
				b.P2 = Point{X: command.Args[i+2], Y: command.Args[i+3]}
				b.P3 = Point{X: command.Args[i+4], Y: command.Args[i+5]}
				current = b.P3
				results = append(results, b)
			}
		} else if command.Type == "L" {
			b := Bezier{}
			b.P0 = current
			b.P1 = current
			b.P2 = Point{X: command.Args[0], Y: command.Args[1]}
			b.P3 = Point{X: command.Args[0], Y: command.Args[1]}
			results = append(results, b)
			current = b.P3
		} else if command.Type == "l" {
			b := Bezier{}
			b.P0 = current
			b.P1 = current
			b.P2 = Point{X: command.Args[0] + current.X, Y: command.Args[1] + current.Y}
			b.P3 = Point{X: command.Args[0] + current.X, Y: command.Args[1] + current.Y}
			results = append(results, b)
			current = b.P3
		} else if command.Type == "Z" || command.Type == "z" {
			b := Bezier{}
			b.P0 = current
			b.P1 = current
			b.P2 = *zPoint
			b.P3 = *zPoint
			results = append(results, b)
			current = b.P3
			zPoint = nil
		} else {
			panic("GetBeziersFromCommands: Unsupported command " + command.Type)
		}
		if zPoint == nil {
			zPoint = &Point{X: current.X, Y: current.Y}
		}
	}
	return results
}

// Search for elements named as the group and calculate their lower coordinates
func FindLowestPadding(g Group, anchorLabel string) (float64, float64) {
	x, y := math.MaxFloat64, math.MaxFloat64
	for _, g := range g.Groups {
		x1, y1 := FindLowestPadding(g, anchorLabel)
		if x1 < x {
			x = x1
		}
		if y1 < y {
			y = y1
		}
	}

	for _, p := range g.Paths {
		if anchorLabel != "" && p.Label != anchorLabel {
			continue
		}
		cmds := ParseD(p.D)
		beziers := GetBeziersFromCommands(cmds)
		for _, bezier := range beziers {
			points := []Point{bezier.P0, bezier.P1, bezier.P2, bezier.P3}
			for _, point := range points {
				x1, y1 := point.X, point.Y
				if x1 < x {
					x = x1
				}
				if y1 < y {
					y = y1
				}
			}
		}
	}

	for _, e := range g.Ellipses {
		if anchorLabel != "" && e.Label != anchorLabel {
			continue
		}
		x1, y1 := e.CX, e.CY
		if x1 < x {
			x = x1
		}
		if y1 < y {
			y = y1
		}
	}

	for _, e := range g.Circles {
		if anchorLabel != "" && e.Label != anchorLabel {
			continue
		}
		x1, y1 := e.CX, e.CY
		if x1 < x {
			x = x1
		}
		if y1 < y {
			y = y1
		}
	}

	return x, y
}

func Unpad(g Group, x float64, y float64) Group {
	ellipsis := make([]Ellipse, 0)
	for _, el := range g.Ellipses {
		el.CX -= x
		el.CY -= y
		ellipsis = append(ellipsis, el)
	}
	circles := make([]Circle, 0)
	for _, ci := range g.Circles {
		ci.CX -= x
		ci.CY -= y
		circles = append(circles, ci)
	}
	paths := make([]Path, 0)
	for _, p := range g.Paths {
		cmds := ParseD(p.D)
		d := ""
		for i := 0; i < len(cmds); i++ {
			if strings.ToLower(cmds[i].Type) == "m" {
				for u := 0; u < len(cmds[i].Args); u += 2 {
					cmds[i].Args[u] -= x
					cmds[i].Args[u+1] -= y
				}
			}
			d = d + " " + MarshalD(cmds[i])
		}
		p.D = d
		paths = append(paths, p)
	}
	groups := make([]Group, 0)
	for _, g := range g.Groups {
		g := Unpad(g, x, y)
		groups = append(groups, g)
	}
	group := Group{
		ID:       g.ID,
		Label:    g.Label,
		Paths:    paths,
		Groups:   groups,
		Ellipses: ellipsis,
		Circles:  circles,
	}
	return group
}

func RetrievePoints(group *Group, rootLabel string) []Point {
	points := make([]Point, 0)
	for _, g := range group.Groups {
		points = append(points, RetrievePoints(&g, rootLabel)...)
	}
	i := 0
	for _, e := range group.Ellipses {
		if e.Label == rootLabel {
			group.Ellipses[i] = e
			i++
		} else {
			points = append(points, Point{X: e.CX, Y: e.CY, Type: BodypartType(e.Label), Transform: e.Transform})
		}
	}
	group.Ellipses = group.Ellipses[:i]
	i = 0
	for _, e := range group.Circles {
		if e.Label == rootLabel {
			group.Circles[i] = e
			i++
		} else {
			points = append(points, Point{X: e.CX, Y: e.CY, Type: BodypartType(e.Label), Transform: e.Transform})
		}
	}
	group.Circles = group.Circles[:i]
	i = 0
	for _, p := range group.Paths {
		if p.Label == rootLabel {
			group.Paths[i] = p
			i++
		} else {
			commands := ParseD(p.D)
			index := slices.IndexFunc(commands, func(c Command) bool {
				return strings.ToLower(c.Type) == "m"
			})
			if index == -1 {
				continue
			}
			points = append(points, Point{X: commands[index].Args[0], Y: commands[index].Args[1], Type: BodypartType(p.Label), Transform: p.Transform})
		}
	}
	group.Paths = group.Paths[:i]
	slices.SortFunc(points, func(a Point, b Point) int {
		return PointsOrder[string(a.Type)] - PointsOrder[string(b.Type)]
	})
	return points
}

func CleanGroup(group *Group) {
	for i := 0; i < len(group.Groups); i++ {
		CleanGroup(&group.Groups[i])
	}
	group.Circles = slices.DeleteFunc(group.Circles, func(c Circle) bool { return c.Label == group.Label })
	group.Ellipses = slices.DeleteFunc(group.Ellipses, func(c Ellipse) bool { return c.Label == group.Label })
	group.Paths = slices.DeleteFunc(group.Paths, func(c Path) bool { return c.Label == group.Label })
}

func parseBody(name string, group Group) Body {
	x, y := FindLowestPadding(group, "body")
	group = Unpad(group, x, y)
	anchors := RetrievePoints(&group, group.Label)
	svg := SVG{
		XMLName: xml.Name{
			Space: group.Label,
			Local: name,
		},
		Xmlns:   "http://www.w3.org/2000/svg",
		Width:   "75",
		Height:  "75",
		ViewBox: "-25 -25 75 75",
		Groups:  []Group{group},
	}
	return Body{
		Points: anchors,
		Svg:    svg,
	}
}

func parseBodypart(body Body, group Group) BodyPart {
	x, y := FindLowestPadding(group, group.Label)
	group = Unpad(group, x, y)
	anchorIndex := slices.IndexFunc(body.Points, func(p Point) bool { return p.Type == BodypartType(group.Label) })
	anchor := body.Points[anchorIndex]
	paths := GetPathsInSVG(body.Svg)
	t, b := findClosestPointInPaths(paths, anchor, 2)
	if t >= 0 {
		angle := GetRotationFromBezier(b, t) * -1
		group = GroupApplyTransformation(group, Transformation{Rotation: angle})
	}
	CleanGroup(&group)
	svg := SVG{
		XMLName: xml.Name{
			Space: group.Label,
			Local: group.ID,
		},
		Xmlns:   "http://www.w3.org/2000/svg",
		Width:   "100",
		Height:  "100",
		ViewBox: "-50 -50 100 100",
		Groups:  []Group{group},
	}
	return BodyPart{
		Svg:  svg,
		Type: BodypartType(group.Label),
	}
}

func Sort(root SVG) ([]Body, []BodyPart) {
	bodies := make([]Body, 0)
	bodyparts := make([]BodyPart, 0)

	for _, character := range root.Groups {
		bodyIndex := slices.IndexFunc(character.Groups, func(g Group) bool { return g.Label == "body" })
		body := parseBody(character.Label, character.Groups[bodyIndex])
		for _, group := range character.Groups {
			if group.Label != "body" {
				bodyparts = append(bodyparts, parseBodypart(body, group))
			}
		}
		bodies = append(bodies, body)
	}

	return bodies, bodyparts
}
