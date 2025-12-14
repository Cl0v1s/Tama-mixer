package main

import (
	"encoding/xml"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func Save(basepath string, bodyparts []Group, lowest float64, highest float64) {
	_ = os.MkdirAll(basepath, 0755)
	w := int64(highest-lowest) * 2
	for i := 0; i < len(bodyparts); i++ {
		part := bodyparts[i]
		path := basepath + "/" + part.ID + ".svg"
		svg := SVG{
			XMLName: xml.Name{
				Space: "",
				Local: part.Label,
			},
			Xmlns:   "http://www.w3.org/2000/svg",
			Groups:  []Group{part},
			Width:   strconv.FormatInt(w, 10),
			Height:  strconv.FormatInt(w, 10),
			ViewBox: fmt.Sprintf("%d %d %d %d", int64(lowest*2), int64(lowest*2), int64(w), int64(w)),
		}
		out, _ := xml.MarshalIndent(svg, " ", "  ")
		os.WriteFile(path, out, 0644)
	}
}

// Search for elements named as the group and calculate their lower coordinates
func FindLowerPadding(bodypart Group, label string) (float64, float64) {
	x, y := math.MaxFloat64, math.MaxFloat64
	for _, g := range bodypart.Groups {
		x1, y1 := FindLowerPadding(g, label)
		if x1 < x {
			x = x1
		}
		if y1 < y {
			y = y1
		}
	}

	for _, p := range bodypart.Paths {
		if label != "" && p.Label != label {
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

	for _, e := range bodypart.Ellipses {
		if label != "" && e.Label != label {
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

	for _, e := range bodypart.Circles {
		if label != "" && e.Label != label {
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

func Unpad(bodypart Group, x float64, y float64) Group {
	ellipsis := make([]Ellipse, 0)
	for _, el := range bodypart.Ellipses {
		el.CX -= x
		el.CY -= y
		ellipsis = append(ellipsis, el)
	}
	circles := make([]Circle, 0)
	for _, ci := range bodypart.Circles {
		ci.CX -= x
		ci.CY -= y
		circles = append(circles, ci)
	}
	paths := make([]Path, 0)
	for _, p := range bodypart.Paths {
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
	for _, g := range bodypart.Groups {
		g := Unpad(g, x, y)
		groups = append(groups, g)
	}
	group := Group{
		ID:       bodypart.ID,
		Label:    bodypart.Label,
		Paths:    paths,
		Groups:   groups,
		Ellipses: ellipsis,
		Circles:  circles,
	}

	return group
}

func Sort(root SVG) ([]Group, []Group) {
	bodies := make([]Group, 0)
	bodyparts := make([]Group, 0)
	for i := 0; i < len(root.Groups); i++ {
		element := root.Groups[i]
		for u := 0; u < len(element.Groups); u++ {
			node := element.Groups[u]
			if node.Label == "body" {
				bodies = append(bodies, node)
			} else {
				bodyparts = append(bodyparts, node)
			}
		}
	}
	return bodies, bodyparts
}

func Place(body Group, bodyparts []Group) Group {
	groups := append(make([]Group, 0), body.Groups...)
	for _, c := range body.Circles {
		var part *Group = nil
		for _, p := range bodyparts {
			if p.Label == c.Label {
				part = &p
			}
		}
		if part == nil {
			continue
		}
		g := Group{
			ID:        part.ID,
			Label:     part.Label,
			Groups:    part.Groups,
			Paths:     part.Paths,
			Ellipses:  part.Ellipses,
			Circles:   part.Circles,
			Transform: part.Transform + " translate(" + strconv.FormatFloat(c.CX, 'f', -1, 64) + "," + strconv.FormatFloat(c.CY, 'f', -1, 64) + ")",
		}
		groups = append(groups, g)
	}

	for _, c := range body.Ellipses {
		var part *Group = nil
		for _, p := range bodyparts {
			if p.Label == c.Label {
				part = &p
			}
		}
		if part == nil {
			continue
		}
		g := Group{
			ID:        part.ID,
			Label:     part.Label,
			Groups:    part.Groups,
			Paths:     part.Paths,
			Ellipses:  part.Ellipses,
			Circles:   part.Circles,
			Transform: part.Transform + " translate(" + strconv.FormatFloat(c.CX, 'f', -1, 64) + "," + strconv.FormatFloat(c.CY, 'f', -1, 64) + ")",
		}
		groups = append(groups, g)
	}

	return Group{
		ID:       body.ID,
		Label:    body.Label,
		Paths:    body.Paths,
		Ellipses: body.Ellipses,
		Circles:  body.Circles,
		Groups:   groups,
	}
}
