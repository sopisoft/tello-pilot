package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"time"

	"pilot/input"
	"pilot/state"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"gobot.io/x/gobot/platforms/dji/tello"
)

var ModeNames = map[input.ControllerMode]string{
	input.ControllerModeGame: "Gaming",
	input.ControllerMode1:    "Mode 1",
	input.ControllerMode2:    "Mode 2",
}

func nextControllerMode(mode input.ControllerMode) input.ControllerMode {
	return (mode + 1) % 3
}

type PilotState struct {
	ControllerMode                      input.ControllerMode
	FlightData                          *tello.FlightData
	xboxImage                           *ebiten.Image
	elevator, throttle, aileron, rudder float32
}

func (sm *PilotState) Init(
	stackdeep int,
) {
	var err error
	sm.xboxImage, _, err = ebitenutil.NewImageFromFile("./assets/images/xbox.png")
	if err != nil {
		log.Fatal(err)
	}

	if Tello != nil {
		go func() {
			err := Tello.On(tello.FlightDataEvent, func(data interface{}) {
				sm.FlightData = data.(*tello.FlightData)
			})
			if err != nil {
				log.Fatal(err)
			}

		}()

		go func() {
			for range time.Tick(120 * time.Millisecond) {
				err := Tello.SetVector(sm.elevator, sm.throttle, sm.aileron, sm.rudder)
				if err != nil {
					log.Fatal(err)
				}
			}
		}()

	}

}

func (sm *PilotState) Update(
	stackdeep int,
) state.Result {
	go GamepadUpdate()

	maxAxis := ebiten.GamepadAxisType(ebiten.GamepadAxisCount(GamepadId))
	for i := 0; i < int(maxAxis); i++ {
		val := ebiten.GamepadAxisValue(GamepadId, ebiten.GamepadAxisType(i))
		if math.Abs(float64(val)) > 0.2 && math.Abs(float64(val)) < 0.8 {
			log.Printf("Axis %d: %f", i, val)
		}
	}

	pressedButtons := []ebiten.GamepadButton{}
	maxButtons := ebiten.GamepadButtonCount(GamepadId)
	for i := 0; i < int(maxButtons); i++ {
		if ebiten.IsGamepadButtonPressed(GamepadId, ebiten.GamepadButton(i)) {
			pressedButtons = append(pressedButtons, ebiten.GamepadButton(i))
		}
	}
	if len(pressedButtons) > 0 {
		log.Printf("Pressed Buttons: %v", pressedButtons)
	}

	sm.elevator, sm.throttle, sm.aileron, sm.rudder = input.Control(GamepadId, sm.ControllerMode)

	if Tello != nil {
		go func() {
			// switch true {
			// case inpututil.IsStandardGamepadButtonJustPressed(GamepadId, ebiten.StandardGamepadButtonLeftTop):
			// 	Tello.TakeOff()
			// case inpututil.IsStandardGamepadButtonJustPressed(GamepadId, ebiten.StandardGamepadButtonLeftBottom):
			// 	Tello.Land()
			// case inpututil.IsStandardGamepadButtonJustPressed(GamepadId, ebiten.StandardGamepadButtonLeftRight):
			// 	Tello.RightFlip()
			// case inpututil.IsStandardGamepadButtonJustPressed(GamepadId, ebiten.StandardGamepadButtonLeftLeft):
			// 	Tello.LeftFlip()
			// }
			if inpututil.IsGamepadButtonJustPressed(GamepadId, ebiten.GamepadButton10) {
				Tello.TakeOff()
			} else if inpututil.IsGamepadButtonJustPressed(GamepadId, ebiten.GamepadButton12) {
				Tello.Land()
			} else if inpututil.IsGamepadButtonJustPressed(GamepadId, ebiten.GamepadButton11) {
				Tello.RightFlip()
			} else if inpututil.IsGamepadButtonJustPressed(GamepadId, ebiten.GamepadButton13) {
				Tello.LeftFlip()
			}
		}()

		if sm.FlightData != nil && sm.throttle > 0.8 && sm.FlightData.OnGround {
			go Tello.TakeOff()
		}
	}

	// if inpututil.IsStandardGamepadButtonJustPressed(GamepadId, ebiten.StandardGamepadButtonCenterRight) {
	if inpututil.IsGamepadButtonJustPressed(GamepadId, ebiten.GamepadButton7) {
		sm.ControllerMode = nextControllerMode(sm.ControllerMode)
	}

	return state.Result{}
}

func (sm *PilotState) Draw(screen *ebiten.Image, stackdeep int) {
	white := ebiten.ColorScale{}
	white.ScaleWithColor(color.White)

	imgH := sm.xboxImage.Bounds().Dy() / 2
	geoM := ebiten.GeoM{}
	geoM.Scale(0.5, 0.5)
	geoM.Translate(10, 10)
	screen.DrawImage(sm.xboxImage, &ebiten.DrawImageOptions{
		GeoM: geoM,
	})
	WriteLeft(screen, "Now in "+ModeNames[sm.ControllerMode], 10, float64(imgH+40), 32, white)
	Write(screen, "Press Button 6 to change controller mode", 250, float64(imgH+72), 18, white)

	data := []struct {
		label string
		value string
	}{
		{"Elevator", fmt.Sprintf("%0.2f", sm.elevator)},
		{"Throttle", fmt.Sprintf("%0.2f", sm.throttle)},
		{"Aileron ", fmt.Sprintf("%0.2f", sm.aileron)},
		{"Rudder  ", fmt.Sprintf("%0.2f", sm.rudder)},
	}
	y := float64(40)
	for _, d := range data {
		WriteLeft(screen, fmt.Sprintf("%s : %s", d.label, d.value), 450, y, 22, white)
		y += 28
	}

	if sm.FlightData != nil {
		data := []struct {
			label string
			value string
		}{
			{"Battery", fmt.Sprintf("%d%%", sm.FlightData.BatteryPercentage)},
			{"Height", fmt.Sprintf("%dcm", sm.FlightData.Height)},
			{"Air Speed", fmt.Sprintf("%fcm/s", sm.FlightData.AirSpeed())},
			{"Ground Speed", fmt.Sprintf("%fcm/s", sm.FlightData.GroundSpeed())},
			{"Vertical Speed", fmt.Sprintf("%d", sm.FlightData.VerticalSpeed)},
		}

		x := float64(450)
		y := float64(160)
		batteryStatus := "OK"
		if sm.FlightData.BatteryLower && sm.FlightData.BatteryLow {
			batteryStatus = "Low"
		} else if sm.FlightData.BatteryLower {
			batteryStatus = "Lower"
		}
		WriteLeft(screen, batteryStatus, x+140, y, 18, white)

		for _, d := range data {
			WriteLeft(screen, fmt.Sprintf("%s: %s", d.label, d.value), x, y, 18, white)
			y += 20
		}

	}

	instructions := input.GetInstructions(sm.ControllerMode)
	usage_init_y := float64(350)
	for _, instruction := range instructions {
		WriteLeft(screen, instruction.Label, 50, usage_init_y, 24, white)
		Write(screen, instruction.Value, 280, usage_init_y, 24, white)
		usage_init_y += 30
	}
}
