package hcsr04_streaming

// SENSOR READING COMING FROM https://github.com/hyhc2016/hc-sr04/blob/master/sensor.go

import (
	"context"
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/stianeikeland/go-rpio"
	"os"
	"strconv"
	"time"
)

// log is the default package logger
var log = logger.GetLogger("trigger-hc-sr04-rpi")

var interval = 1000

var (
	pin_trig rpio.Pin = rpio.Pin(23)
	pin_echo rpio.Pin = rpio.Pin(18)
)

// HCSR04Factory My Trigger factory
type HCSR04Factory struct {
	metadata *trigger.Metadata
}

//NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &HCSR04Factory{metadata: md}
}

//New Creates a new trigger instance for a given id
func (t *HCSR04Factory) New(config *trigger.Config) trigger.Trigger {
	return &HCSR04Trigger{metadata: t.metadata, config: config}
}

// HCSR04Trigger is a stub for your Trigger implementation
type HCSR04Trigger struct {
	metadata *trigger.Metadata
	runner   action.Runner
	config   *trigger.Config
}

func doEvery(d time.Duration, f func()) {
	for _ = range time.Tick(d) {
		f()
	}
}

// Init implements trigger.Trigger.Init
func (t *HCSR04Trigger) Init(runner action.Runner) {
	t.runner = runner

	if t.config.Settings == nil {
		log.Info("No configuration set for the trigger... Using default configuration...")
	} else {
		if t.config.Settings["delay_ms"] != nil && t.config.Settings["delay_ms"] != "" {
			interval, _ = strconv.Atoi(t.config.GetSetting("delay_ms"))
		} else {
			log.Infof("No delay has been set. Using default value (", interval, "ms)")
		}
	}
	log.Info("Opening GPIO connection...")
	err := rpio.Open()
	if err != nil {
		fmt.Println(err)
		log.Errorf("An error occured while opening GPIO port. [%s]. Exiting.", err)
		os.Exit(1)
	}
	defer rpio.Close()
	pin_trig.Output()
	pin_echo.Input()

	time.Sleep(time.Second * 2)

}

// Metadata implements trigger.Trigger.Metadata
func (t *HCSR04Trigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Start implements trigger.Trigger.Start
func (t *HCSR04Trigger) Start() error {
	// start the trigger
	log.Debug("Start Trigger HC-SR04 for Raspberry PI")
	handlers := t.config.Handlers

	log.Debug("Processing handlers")
	for _, handler := range handlers {
		t.scheduleRepeating(handler)
		log.Debugf("Processing Handler: %s", handler.ActionId)
	}
	return nil
}

// Stop implements trigger.Trigger.Start
func (t *HCSR04Trigger) Stop() error {
	// stop the trigger
	return nil
}

func (t *HCSR04Trigger) scheduleRepeating(endpoint *trigger.HandlerConfig) {

	log.Debug("Repeating every ", interval, "ms")

	fn2 := func() {
		act := action.Get(endpoint.ActionId)
		data := make(map[string]interface{})

		distance, err := t.checkDistance(endpoint)
		if err != nil {
			log.Error("Error while reading sensor data. Err: ", err.Error())
		}

		data["distance"] = distance

		log.Debugf("Distance: [%fmm]", distance)
		startAttrs, err := t.metadata.OutputsToAttrs(data, true)

		if err != nil {
			log.Errorf("After run error' %s'\n", err)
		}

		ctx := trigger.NewContext(context.Background(), startAttrs)
		results, err := t.runner.RunAction(ctx, act, nil)

		if err != nil {
			log.Errorf("An error occured while starting the flow. Err:", err)
		}
		log.Info("Exec: ", results)
	}

	// schedule repeating
	doEvery(time.Duration(interval)*time.Millisecond, fn2)
}

func (t *HCSR04Trigger) checkDistance(endpoint *trigger.HandlerConfig) (distance float64, err error) {
	/*pin_trig.Low()
	time.Sleep(time.Microsecond * 30)
	pin_trig.High()
	time.Sleep(time.Microsecond * 30)
	pin_trig.Low()
	time.Sleep(time.Microsecond * 30)
	for {
		status := pin_echo.Read()
		if status == rpio.High {
			break
		}
	}
	begin := time.Now()
	for {
		status := pin_echo.Read()
		if status == rpio.Low {
			break
		}
	}
	end := time.Now()
	diff := end.Sub(begin)
	result_sec := float64(diff.Nanoseconds()) / 1000000000.0
	return result_sec * 17150, nil*/
	return 2000, nil
}
