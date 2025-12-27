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

func Mix(bodies []Body, bodyparts []BodyPart) []Body {
	ready := make([]Body, 0)
	bucket := make([]Body, 0)

	// we mix all existing bodypart with all compatible bodies
	for _, part := range bodyparts {
		for _, b := range bodies {
			if !BodyIsCompatible(b, part.Label) {
				continue
			}
			body := BodyCopy(b)
			body.Parts = append(body.Parts, part)
			bucket = append(bucket, body)

		}
	}

	for len(bucket) > 0 {
		body := bucket[0]
		err, point := BodyGetMissingPart(body)
		if err != nil {
			ready = append(ready, body)
		} else {
			pretendants := make([]BodyPart, len(bodyparts))
			copy(pretendants, bodyparts)
			pretendants = slices.DeleteFunc(pretendants, func(p BodyPart) bool { return p.Label != point.Label })
			for _, pr := range pretendants {
				c := BodyCopy(body)
				c.Parts = append(c.Parts, pr)
				bucket = append(bucket, c)
			}
		}
		bucket = bucket[1:]
	}
	return ready
}
