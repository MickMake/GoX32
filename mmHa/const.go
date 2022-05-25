package mmHa

const (
	// LabelSwitch       = "switch"
	// LabelSelect       = "select"
	// LabelBinarySensor = "binary"

	ConfigTopicSuffix = "config"
	StateTopicSuffix = "state"
	CmdTopicSuffix   = "set"
)

type Labels []string

func (l *Labels) ValueExists(value string) bool {
	var ok bool
	for _, l := range *l {
		if l == value {
			ok = true
		}
	}
	return ok
}