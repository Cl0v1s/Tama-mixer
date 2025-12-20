package main

import (
	"slices"
	"strconv"
)

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
			location = GetPointFromBezier(bezier, t)
		}
		for _, group := range bodyparts[index].Svg.Groups {
			group.Transform = "translate(" + strconv.FormatFloat(location.X, 'f', -1, 64) + "," + strconv.FormatFloat(location.Y, 'f', -1, 64) + ") rotate(" + strconv.FormatFloat(angle, 'f', -1, 64) + ") " + group.Transform
			svg.Groups = append(svg.Groups, group)
		}
	}

	return svg
}
