package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func ParseFrame(filename string) Frame {
	frame := Frame{}
	frame.Filename = filename
	formRexg := regexp.MustCompile("([A-Za-z]+)_([0-9]+).+.svg")
	matches := formRexg.FindStringSubmatch(filename)
	if len(matches) < 3 {
		panic("Bad filename " + filename)
	}
	frame.Form = matches[1]
	bodyFrame, err := strconv.ParseInt(matches[2], 10, 32)
	if err != nil {
		panic(err)
	}
	frame.BodyFrame = int(bodyFrame)

	stepsRegexg := regexp.MustCompile(`.+\@(.+)\.svg`)
	matches = stepsRegexg.FindStringSubmatch(filename)
	if len(matches) < 2 {
		panic("Bad filename " + filename)
	}
	steps := strings.Split(matches[1], "+")
	completeRegx := regexp.MustCompile("(.+)=.+([0-9]+)")
	partialRegx := regexp.MustCompile("(.+)=.+")
	stepsMap := make(map[string]int)
	for _, step := range steps {
		if step == "" {
			continue
		}
		matches = completeRegx.FindStringSubmatch(step)
		if len(matches) < 3 {
			matches = partialRegx.FindStringSubmatch(step)
			if len(matches) < 2 {
				panic("Bad step" + step)
			}
			stepsMap[matches[1]] = 0
		} else {
			f, err := strconv.ParseInt(matches[2], 10, 32)
			if err != nil {
				panic(err)
			}
			stepsMap[matches[1]] = int(f)
		}
	}
	fmt.Println(stepsMap)

	return frame
}

func ParseFrames(path string) []Frame {
	entries, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	frames := make([]Frame, 0)
	for _, e := range entries {
		frames = append(frames, ParseFrame(e.Name()))
	}
	return frames
}
