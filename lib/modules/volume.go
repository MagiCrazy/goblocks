package modules

import (
	"fmt"
	"os/exec"
	"regexp"

	"github.com/davidscholberg/go-i3barjson"
)

// Volume represents the configuration for the volume display block.
type Volume struct {
	BlockConfigBase `yaml:",inline"`
	MuteColor       string `yaml:"mute_color"`
}

// UpdateBlock updates the volume display block.
// Currently, only the ALSA master channel volume is supported.
func (c Volume) UpdateBlock(b *i3barjson.Block) {
	b.Color = c.Color
	fullTextFmt := fmt.Sprintf("%s%%s", c.Label)
	amixerCmd := "amixer"
	amixerArgs := []string{"-D", "default", "get", "Master"}
	out, err := exec.Command(amixerCmd, amixerArgs...).Output()
	if err != nil {
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}
	outStr := string(out)
	re := regexp.MustCompile(`\[([0-9]+%)\]\s\[(on|off)\]`)
	extractedInfo := re.FindAllStringSubmatch(outStr, 1)
	if extractedInfo == nil {
		b.FullText = fmt.Sprintf(fullTextFmt, "cannot parse amixer output")
		return
	}
	if extractedInfo[0][2] == "off" {
		b.Color = c.MuteColor
	}
	b.FullText = fmt.Sprintf(fullTextFmt, extractedInfo[0][1])
}
