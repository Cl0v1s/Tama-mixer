package main

import (
	"io"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// We order by body frame
func orderFramesIdle(frames []Frame) []Frame {
	slices.SortFunc(frames, func(a Frame, b Frame) int {
		if a.BodyFrame != b.BodyFrame {
			return a.BodyFrame - b.BodyFrame
		}
		if a.MouthFrame != b.MouthFrame {
			return a.MouthFrame - b.MouthFrame
		}
		if a.Arm1Frame != b.Arm1Frame {
			return a.Arm1Frame - b.Arm1Frame
		}
		if a.Arm2Frame != b.Arm2Frame {
			return a.Arm2Frame - b.Arm2Frame
		}
		if a.Leg1Frame != b.Leg1Frame {
			return a.Leg1Frame - b.Leg1Frame
		}
		if a.Leg2Frame != b.Leg2Frame {
			return a.Leg2Frame - b.Leg2Frame
		}
		return 0
	})
	// we delete frames with open mouth and moving legs or arms
	frames = slices.DeleteFunc(frames, func(frame Frame) bool {
		if frame.MouthFrame != 0 && (frame.Leg1Frame != 0 || frame.Leg2Frame != 0 || frame.Arm1Frame != 0 || frame.Arm2Frame != 0) {
			return true
		}
		return false
	})
	return frames
}

func orderFramesHappy(frames []Frame) []Frame {
	slices.SortFunc(frames, func(a Frame, b Frame) int {
		if a.BodyFrame != b.BodyFrame {
			return a.BodyFrame - b.BodyFrame
		}
		if a.MouthFrame != b.MouthFrame {
			return a.MouthFrame - b.MouthFrame
		}
		if a.Arm1Frame != b.Arm1Frame {
			return a.Arm1Frame - b.Arm1Frame
		}
		if a.Arm2Frame != b.Arm2Frame {
			return a.Arm2Frame - b.Arm2Frame
		}
		return 0
	})
	// we delete frames with open mouth and moving legs or arms
	frames = slices.DeleteFunc(frames, func(frame Frame) bool {
		if frame.BodyFrame != 0 || frame.Leg1Frame != frame.Leg2Frame {
			return true
		}
		return false
	})
	return frames
}

func orderFramesAngry(frames []Frame) []Frame {
	slices.SortFunc(frames, func(a Frame, b Frame) int {
		if a.BodyFrame != b.BodyFrame {
			return a.BodyFrame - b.BodyFrame
		}
		if a.MouthFrame != b.MouthFrame {
			return a.MouthFrame - b.MouthFrame
		}
		if a.Arm1Frame != b.Arm1Frame {
			return a.Arm1Frame - b.Arm1Frame
		}
		if a.Arm2Frame != b.Arm2Frame {
			return a.Arm2Frame - b.Arm2Frame
		}
		if a.Leg1Frame != b.Leg1Frame {
			return a.Leg1Frame - b.Leg1Frame
		}
		if a.Leg2Frame != b.Leg2Frame {
			return a.Leg2Frame - b.Leg2Frame
		}
		return 0
	})
	return frames
}

func orderFramesByExpression(frames []Frame) []Frame {
	if len(frames) == 0 {
		return frames
	}
	expression := frames[0].Expression
	switch expression {
	case int(Expression_Idle):
		return orderFramesIdle(frames)
	case int(Expression_Happy):
		return orderFramesHappy(frames)
	case int(Expression_Angry):
		return orderFramesAngry(frames)
	}
	return frames
}

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
	stepRegex := regexp.MustCompile("(.+)=(.+)-([0-9]+)")
	stepsMap := make(map[string]int)
	for _, step := range steps {
		if step == "" {
			continue
		}
		matches = stepRegex.FindStringSubmatch(step)
		if len(matches) < 4 {
			panic("Bad step" + step)
		}
		f, err := strconv.ParseInt(matches[3], 10, 32)
		if err != nil {
			panic(err)
		}
		stepsMap[matches[1]] = int(f)
	}
	value, ok := stepsMap[string(BodypartType_Eye)]
	if ok {
		frame.Expression = value
	}
	value, ok = stepsMap[string(BodypartType_Arm1)]
	if ok {
		frame.Arm1Frame = value
	}
	value, ok = stepsMap[string(BodypartType_Arm2)]
	if ok {
		frame.Arm2Frame = value
	}
	value, ok = stepsMap[string(BodypartType_Leg1)]
	if ok {
		frame.Leg1Frame = value
	}
	value, ok = stepsMap[string(BodypartType_Leg2)]
	if ok {
		frame.Leg2Frame = value
	}
	value, ok = stepsMap[string(BodypartType_Mouth)]
	if ok {
		frame.MouthFrame = value
	}
	return frame
}

func groupFramesByForm(frames []Frame) [][]Frame {
	if len(frames) == 0 {
		return [][]Frame{}
	}
	slices.SortFunc(frames, func(a Frame, b Frame) int {
		return strings.Compare(a.Form, b.Form)
	})
	results := make([][]Frame, 0)
	current := make([]Frame, 0)
	currentForm := frames[0].Form
	for _, frame := range frames {
		if frame.Form != currentForm {
			currentForm = frame.Form
			results = append(results, current)
			current = make([]Frame, 0)
		}
		current = append(current, frame)
	}
	if len(current) > 0 {
		results = append(results, current)
	}
	return results
}

func groupFramesByExpression(frames []Frame) [][]Frame {
	if len(frames) == 0 {
		return [][]Frame{}
	}
	slices.SortFunc(frames, func(a Frame, b Frame) int {
		return a.Expression - b.Expression
	})
	results := make([][]Frame, 0)
	current := make([]Frame, 0)
	currentExpression := frames[0].Expression
	for _, frame := range frames {
		if frame.Expression != currentExpression {
			currentExpression = frame.Expression
			results = append(results, current)
			current = make([]Frame, 0)
		}
		current = append(current, frame)
	}
	if len(current) > 0 {
		results = append(results, current)
	}
	return results
}

func ParseFrames(path string) [][][]Frame {
	entries, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	frames := make([]Frame, 0)
	for _, e := range entries {
		frames = append(frames, ParseFrame(e.Name()))
	}

	framesByForm := groupFramesByForm(frames)
	framesByFormAndExpression := make([][][]Frame, 0)
	for i := 0; i < len(framesByForm); i++ {
		framesByExpression := groupFramesByExpression(framesByForm[i])
		for u := 0; u < len(framesByExpression); u++ {
			framesByExpression[u] = orderFramesByExpression(framesByExpression[u])
		}
		framesByFormAndExpression = append(framesByFormAndExpression, framesByExpression)
	}

	return framesByFormAndExpression
}

func SaveFrames(sourceFolder string, outFolder string, frames []Frame) {
	for index, frame := range frames {
		src := sourceFolder + "/" + frame.Filename
		dstdir := outFolder + "/" + frame.Form + "/" + strconv.FormatInt(int64(frame.Expression), 10)
		dst := dstdir + "/" + strconv.FormatInt(int64(index), 10) + ".svg"
		os.MkdirAll(dstdir, 0750)
		fin, err := os.Open(src)
		if err != nil {
			panic(err)
		}
		defer fin.Close()
		fout, err := os.Create(dst)
		if err != nil {
			panic(err)
		}
		defer fout.Close()
		_, err = io.Copy(fout, fin)
		if err != nil {
			panic(err)
		}
	}
}
