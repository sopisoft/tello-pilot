package input

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type ControllerMode int

const (
	ControllerMode2 ControllerMode = iota
	ControllerMode1
	ControllerModeGame
)

func smoothen(val float64) float32 {
	if math.Abs(val) < 0.1 {
		return 0
	} else {
		return float32(math.Round(val*100) / 100)
	}
}

func Control(id ebiten.GamepadID, mode ControllerMode) (e, t, a, r float32) {
	// Elevator - Positive : Forward
	// Throttle - Positive : Up
	// Aileron  - Positive : Right
	// Rudder   - Positive : Clockwise
	var elevator, throttle, aileron, rudder float64 // [-1.0 - 1.0]

	var (
		// lv float64 = ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickVertical) * -1
		lv float64 = ebiten.GamepadAxisValue(id, 1) * -1
		// rv float64 = ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisRightStickVertical) * -1
		rv float64 = ebiten.GamepadAxisValue(id, 3) * -1
	)
	var (
		// lh float64 = ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickHorizontal)
		lh float64 = ebiten.GamepadAxisValue(id, 0)
		// rh float64 = ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisRightStickHorizontal)
		rh float64 = ebiten.GamepadAxisValue(id, 2)
	)

	// log.Printf("id %d, lv: %f, rv: %f, lh: %f, rh: %f", id, lv, rv, lh, rh)

	switch mode {
	case ControllerMode1:
		{
			elevator = lv
			throttle = rv
			aileron = rh
			rudder = lh
		}
	case ControllerMode2:
		{
			elevator = rv
			throttle = lv
			aileron = rh
			rudder = lh
		}
	case ControllerModeGame:
		{
			// throttleUp := ebiten.StandardGamepadButtonValue(id, ebiten.StandardGamepadButtonFrontTopRight)
			// throttleDown := ebiten.StandardGamepadButtonValue(id, ebiten.StandardGamepadButtonFrontTopLeft)
			throttleUp := ebiten.IsGamepadButtonPressed(id, ebiten.GamepadButton5)
			throttleDown := ebiten.IsGamepadButtonPressed(id, ebiten.GamepadButton4)

			elevator = lv
			// throttle = (throttleUp - throttleDown) * 0.6 // [-0.6 - 0.6]
			if throttleUp && !throttleDown {
				throttle = 0.6
			} else if !throttleUp && throttleDown {
				throttle = -0.6
			} else {
				throttle = 0
			}
			aileron = lh
			rudder = rh
		}
	}

	return smoothen(elevator), smoothen(throttle), smoothen(aileron), smoothen(rudder)
}
