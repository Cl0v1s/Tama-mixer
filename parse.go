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

func ArcToBeziers(current Point, rx, ry, xAxisRotation float64, largeArcFlag, sweepFlag int, x, y float64) []Bezier {
	// Convertir rotation en radians
	phi := xAxisRotation * math.Pi / 180.0

	// Point final
	P1 := Point{X: x, Y: y}

	// Déplacement si start == end
	if current.X == P1.X && current.Y == P1.Y {
		return nil
	}

	// Correction des rayons si nécessaire
	rx = math.Abs(rx)
	ry = math.Abs(ry)
	if rx == 0 || ry == 0 {
		return []Bezier{{P0: current, P1: current, P2: P1, P3: P1}}
	}

	// Convertir coordonnées dans le repère de l'ellipse
	dx2 := (current.X - P1.X) / 2.0
	dy2 := (current.Y - P1.Y) / 2.0
	x1p := math.Cos(phi)*dx2 + math.Sin(phi)*dy2
	y1p := -math.Sin(phi)*dx2 + math.Cos(phi)*dy2

	// Calcul du centre cx', cy'
	rx2 := rx * rx
	ry2 := ry * ry
	x1p2 := x1p * x1p
	y1p2 := y1p * y1p

	var sign float64 = 1
	if largeArcFlag == sweepFlag {
		sign = -1
	}

	sq := ((rx2*ry2 - rx2*y1p2 - ry2*x1p2) / (rx2*y1p2 + ry2*x1p2))
	if sq < 0 {
		sq = 0
	}
	coef := sign * math.Sqrt(sq)
	cxp := coef * (rx * y1p / ry)
	cyp := coef * -(ry * x1p / rx)

	// Centre dans le repère original
	cx := math.Cos(phi)*cxp - math.Sin(phi)*cyp + (current.X+P1.X)/2
	cy := math.Sin(phi)*cxp + math.Cos(phi)*cyp + (current.Y+P1.Y)/2

	// Angles de début et de fin
	theta1 := math.Atan2((y1p-cyp)/ry, (x1p-cxp)/rx)
	deltaTheta := math.Atan2((-y1p-cyp)/ry, (-x1p-cxp)/rx) - theta1

	if sweepFlag == 0 && deltaTheta > 0 {
		deltaTheta -= 2 * math.Pi
	} else if sweepFlag == 1 && deltaTheta < 0 {
		deltaTheta += 2 * math.Pi
	}

	// Diviser en segments ≤ 90°
	segments := int(math.Ceil(math.Abs(deltaTheta) / (math.Pi / 2)))
	delta := deltaTheta / float64(segments)

	var beziers []Bezier
	for i := 0; i < segments; i++ {
		t1 := theta1 + float64(i)*delta
		t2 := t1 + delta

		// Points de Bézier pour ce segment
		alpha := math.Sin(delta) * (math.Sqrt(4+3*math.Pow(math.Tan(delta/2), 2)) - 1) / 3

		p0 := Point{
			X: cx + rx*math.Cos(phi)*math.Cos(t1) - ry*math.Sin(phi)*math.Sin(t1),
			Y: cy + rx*math.Sin(phi)*math.Cos(t1) + ry*math.Cos(phi)*math.Sin(t1),
		}
		p3 := Point{
			X: cx + rx*math.Cos(phi)*math.Cos(t2) - ry*math.Sin(phi)*math.Sin(t2),
			Y: cy + rx*math.Sin(phi)*math.Cos(t2) + ry*math.Cos(phi)*math.Sin(t2),
		}
		dx := p3.X - p0.X
		dy := p3.Y - p0.Y
		p1 := Point{X: p0.X + alpha*dx, Y: p0.Y + alpha*dy}
		p2 := Point{X: p3.X - alpha*dx, Y: p3.Y - alpha*dy}

		beziers = append(beziers, Bezier{P0: p0, P1: p1, P2: p2, P3: p3})
	}

	return beziers
}

