package main

import (
	"encoding/xml"
	"math"
	"os"
	"regexp"
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

func CompileD(commands []Command) string {
	results := make([]string, 0)
	for _, command := range commands {
		args := make([]string, 0)
		if command.Type == "A" || command.Type == "a" {
			for i := 0; i < len(command.Args); i += 7 {
				rx := strconv.FormatFloat(command.Args[i], 'f', 8, 64)
				ry := strconv.FormatFloat(command.Args[i+1], 'f', 8, 64)
				rotation := strconv.FormatFloat(command.Args[i+2], 'f', 8, 64)

				// Flags : large-arc et sweep, forcés à 0 ou 1
				largeArcFlag := "0"
				if command.Args[i+3] != 0 {
					largeArcFlag = "1"
				}
				sweepFlag := "0"
				if command.Args[i+4] != 0 {
					sweepFlag = "1"
				}

				x := strconv.FormatFloat(command.Args[i+5], 'f', 8, 64)
				y := strconv.FormatFloat(command.Args[i+6], 'f', 8, 64)

				args = append(args, rx, ry, rotation, largeArcFlag, sweepFlag, x, y)
			}
		} else {
			for _, arg := range command.Args {
				args = append(args, strconv.FormatFloat(arg, 'f', 8, 64))
			}
		}
		results = append(results, command.Type+" "+strings.Join(args, " "))
	}
	return strings.Join(results, " ")
}

func ParseD(d string) []Command {
	// parsing path string to commands array
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

	// converting relative commands to absolute ones
	current := Point{}
	var zPoint *Point

	for i := 0; i < len(commands); i++ {
		cmd := &commands[i]
		switch cmd.Type {
		case "M", "m", "L", "l":
			for u := 0; u < len(cmd.Args); u += 2 {
				x, y := cmd.Args[u], cmd.Args[u+1]
				if cmd.Type == "m" || cmd.Type == "l" {
					x += current.X
					y += current.Y
				}
				cmd.Args[u], cmd.Args[u+1] = x, y
				current = Point{X: x, Y: y}
			}
			// point de fermeture pour 'Z'
			zPoint = &Point{X: current.X, Y: current.Y}

		case "C", "c":
			for u := 0; u < len(cmd.Args); u += 6 {
				if cmd.Type == "c" {
					cmd.Args[u] += current.X
					cmd.Args[u+1] += current.Y
					cmd.Args[u+2] += current.X
					cmd.Args[u+3] += current.Y
					cmd.Args[u+4] += current.X
					cmd.Args[u+5] += current.Y
				}
				current = Point{X: cmd.Args[u+4], Y: cmd.Args[u+5]}
			}

		case "A", "a":
			for u := 0; u < len(cmd.Args); u += 7 {
				if cmd.Type == "a" {
					// seule la position finale du point change
					cmd.Args[u+5] += current.X
					cmd.Args[u+6] += current.Y
				}
				current = Point{X: cmd.Args[u+5], Y: cmd.Args[u+6]}
			}

		case "H", "h": // ligne horizontale
			for u := 0; u < len(cmd.Args); u++ {
				if cmd.Type == "h" {
					cmd.Args[u] += current.X
				}
				current.X = cmd.Args[u]
			}

		case "V", "v": // ligne verticale
			for u := 0; u < len(cmd.Args); u++ {
				if cmd.Type == "v" {
					cmd.Args[u] += current.Y
				}
				current.Y = cmd.Args[u]
			}

		case "Z", "z":
			if zPoint != nil {
				current = *zPoint
			}
			zPoint = nil
		}

		cmd.Type = strings.ToUpper(cmd.Type)
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

func GetRotationFromBezierRadian(bezier Bezier, t float64) float64 {
	dx := 3*(1-t)*(1-t)*(bezier.P1.X-bezier.P0.X) + 6*(1-t)*t*(bezier.P2.X-bezier.P1.X) + 3*t*t*(bezier.P3.X-bezier.P2.X)
	dy := 3*(1-t)*(1-t)*(bezier.P1.Y-bezier.P0.Y) + 6*(1-t)*t*(bezier.P2.Y-bezier.P1.Y) + 3*t*t*(bezier.P3.Y-bezier.P2.Y)
	if dx == 0 {
		if dy > 0 {
			return 0.5 * math.Pi
		} else {
			return 0.5 * math.Pi
		}
	}
	angle := math.Atan(dy / dx)
	return angle
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

func RetrievePoints(group *Group, rootLabel string) ([]Point, Point) {
	points := make([]Point, 0)
	for _, g := range group.Groups {
		pts, _ := RetrievePoints(&g, rootLabel)
		points = append(points, pts...)
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

	// get barycentre from anchor points
	baryCentre := Point{X: 0, Y: 0}
	for _, point := range points {
		baryCentre = baryCentre.Add(point)
	}
	baryCentre.X = baryCentre.X / float64(len(points))
	baryCentre.Y = baryCentre.Y / float64(len(points))

	// calculate all points around bodyshape and get tan corrected by quadrant for each of them
	size := Point{}
	bodyPoints := make([]Point, 0)
	for _, path := range GetPathsInGroup(*group) {
		cmds := ParseD(path.D)
		bzs := GetBeziersFromCommands(cmds)
		for _, bz := range bzs {
			for u := 0.0; u <= 1.0; u += 0.05 {
				location := GetPointFromBezier(bz, u)
				normalizedLocation := location.Sub(baryCentre)
				quadrant := normalizedLocation.Quadrant()
				k := 1.0
				if quadrant == 2 || quadrant == 3 {
					k = 0
				}
				angle := (GetRotationFromBezierRadian(bz, u) + k*math.Pi) * 180 / math.Pi
				location.T = angle

				if location.X > size.X {
					size.X = location.X
				}
				if location.Y > size.Y {
					size.Y = location.Y
				}

				bodyPoints = append(bodyPoints, location)
			}
		}
	}

	// Place anchor on exact body point
	for u := 0; u < len(points); u++ {
		shortest := math.MaxFloat64
		for _, point := range bodyPoints {
			distance := points[u].Distance(point)
			if shortest > distance && distance < 2 {
				points[u].X = point.X
				points[u].Y = point.Y
				points[u].T = point.T
				shortest = distance
			}
		}
	}

	slices.SortFunc(points, func(a Point, b Point) int {
		return PointsOrder[string(a.Type)] - PointsOrder[string(b.Type)]
	})
	return points, size
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
	group = group.Transform(Transformation{Translation: Point{X: -x, Y: -y}})
	anchors, size := RetrievePoints(&group, group.Label)
	frameReg := regexp.MustCompile("(.+)-([0-9]+)")
	matches := frameReg.FindStringSubmatch(group.ID)
	if len(matches) < 3 {
		panic(group.Label + "-" + group.ID + " bad name")
	}
	frame, err := strconv.ParseInt(matches[2], 10, 64)
	if err != nil {
		panic(err)
	}
	return Body{
		Path:   group.GetPath().D,
		Points: anchors,
		Frame:  int(frame),
		Name:   matches[1],
		Size:   size,
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
		tailDist := tail.Distance(Point{X: 0, Y: 0})
		pointDist := point.Distance(Point{X: 0, Y: 0})
		if tailDist < pointDist {
			tail = point
		}
	}
	quadrant := tail.Quadrant()
	k := 1.0
	if quadrant == 2 || quadrant == 3 {
		k = 0
	}
	a := 0.0
	if tail.Y != 0 {
		a = math.Atan(tail.X/tail.Y) + k*math.Pi
	} else if tail.X != 0 {
		a = math.Atan(-tail.Y/tail.X) + k*math.Pi
	}
	return group.Transform(Transformation{Rotation: a * 180 / math.Pi})
}

func parseBodypart(g Group) BodyPart {
	group := GroupCopy(g)
	group.ID = group.Label
	group.Label = group.Ellipses[0].Label
	anchor := findElementPosition(group, group.Label)
	group = group.Transform(Transformation{Translation: Point{X: anchor.X * -1, Y: anchor.Y * -1}})
	CleanGroup(&group)
	if group.Label != "eye" && group.Label != "mouth" {
		group = GroupNormalizeRotation(group)
	}
	path := group.GetPath()
	bb := path.GetBoundingBox()

	frameReg := regexp.MustCompile("(.+)-([0-9]+)")
	matches := frameReg.FindStringSubmatch(group.ID)
	if len(matches) < 3 {
		panic(group.Label + "-" + group.ID + " bad name")
	}
	frame, err := strconv.ParseInt(matches[2], 10, 64)
	if err != nil {
		panic(err)
	}
	return BodyPart{
		BoundingBox: bb,
		Path:        path.D,
		Type:        BodypartType(group.Label),
		Frame:       int(frame),
		Name:        matches[1],
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

func GroupBodyParts(bodyparts []BodyPart) [][]BodyPart {
	if len(bodyparts) == 0 {
		return [][]BodyPart{}
	}

	slices.SortFunc(bodyparts, func(a, b BodyPart) int {
		if cmp := strings.Compare(a.Name, b.Name); cmp != 0 {
			return cmp
		}
		return strings.Compare(string(a.Type), string(b.Type))
	})

	results := make([][]BodyPart, 0)
	currentGroup := make([]BodyPart, 0)

	prevName := bodyparts[0].Name
	prevType := bodyparts[0].Type

	for i := range bodyparts {
		bp := bodyparts[i]

		if bp.Name != prevName || bp.Type != prevType {
			if len(currentGroup) > 0 {
				results = append(results, currentGroup)
				currentGroup = make([]BodyPart, 0)
			}
			prevName = bp.Name
			prevType = bp.Type
		}

		currentGroup = append(currentGroup, bp)
	}

	if len(currentGroup) > 0 {
		results = append(results, currentGroup)
	}

	return results
}

func GroupBodies(bodies []Body) [][]Body {
	if len(bodies) == 0 {
		return [][]Body{}
	}
	slices.SortFunc(bodies, func(a Body, b Body) int {
		return strings.Compare(a.Name, b.Name)
	})
	results := make([][]Body, 0)
	current := make([]Body, 0)
	currentExpression := bodies[0].Name
	for _, bodyparts := range bodies {
		if bodyparts.Name != currentExpression {
			currentExpression = bodyparts.Name
			results = append(results, current)
			current = make([]Body, 0)
		}
		current = append(current, bodyparts)
	}
	if len(current) > 0 {
		results = append(results, current)
	}
	return results
}
