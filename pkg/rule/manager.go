package rule

import (
	"fmt"
	"home-solar-pi/pkg/device"
	"log"
	"os"
	"path"
	"time"

	"github.com/Knetic/govaluate"
	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
)

var logger = log.Default()

type RuleManager struct {
	rules     []Rule
	dm        device.DeviceManager
	rulesPath string
}

func NewRuleManager(rulesPath string, dm device.DeviceManager) *RuleManager {

	rules, err := readRulesFromFile(rulesPath)

	if err != nil {
		panic("Rules file error")
	}

	return &RuleManager{
		rules:     rules,
		dm:        dm,
		rulesPath: rulesPath,
	}

}

func (m *RuleManager) GetAllRules() []Rule {
	return m.rules
}

// starts the rule server
// automatically watches rules change in files.
// runs rules indipendently in go rutines
// closes everything if a panic event occurs
func (m *RuleManager) StartRuleServer(panicChan chan error) {

	restartRules := make(chan bool)

	// watches the file changes inside directory
	go m.watchRuleChanges(restartRules)

	for {

		// creating channel array for closing rules
		closeChannels := make([]chan bool, len(m.GetAllRules()))
		for ir, rule := range m.GetAllRules() {

			// running new rules adding channel
			go func() {
				close := make(chan bool)
				closeChannels[ir] = close

				// running rule
				err := m.runRule(rule, close)
				if err != nil {
					panicChan <- err
				}
			}()
		}

		// waits 5 seconds because the IDE writes multiple times the files
		// wtf
		time.Sleep(time.Second * 5)

		// event of file rule.yml modified
		<-restartRules

		// close all running rules
		for _, close := range closeChannels {
			close <- true
		}

		// reading new rules from files
		newRules, err := readRulesFromFile(m.rulesPath)
		if err != nil {
			panicChan <- err
			return
		}

		// setting new rules
		m.rules = newRules

	}

}

// Watches the rule dir for changes
// if any change occurs sends an event on the restartRules channel
func (m *RuleManager) watchRuleChanges(restartRules chan bool) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				// logger.Println("event:", event)
				if !ok {
					continue
				}
				if event.Has(fsnotify.Write) {
					// logger.Println("modified file:", event.Name)
					restartRules <- true
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					continue
				}
				logger.Println("error:", err)
			}
		}
	}()

	// Add a path.
	err = watcher.Add(m.rulesPath)
	if err != nil {
		println(m.rulesPath)
		logger.Fatal(err)
		return
	}

	<-make(chan struct{})

}

// Runs a rule. Exits if close chan send an event
func (m *RuleManager) runRule(r Rule, close chan bool) error {

	for {

		select {
		case <-close:
			return nil
		default:
			{
				resultCondition, err := m.evalCondition(r.Condition)
				if err != nil {
					// return err
					logger.Printf("Error evaluating condition : %s\n", err.Error())
				} else {
					m.performAction(resultCondition, r)
				}

				time.Sleep(time.Second * time.Duration(r.RefreshInterval))
			}
		}

	}

}

// performs action for a device rule
// TODO implement other actions
func (m *RuleManager) performAction(conditionResult bool, r Rule) error {

	device, err := m.dm.GetDeviceByName(r.Device)
	if err != nil {
		return err
	}

	var actionToPerform RuleAction
	if conditionResult {
		actionToPerform = r.Action
	} else {
		actionToPerform = r.InverseAction
	}

	logger.Printf("Rule %s Action %s on %s\n", r.Name, actionToPerform, device.GetDeviceName())
	switch actionToPerform {
	case POWER_ON:
		device.PowerOn()
	case POWER_OFF:
		device.PowerOff()
	default:
		return fmt.Errorf("action not found :%s\n", actionToPerform)
	}

	return nil
}

// evaluates the condition of a rule using govaluate package
// automatically detects and injects values referring to other devices
// e.g. inverter > 500 -> inverter variable taken from inverterDevice.ReadValue()
func (m *RuleManager) evalCondition(cond string) (bool, error) {
	exp, err := govaluate.NewEvaluableExpression(cond)
	if err != nil {
		return false, err
	}

	vars := exp.Vars()

	parameters := make(map[string]any)

	for _, deviceName := range vars {

		deviceToGet, err := m.dm.GetDeviceByName(deviceName)
		if err != nil {
			return false, err
		}

		val, err := deviceToGet.ReadValue()
		if err != nil {
			return false, err
		}

		parameters[deviceName] = val

	}

	// fmt.Printf("%+v\n", parameters)

	evaluateResult, err := exp.Evaluate(parameters)
	if err != nil {
		return false, err
	}

	switch evaluateResult := evaluateResult.(type) {
	case bool:
		return evaluateResult, nil
	default:
		return false, fmt.Errorf("condition not a boolean := %s", cond)
	}
}

func readRulesFromFile(rulesPath string) ([]Rule, error) {
	rulesDir, err := os.ReadDir(rulesPath)
	if err != nil {
		return nil, err
	}

	rules := make([]Rule, 0)

	for _, ruleFile := range rulesDir {

		ruleYaml, err := os.ReadFile(path.Join(rulesPath, ruleFile.Name()))
		if err != nil {
			println("Failed -", ruleFile.Name())
			continue
		}

		var rule Rule
		err = yaml.Unmarshal(ruleYaml, &rule)
		if err != nil {
			log.Fatalf("Unmarshal: %v", err)
			continue
		}

		rules = append(rules, rule)
	}

	return rules, nil
}
