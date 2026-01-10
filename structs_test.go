package main

import (
	"encoding/xml"
	"strings"
	"testing"
)

const same = `
<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<svg
   width="210mm"
   height="297mm"
   viewBox="0 0 210 297"
>
<g
     inkscape:groupmode="layer"
     id="layer38"
     inkscape:label="ball-0">
    <ellipse
       style="fill:#ff9631;fill-opacity:1;stroke:none;stroke-width:0.235414;stroke-dasharray:none;stroke-opacity:1"
       id="path10-77-60"
       cx="16.37044"
       cy="109.05718"
       rx="0.72288907"
       ry="0.76093584"
       inkscape:label="eye" />
    <ellipse
       style="fill:#ff9631;fill-opacity:1;stroke:none;stroke-width:0.235414;stroke-dasharray:none;stroke-opacity:1"
       id="path10-77-8"
       cx="20.13007"
       cy="109.0287"
       rx="0.72288907"
       ry="0.76093584"
       inkscape:label="eye" />
    <ellipse
       style="fill:#ff9631;fill-opacity:1;stroke:none;stroke-width:0.235414;stroke-dasharray:none;stroke-opacity:1"
       id="path10-77-2"
       cx="18.746181"
       cy="112.84247"
       rx="0.72288907"
       ry="0.76093584"
       inkscape:label="mouth" />
    <path
       id="path15"
       style="fill:none;stroke:#000000;stroke-width:0.264583"
       d="m 25.907079,111.16 c -3e-6,3.87473 -3.77164,5.95749 -7.853023,5.95749 -4.081384,0 -7.257709,-2.94265 -7.257711,-6.81738 -3e-6,-3.87474 3.308615,-5.56062 7.390003,-5.56062 4.081387,0 7.720734,2.54578 7.720731,6.42051 z"
       sodipodi:nodetypes="sssss"
       inkscape:label="body" />
  </g>
</svg>
`

func TestBodyGetMissingPart(t *testing.T) {
	var svg SVG
	file := strings.NewReader(same)
	decoder := xml.NewDecoder(file)
	decoder.Decode(&svg)
	bodies, _ := Sort(svg)
	body := bodies[0]
	_, got := BodyGetMissingPart(body)
	want := "eye"
	if got.Type != BodypartType(want) {
		t.Errorf("Got %s expected %s", got.Type, want)
	}
	body.Parts = append(body.Parts, BodyPart{Type: "eye"})
	_, got = BodyGetMissingPart(body)
	want = "eye"
	if got.Type != BodypartType(want) {
		t.Errorf("Got %s expected %s", got.Type, want)
	}
	body.Parts = append(body.Parts, BodyPart{Type: "eye"})
	_, got = BodyGetMissingPart(body)
	want = "mouth"
	if got.Type != BodypartType(want) {
		t.Errorf("Got %s expected %s", got.Type, want)
	}
}
