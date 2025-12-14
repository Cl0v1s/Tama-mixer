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
		groups := bodyparts[index].Svg.Groups
		for i := 0; i < len(groups); i++ {
			groups[i].Transform = "translate(" + strconv.FormatFloat(point.X, 'f', -1, 64) + "," + strconv.FormatFloat(point.Y, 'f', -1, 64) + ")"
		}
		svg.Groups = append(svg.Groups, groups...)
	}

	return svg
}
