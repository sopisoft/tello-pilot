package main

import (
	"log"

	"pilot/state"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
)

const (
	screenWidth  = 1280
	screenHeight = 720
)

var (
	Robot     *gobot.Robot
	Tello     *tello.Driver
	GamepadId ebiten.GamepadID
)

func GamepadUpdate() {
	gamepadIDsBuf := inpututil.AppendJustConnectedGamepadIDs(nil)
	for _, id := range gamepadIDsBuf {
		gamepadIDs := map[ebiten.GamepadID]struct{}{}
		gamepadIDs[id] = struct{}{}
		if ebiten.IsStandardGamepadLayoutAvailable(id) {
			GamepadId = id
		}
	}
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Tello Pilot")
	ebiten.SetVsyncEnabled(true)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)

	gms := &state.Machine{
		LayoutWidth:  screenWidth,
		LayoutHeight: screenHeight,
	}
	Titlesm := &InitialState{}

	gms.StateAdd(Titlesm)

	if err := ebiten.RunGame(gms); err != nil {
		log.Fatal(err)
	}
}
