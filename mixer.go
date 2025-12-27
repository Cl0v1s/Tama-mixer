package main

import (
	"slices"
	"strconv"
)

func bezierToD(beziers []Bezier) string {
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

func mergePaths(base Body, toMerge Group, kernel float64) Group {
	bodyBeziers := make([]Bezier, 0)
	for _, p := range GetPathsInSVG(base.Svg) {
		commands := ParseD(p.D)
		bodyBeziers = append(bodyBeziers, GetBeziersFromCommands(commands)...)
	}
	paths := GetPathsInGroup(toMerge)
	for u := 0; u < len(paths); u++ {
		commands := ParseD(paths[u].D)
		beziers := GetBeziersFromCommands(commands)
		for v := 0; v < len(beziers); v++ {
			// close to P0
			closeStartP0 := slices.IndexFunc(bodyBeziers, func(b Bezier) bool {
				return PointDistance(b.P0, beziers[v].P0) <= kernel
			})
			if closeStartP0 != -1 {
				beziers[v].P0 = bodyBeziers[closeStartP0].P0
			} else {
				closeStartP3 := slices.IndexFunc(bodyBeziers, func(b Bezier) bool {
					return PointDistance(b.P3, beziers[v].P0) <= kernel
				})
				if closeStartP3 != -1 {
					beziers[v].P0 = bodyBeziers[closeStartP3].P3
				}
			}
			// close to P3
			closeEndP0 := slices.IndexFunc(bodyBeziers, func(b Bezier) bool {
				return PointDistance(b.P0, beziers[v].P3) <= kernel
			})
			if closeEndP0 != -1 {
				beziers[v].P0 = bodyBeziers[closeEndP0].P0
			} else {
				closeEndP3 := slices.IndexFunc(bodyBeziers, func(b Bezier) bool {
					return PointDistance(b.P3, beziers[v].P3) <= kernel
				})
				if closeEndP3 != -1 {
					beziers[v].P0 = bodyBeziers[closeEndP3].P3
				}
			}
		}
		paths[u].D = bezierToD(beziers)
		bodyBeziers = append(bodyBeziers, beziers...)
	}
	toMerge.Paths = paths
	return toMerge
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
			location = GetPointFromBezier(bezier, t)
		}
		for _, group := range bodyparts[index].Svg.Groups {
			group = GroupApplyTransformation(group, Transformation{Rotation: angle})
			group = GroupApplyTransformation(group, Transformation{Translation: location})
			// group = mergePaths(body, group, 3)
			svg.Groups = append(svg.Groups, group)
		}
	}

	return svg
}
