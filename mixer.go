package main

import (
	"math"
	"slices"
	"strconv"
)

const RANGE = 5

func findClosestPointInPaths(paths []Path, point Point) (float64, Bezier) {
	for _, p := range paths {
		commands := ParseD(p.D)
		beziers := GetBeziersFromCommands(commands)
		for _, b := range beziers {
			// TODO: peut largement être optimisé
			for i := 0.0; i < 1; i += 0.1 {
				p := GetPointFromBezier(b, i)
				d := math.Abs(p.X-point.X) + math.Abs(p.Y-point.Y)
				if d <= RANGE {
					return i, b
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
		groups := bodyparts[index].Svg.Groups
		for i := 0; i < len(groups); i++ {
			groups[i].Transform = "translate(" + strconv.FormatFloat(point.X, 'f', -1, 64) + "," + strconv.FormatFloat(point.Y, 'f', -1, 64) + ")"
		}
		svg.Groups = append(svg.Groups, groups...)
	}

	return svg
}
