package rpicamera

import (
	"encoding/base64"
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/dhowden/raspicam"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var log = logger.GetLogger("activity-rpicamera")

// List of input and output variables names
const (
	ivWidth        = "picWidth"
	ivHeight       = "picHeight"
	ivFlipH        = "HFlip"
	ivFlipV        = "VFlip"
	ivBrightness   = "Brightness"
	ivFolderOut    = "folderOut"
	ivOutputBase64 = "outputBase64"
	ovPath         = "picFile"
	ovBase64       = "base64"
)

// RPICamera is a stub for your Activity implementation
type RPICamera struct {
	metadata *activity.Metadata
}

// NewActivity creates a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &RPICamera{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *RPICamera) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *RPICamera) Eval(context activity.Context) (done bool, err error) {

	log.Debugf("Starting Eval function from activity rpicamera")

	// Set inputs
	picWidth := context.GetInput(ivWidth).(int)
	picHeight := context.GetInput(ivHeight).(int)
	flipH := context.GetInput(ivFlipH).(bool)
	flipV := context.GetInput(ivFlipV).(bool)
	brightness, _ := strconv.Atoi(ivBrightness)
	outputBase64 := context.GetInput(ivOutputBase64).(bool)
	targetFile := context.GetInput(ivFolderOut).(string)

	log.Debugf("Input defined: [picWidth = %d], [picHeight = %d], [flipH = %t], [flipV = %t], [brightness = %d], [targetFile = %s], [outputBase64 = %t]", picWidth, picHeight, flipH, flipV, brightness, targetFile, outputBase64)

	fileName := time.Now().Format("20060102150405")

	targetFile = filepath.Join(targetFile, fileName + ".jpg")
	log.Debugf("Picture will be saved to [%s]", targetFile)

	img, err := os.Create(targetFile)
	if err != nil {
		log.Errorf("error while creating output file: %v", err)
		return false, err
	}
	defer img.Close()

	log.Debug("Setting capture parameters...")
	s := raspicam.NewStill()
	s.BaseStill.Timeout = 1 * time.Millisecond
	s.BaseStill.Width = picWidth
	s.BaseStill.Height = picHeight
	s.BaseStill.Camera.Brightness = brightness
	s.BaseStill.Camera.HFlip = flipH
	s.BaseStill.Camera.VFlip = flipV

	errCh := make(chan error)
	go func() {
		for x := range errCh {
			log.Errorf("Error while capturing picture... [%v]", x)
		}
	}()

	log.Debug("Capturing image...")
	raspicam.Capture(s, img, errCh)

	/*if errCh != nil {
		return false, fmt.Errorf("error while capture picture... Please see logs.")
	}*/

	context.SetOutput(ovPath, targetFile)
	log.Debugf("Output variable picFile set to [%s]", targetFile)

	if outputBase64 {
		log.Debug("Export to base64 selected")
		log.Debugf("Reading picture file [%s]", targetFile)

		buf, err := ioutil.ReadFile(targetFile)
		if err != nil {
			log.Errorf("Error while reading picture... [%v]", err)
			return false, err
		}
		log.Debug("Encoding picture...")
		oBase64 := base64.StdEncoding.EncodeToString(buf)
		log.Debug("Picture encoded.")
		context.SetOutput(ovBase64, oBase64)
	}
	return true, nil
}
