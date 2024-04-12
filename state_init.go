package main

import (
	"bytes"
	"image/color"
	"log"
	"os/exec"
	"sync"
	"time"

	"pilot/state"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
)

type InitialState struct {
	dots           string
	botStart       sync.Once
	telloConnected bool
	info           string
}

func (sm *InitialState) Init(
	stackdeep int,
) {
	sm.dots = "    "
	sm.telloConnected = false
	sm.info = "Now Initializing" + sm.dots

	go func() {
		var i int
		dot := []byte(".")
		whitespace := []byte(" ")
		t := time.NewTicker(1 * time.Second)
		for range t.C {
			i++
			n := i % 4
			if Robot == nil {
				sm.info = "Robot is nil"
			} else if Robot.Running() && !sm.telloConnected {
				sm.dots = string(bytes.Repeat(dot, n)) + string(bytes.Repeat(whitespace, 3-n))
				sm.info = "Not Connected" + sm.dots
			} else {
				t.Stop()
				break
			}
		}
	}()

	if Tello == nil {
		drone := tello.NewDriver("8888")
		work := func() {
			drone.On(tello.ConnectedEvent, func(data interface{}) {
				sm.telloConnected = true
				sm.info = "Press A to start"

				err := drone.StartVideo()
				drone.SetVideoEncoderRate(4)
				gobot.Every(60*time.Millisecond, func() {
					err := drone.StartVideo()
					if err != nil {
						log.Fatal(err)
					}

				})
				if err != nil {
					log.Fatal(err)
				}

				cmd := exec.Command("ffplay", "-probesize", "32", "-sync", "ext", "udp://127.0.0.1:11111", "-framerate", "30")
				ffplay_err := cmd.Start()
				if ffplay_err != nil {
					log.Fatal(err)
				}
			})
		}
		robot := gobot.NewRobot("tello",
			[]gobot.Connection{},
			[]gobot.Device{drone},
			work,
		)

		sm.botStart.Do(func() {
			go func() {
				var err = robot.Start()
				if err != nil {
					log.Fatal(err)
				}
			}()
		})

		Robot = robot
		Tello = drone
	}
}

func (sm *InitialState) Update(
	stackdeep int,
) state.Result {
	GamepadUpdate()
	// if inpututil.IsStandardGamepadButtonJustPressed(GamepadId, ebiten.StandardGamepadButtonCenterRight) && sm.telloConnected {
	if inpututil.IsGamepadButtonJustPressed(GamepadId, ebiten.GamepadButton0) && sm.telloConnected {
		return state.Result{
			Code:      state.Change,
			NextState: &PilotState{},
		}
	}
	return state.Result{}
}

func (sm *InitialState) Draw(screen *ebiten.Image, stackdeep int) {
	x, y := float64(screenWidth/2), float64(screenHeight/2)
	white := ebiten.ColorScale{}
	white.ScaleWithColor(color.White)
	Write(screen, "Tello", x, (y), 108, white)
	Write(screen, sm.info, x, (y + 100), 32, white)
}
