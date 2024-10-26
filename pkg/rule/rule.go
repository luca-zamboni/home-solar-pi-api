package rule

type RuleAction string

const (
	POWER_ON  RuleAction = "power_on"
	POWER_OFF RuleAction = "power_off"
)

type Rule struct {
	Name            string     `yaml:"Name"`
	RefreshInterval int        `yaml:"RefreshInterval"`
	Action          RuleAction `yaml:"Action"`
	InverseAction   RuleAction `yaml:"InverseAction"`
	Condition       string     `yaml:"Condition"`
	Device          string     `yaml:"Device"`
}