// Fonction principale
func GetBeziersFromCommands(commands []Command) []Bezier {
	results := make([]Bezier, 0)
	current := Point{X: 0, Y: 0}
	var zPoint *Point = nil

	for _, command := range commands {
		switch command.Type {
		case "M":
			current.X = command.Args[0]
			current.Y = command.Args[1]
			if zPoint == nil {
				zPoint = &Point{X: current.X, Y: current.Y}
			}
		case "m":
			current.X += command.Args[0]
			current.Y += command.Args[1]
			if zPoint == nil {
				zPoint = &Point{X: current.X, Y: current.Y}
			}
		case "L", "l":
			x, y := command.Args[0], command.Args[1]
			if command.Type == "l" {
				x += current.X
				y += current.Y
			}
			results = append(results, Bezier{P0: current, P1: current, P2: Point{X: x, Y: y}, P3: Point{X: x, Y: y}})
			current = Point{X: x, Y: y}
		case "C", "c":
			for i := 0; i < len(command.Args); i += 6 {
				var b Bezier
				b.P0 = current
				if command.Type == "c" {
					b.P1 = Point{X: current.X + command.Args[i], Y: current.Y + command.Args[i+1]}
					b.P2 = Point{X: current.X + command.Args[i+2], Y: current.Y + command.Args[i+3]}
					b.P3 = Point{X: current.X + command.Args[i+4], Y: current.Y + command.Args[i+5]}
				} else {
					b.P1 = Point{X: command.Args[i], Y: command.Args[i+1]}
					b.P2 = Point{X: command.Args[i+2], Y: command.Args[i+3]}
					b.P3 = Point{X: command.Args[i+4], Y: command.Args[i+5]}
				}
				current = b.P3
				results = append(results, b)
			}
		case "A", "a":
			for i := 0; i < len(command.Args); i += 7 {
				rx := command.Args[i]
				ry := command.Args[i+1]
				xAxisRotation := command.Args[i+2]
				largeArcFlag := int(command.Args[i+3])
				sweepFlag := int(command.Args[i+4])
				x := command.Args[i+5]
				y := command.Args[i+6]
				if command.Type == "a" {
					x += current.X
					y += current.Y
				}

				beziers := ArcToBeziers(current, rx, ry, xAxisRotation, largeArcFlag, sweepFlag, x, y)
				for _, b := range beziers {
					results = append(results, b)
					current = b.P3
				}
			}
		case "Z", "z":
			if zPoint != nil {
				results = append(results, Bezier{P0: current, P1: current, P2: *zPoint, P3: *zPoint})
				current = *zPoint
			}
			zPoint = nil
		}
	}
	return results
}

