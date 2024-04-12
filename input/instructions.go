package input

type usage struct {
	Name                   string
	GoForward              string
	GoBackward             string
	GoRight                string
	GoLeft                 string
	GoUp                   string
	GoDown                 string
	RotateClockwise        string
	RotateCounterClockwise string
}

var modeUsage = map[ControllerMode]usage{
	ControllerMode1: {
		Name:                   "モード1",
		GoForward:              "左スティック上",
		GoBackward:             "左スティック下",
		GoRight:                "右スティック右",
		GoLeft:                 "右スティック左",
		GoUp:                   "右スティック上",
		GoDown:                 "右スティック下",
		RotateClockwise:        "左スティック右",
		RotateCounterClockwise: "左スティック左",
	},
	ControllerMode2: {
		Name:                   "モード2",
		GoForward:              "右スティック上",
		GoBackward:             "右スティック下",
		GoRight:                "右スティック右",
		GoLeft:                 "右スティック左",
		GoUp:                   "左スティック上",
		GoDown:                 "左スティック下",
		RotateClockwise:        "左スティック右",
		RotateCounterClockwise: "左スティック左",
	},
	ControllerModeGame: {
		Name:                   "ゲームモード",
		GoForward:              "左スティック上",
		GoBackward:             "左スティック下",
		GoRight:                "左スティック右",
		GoLeft:                 "左スティック左",
		GoUp:                   "手前トリガー右",
		GoDown:                 "手前トリガー左",
		RotateClockwise:        "右スティック右",
		RotateCounterClockwise: "右スティック左",
	},
}

func GetInstructions(mode ControllerMode) []struct {
	Label string
	Value string
} {
	usage := modeUsage[mode]
	return []struct {
		Label string
		Value string
	}{
		{"操作方法", usage.Name},
		{"前進", usage.GoForward},
		{"後退", usage.GoBackward},
		{"右", usage.GoRight},
		{"左", usage.GoLeft},
		{"上昇", usage.GoUp},
		{"下降", usage.GoDown},
		{"時計回り", usage.RotateClockwise},
		{"反時計回り", usage.RotateCounterClockwise},
	}
}
