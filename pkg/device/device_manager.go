package device

import (
	"log"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

type DeviceManager struct {
	devicePath    string
	logger        *log.Logger
	deviceDrivers []DeviceDriver
}

func NewDeviceManager(devicePath string, logger *log.Logger) (*DeviceManager, error) {
	deviceManager := DeviceManager{
		devicePath: devicePath,
		logger:     logger,
	}

	var err error
	deviceManager.deviceDrivers, err = deviceManager.GetAllDevices()

	if err != nil {
		return nil, err
	}

	return &deviceManager, nil
}

func (m DeviceManager) GetAllDevices() ([]DeviceDriver, error) {
	devicesFile, err := os.ReadDir(m.devicePath)
	if err != nil {
		log.Fatal(err)
	}

	devices := make([]DeviceDriver, 0)

	for _, deviceFile := range devicesFile {

		yamlFile, err := os.ReadFile(path.Join(m.devicePath, deviceFile.Name()))
		if err != nil {
			println("Failed -", deviceFile.Name())
			continue
		}

		var device Device
		err = yaml.Unmarshal(yamlFile, &device)
		if err != nil {
			log.Fatalf("Unmarshal: %v", err)
			continue
		}

		m.logger.Printf("%+v\n", device)

		switch device.Driver {
		case HeaterType:
			devices = append(devices, NewHeater(device))
		case InverterType:
			devices = append(devices, NewInterver(device))
		default:
			devices = append(devices, DeviceDriver(device))
		}

	}

	return devices, nil

}

func (m DeviceManager) GetDeviceDriver(id DriverType) (DeviceDriver, error) {

	// var device DeviceDriver
	// m.logger.Printf("----------- %+v %+v ----------\n", id, m.deviceDrivers)

	for _, dev := range m.deviceDrivers {
		m.logger.Printf("----------- %s\n", dev.GetDriverName())
		if dev.GetDriverName() == id {
			return dev, nil
		}
	}

	return nil, errNotFound
}

func (m DeviceManager) PowerOn(id DriverType) error {
	device, err := m.GetDeviceDriver(id)

	if err != nil {
		return err
	}

	return device.PowerOn()
}

func (m DeviceManager) PowerOff(id DriverType) error {
	device, err := m.GetDeviceDriver(id)

	if err != nil {
		return err
	}

	return device.PowerOff()
}

func (m DeviceManager) DeviceStatus(id DriverType) (DeviceStatus, error) {
	device, err := m.GetDeviceDriver(id)

	if err != nil {
		return INACTIVE, err
	}

	return device.Status()
}

// Status() (DeviceStatus, error)
// ReadValue() (any, error)
// GetConfig() (DeviceConfig, error)