func findLowestPadding(g Group) (float64, float64) {
	x, y := math.MaxFloat64, math.MaxFloat64
	paths := GetPathsInGroup(g)
	for _, path := range paths {
		commands := ParseD(path.D)
		beziers := GetBeziersFromCommands(commands)
		for _, b := range beziers {
			for i := 0.0; i <= 1; i += 0.1 {
				p := GetPointFromBezier(b, i)
				if x > p.X {
					x = p.X
				}
				if y > p.Y {
					y = p.Y
				}
			}
		}
	}
	return x, y
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
			points = append(points, Point{X: e.CX, Y: e.CY, Type: BodypartType(e.Label)})
		}
	}
	group.Ellipses = group.Ellipses[:i]
	i = 0
	for _, e := range group.Circles {
		if e.Label == rootLabel {
			group.Circles[i] = e
			i++
		} else {
			points = append(points, Point{X: e.CX, Y: e.CY, Type: BodypartType(e.Label)})
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
			points = append(points, Point{X: commands[index].Args[0], Y: commands[index].Args[1], Type: BodypartType(p.Label)})
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

func parseBody(g Group) Body {
	group := GroupCopy(g)
	group.ID = group.Label
	group.Label = "body"
	x, y := findLowestPadding(group)
	group = GroupApplyTransformation(group, Transformation{Translation: Point{X: -x, Y: -y}})
	anchors := RetrievePoints(&group, group.Label)
	svg := SVG{
		XMLName: xml.Name{
			Space: group.Label,
			Local: group.ID,
		},
		Xmlns:   "http://www.w3.org/2000/svg",
		Width:   "100",
		Height:  "100",
		ViewBox: "0 0 100 100",
		Groups:  []Group{group},
	}
	return Body{
		Points: anchors,
		Svg:    svg,
	}
}

// do not considere paths
func findElementPosition(group Group, name string) Point {
	for _, el := range group.Ellipses {
		if el.Label == name {
			return Point{X: el.CX, Y: el.CY}
		}
	}
	for _, ci := range group.Circles {
		if ci.Label == name {
			return Point{X: ci.CX, Y: ci.CY}
		}
	}
	for _, g := range group.Groups {
		p := findElementPosition(g, name)
		if p.X != math.MaxFloat64 && p.Y != math.MaxFloat64 {
			return p
		}
	}
	return Point{X: math.MaxFloat64, Y: math.MaxFloat64}
}

func PointFindQuadrant(p Point) int {
	switch {
	case p.X < 0 && p.Y < 0:
		return 1 // haut-gauche
	case p.X < 0 && p.Y > 0:
		return 2 // bas-gauche
	case p.X > 0 && p.Y > 0:
		return 3 // bas-droite
	case p.X > 0 && p.Y < 0:
		return 4 // haut-droite
	default:
		return 0 // sur un axe (X == 0 ou Y == 0)
	}
}
func GroupNormalizeRotation(group Group) Group {
	paths := GetPathsInGroup(group)
	tail := Point{X: 0, Y: 0}

	var points []Point
	for _, path := range paths {
		commands := ParseD(path.D)
		beziers := GetBeziersFromCommands(commands)
		for _, bz := range beziers {
			points = append(points, bz.P0, bz.P3)
			for t := 0.05; t < 1.0; t += 0.05 {
				points = append(points, GetPointFromBezier(bz, t))
			}
		}
	}
	for _, point := range points {
		tailDist := PointDistance(Point{X: 0, Y: 0}, tail)
		pointDist := PointDistance(Point{X: 0, Y: 0}, point)
		if tailDist < pointDist {
			tail = point
		}
	}
	k := 0.0
	if tail.X <= 0 && tail.Y > 0 {
		// quadrant = 1 // 1 * pi
		k = 1
	} else if tail.X <= 0 && tail.Y < 0 {
		// quadrant = 2 // 0 * pi
		k = 0
	} else if tail.X > 0 && tail.Y <= 0 {
		// quadrant = 3 // 0 * pi
		k = 0
	} else {
		// quadrant = 4 // 1 * pi
		k = 1
	}
	a := 0.0
	if tail.Y != 0 {
		a = math.Atan(tail.X/tail.Y) + k*math.Pi
	} else if tail.X != 0 {
		a = math.Atan(-tail.Y/tail.X) + k*math.Pi
	}
	// if tail.X && tail.Y == 0 -> dont rotate
	// fmt.Printf("%s: %f\n", group.ID, a*180/math.Pi)
	return GroupApplyTransformation(group, Transformation{Rotation: a * 180 / math.Pi})
}

func parseBodypart(g Group) BodyPart {
	group := GroupCopy(g)
	group.ID = group.Label
	group.Label = group.Ellipses[0].Label
	anchor := findElementPosition(group, group.Label)
	group = GroupApplyTransformation(group, Transformation{Translation: Point{X: anchor.X * -1, Y: anchor.Y * -1}})
	CleanGroup(&group)
	if group.Label != "eye" && group.Label != "mouth" {
		group = GroupNormalizeRotation(group)
	}
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
	for _, group := range root.Groups {
		if group.Paths[0].Label == "body" {
			bodies = append(bodies, parseBody(group))
		} else {
			bodyparts = append(bodyparts, parseBodypart(group))
		}
	}
	return bodies, bodyparts
}
