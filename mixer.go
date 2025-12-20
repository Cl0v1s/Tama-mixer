package main

import (
	"fmt"
	"math"
	"slices"
	"strconv"
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

func Place(body Body, bodyparts []BodyPart) SVG {
	svg := body.Svg

	for _, point := range body.Points {
		index := slices.IndexFunc(bodyparts, func(part BodyPart) bool {
			return part.Label == point.Label
		})
		if index == -1 {
			continue
		}
		angle := 0.0
		location := point
		t, bezier := findClosestPointInPaths(GetPathsInSVG(svg), point, 2)
		if t >= 0 {
			angle = GetRotationFromBezier(bezier, t)
			fmt.Println(angle)
			location = GetPointFromBezier(bezier, t)
		}
		groups := bodyparts[index].Svg.Groups
		for i := 0; i < len(groups); i++ {
			groups[i].Transform = "rotate(" + strconv.FormatFloat(angle, 'f', -1, 64) + ", 0, 0) translate(" + strconv.FormatFloat(location.X, 'f', -1, 64) + "," + strconv.FormatFloat(location.Y, 'f', -1, 64) + ")"
		}
		svg.Groups = append(svg.Groups, groups...)
	}

	return svg
}
