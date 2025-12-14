package main

import (
	"encoding/xml"
	"math"
	"os"
	"slices"
	"strings"
)

func Save(basepath string, svg SVG) {
	_ = os.MkdirAll(basepath, 0755)
	path := basepath + "/" + svg.XMLName.Space + "_" + svg.XMLName.Local + ".svg"
	out, _ := xml.MarshalIndent(svg, " ", "  ")
	os.WriteFile(path, out, 0644)
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
		for _, cmd := range cmds {
			if strings.ToLower(cmd.Type) != "m" {
				continue
			}
			x1, y1 := cmd.Args[0], cmd.Args[1]
			if x1 < x {
				x = x1
			}
			if y1 < y {
				y = y1
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
				cmds[i].Args[0] -= x
				cmds[i].Args[1] -= y
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

func RetrieveAnchors(group *Group, rootLabel string) []Point {
	points := make([]Point, 0)
	for _, g := range group.Groups {
		points = append(points, RetrieveAnchors(&g, rootLabel)...)
	}
	i := 0
	for _, e := range group.Ellipses {
		if e.Label == rootLabel {
			group.Ellipses[i] = e
			i++
		} else {
			points = append(points, Point{X: e.CX, Y: e.CY, Label: e.Label, Transform: e.Transform})
		}
	}
	group.Ellipses = group.Ellipses[:i]
	i = 0
	for _, e := range group.Circles {
		if e.Label == rootLabel {
			group.Circles[i] = e
			i++
		} else {
			points = append(points, Point{X: e.CX, Y: e.CY, Label: e.Label, Transform: e.Transform})
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
			points = append(points, Point{X: commands[index].Args[0], Y: commands[index].Args[1], Label: p.Label, Transform: p.Transform})
		}
	}
	group.Paths = group.Paths[:i]
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

func Sort(root SVG) ([]Body, []BodyPart) {
	bodies := make([]Body, 0)
	bodyparts := make([]BodyPart, 0)

	for _, character := range root.Groups {
		for _, group := range character.Groups {
			x, y := FindLowestPadding(group, group.Label)
			group = Unpad(group, x, y)

			if group.Label == "body" {
				anchors := RetrieveAnchors(&group, group.Label)
				svg := SVG{
					XMLName: xml.Name{
						Space: group.ID,
						Local: group.Label,
					},
					Xmlns:   "http://www.w3.org/2000/svg",
					Width:   "100",
					Height:  "100",
					ViewBox: "-10 -10 100 100",
					Groups:  []Group{group},
				}
				bodies = append(bodies, Body{
					Svg:    svg,
					Points: anchors,
				})
			} else {
				CleanGroup(&group)
				svg := SVG{
					XMLName: xml.Name{
						Space: group.ID,
						Local: group.Label,
					},
					Xmlns:   "http://www.w3.org/2000/svg",
					Width:   "100",
					Height:  "100",
					ViewBox: "-10 -10 100 100",
					Groups:  []Group{group},
				}
				bodyparts = append(bodyparts, BodyPart{
					Svg:   svg,
					Label: group.Label,
				})
			}
		}
	}

	return bodies, bodyparts
}
