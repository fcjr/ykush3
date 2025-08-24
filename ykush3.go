// Package ykush3 provides a Go library for controlling YKUSH3 USB switching devices.
//
// YKUSH3 is a 3-port USB switch that allows you to programmatically turn USB ports
// on and off. This is useful for power cycling USB devices, managing USB device
// connections, or automating hardware testing scenarios.
//
// Basic usage:
//
//	// Connect to the first available YKUSH3 device
//	ykush, err := ykush3.New()
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer ykush.Close()
//
//	// Turn on port 1
//	err = ykush.PortUp(ykush3.Port1)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Check port state
//	state, err := ykush.GetPortState(ykush3.Port1)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Port 1 is %s\n", state)
package ykush3

import (
	"errors"
	"fmt"

	"github.com/sstallion/go-hid"
)

const (
	// VendorID is the USB vendor ID for YKUSH3 devices.
	VendorID = 0x04D8
	// ProductID is the USB product ID for YKUSH3 devices.
	ProductID = 0xF11B
	// ReportSize is the HID report size used for communication.
	ReportSize = 64
)

// Port represents a USB port number on the YKUSH3 device.
type Port int

const (
	// Port1 is the first USB port.
	Port1 Port = 1
	// Port2 is the second USB port.
	Port2 Port = 2
	// Port3 is the third USB port.
	Port3 Port = 3
	// AllPorts represents all ports for bulk operations.
	AllPorts Port = 10
)

// PortState represents the on/off state of a USB port.
type PortState bool

const (
	// PortOff indicates the port is turned off.
	PortOff PortState = false
	// PortOn indicates the port is turned on.
	PortOn PortState = true
)

// YKUSH3 represents a connection to a YKUSH3 USB switching device.
type YKUSH3 struct {
	device *hid.Device
	serial string
}

// New creates a new YKUSH3 instance and opens the first available device.
func New() (*YKUSH3, error) {
	return NewWithSerial("")
}

// NewWithSerial creates a new YKUSH3 instance and opens the device with the specified serial number.
// If serial is empty, it opens the first available device.
func NewWithSerial(serial string) (*YKUSH3, error) {
	if err := hid.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize HID library: %w", err)
	}

	var device *hid.Device
	var err error

	if serial == "" {
		device, err = hid.OpenFirst(VendorID, ProductID)
	} else {
		device, err = hid.Open(VendorID, ProductID, serial)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to open YKUSH3 device: %w", err)
	}

	return &YKUSH3{
		device: device,
		serial: serial,
	}, nil
}

// Close closes the connection to the YKUSH3 device and releases resources.
func (y *YKUSH3) Close() error {
	if y.device != nil {
		err := y.device.Close()
		y.device = nil
		return err
	}
	return nil
}

// GetSerial returns the serial number of the connected YKUSH3 device.
func (y *YKUSH3) GetSerial() (string, error) {
	if y.device == nil {
		return "", errors.New("device not connected")
	}
	return y.device.GetSerialNbr()
}

// sendCommand sends a command to the device and returns the response
func (y *YKUSH3) sendCommand(cmd, ctrl byte) ([]byte, error) {
	if y.device == nil {
		return nil, errors.New("device not connected")
	}

	cmdBuf := make([]byte, ReportSize)
	cmdBuf[0] = cmd
	cmdBuf[1] = ctrl

	if _, err := y.device.Write(cmdBuf); err != nil {
		return nil, fmt.Errorf("failed to send command: %w", err)
	}

	respBuf := make([]byte, ReportSize)
	if _, err := y.device.Read(respBuf); err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return respBuf, nil
}

// PortUp turns on the specified USB port.
func (y *YKUSH3) PortUp(port Port) error {
	var cmd byte
	switch port {
	case Port1:
		cmd = 0x11
	case Port2:
		cmd = 0x12
	case Port3:
		cmd = 0x13
	case AllPorts:
		cmd = 0x1A
	default:
		return fmt.Errorf("invalid port: %d", port)
	}

	resp, err := y.sendCommand(cmd, cmd)
	if err != nil {
		return err
	}

	if resp[0] != 0x01 || resp[1] != cmd {
		return fmt.Errorf("unexpected response: status=0x%02x, response=0x%02x", resp[0], resp[1])
	}

	return nil
}

// PortDown turns off the specified USB port.
func (y *YKUSH3) PortDown(port Port) error {
	var cmd byte
	switch port {
	case Port1:
		cmd = 0x01
	case Port2:
		cmd = 0x02
	case Port3:
		cmd = 0x03
	case AllPorts:
		cmd = 0x0A
	default:
		return fmt.Errorf("invalid port: %d", port)
	}

	resp, err := y.sendCommand(cmd, cmd)
	if err != nil {
		return err
	}

	if resp[0] != 0x01 || resp[1] != cmd {
		return fmt.Errorf("unexpected response: status=0x%02x, response=0x%02x", resp[0], resp[1])
	}

	return nil
}

// GetPortState returns the current on/off state of the specified USB port.
func (y *YKUSH3) GetPortState(port Port) (PortState, error) {
	var cmd byte
	switch port {
	case Port1:
		cmd = 0x21
	case Port2:
		cmd = 0x22
	case Port3:
		cmd = 0x23
	default:
		return PortOff, fmt.Errorf("invalid port for state query: %d", port)
	}

	resp, err := y.sendCommand(cmd, cmd)
	if err != nil {
		return PortOff, err
	}

	if resp[0] != 0x01 {
		return PortOff, fmt.Errorf("command failed: status=0x%02x", resp[0])
	}

	switch resp[1] {
	case 0x01, 0x02, 0x03:
		return PortOff, nil
	case 0x11, 0x12, 0x13:
		return PortOn, nil
	default:
		return PortOff, fmt.Errorf("unexpected state response: 0x%02x", resp[1])
	}
}

// SetPortState sets the specified USB port to the given state (on or off).
func (y *YKUSH3) SetPortState(port Port, state PortState) error {
	if state {
		return y.PortUp(port)
	}
	return y.PortDown(port)
}

// AllPortsUp turns on all USB ports simultaneously.
func (y *YKUSH3) AllPortsUp() error {
	return y.PortUp(AllPorts)
}

// AllPortsDown turns off all USB ports simultaneously.
func (y *YKUSH3) AllPortsDown() error {
	return y.PortDown(AllPorts)
}

// GetAllPortsState returns the current state of all individual USB ports as a map.
func (y *YKUSH3) GetAllPortsState() (map[Port]PortState, error) {
	states := make(map[Port]PortState)

	ports := []Port{Port1, Port2, Port3}
	for _, port := range ports {
		state, err := y.GetPortState(port)
		if err != nil {
			return nil, fmt.Errorf("failed to get state for port %d: %w", port, err)
		}
		states[port] = state
	}

	return states, nil
}

// ListDevices returns information about all connected YKUSH3 devices on the system.
func ListDevices() ([]hid.DeviceInfo, error) {
	if err := hid.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize HID library: %w", err)
	}

	var devices []hid.DeviceInfo
	err := hid.Enumerate(VendorID, ProductID, func(info *hid.DeviceInfo) error {
		devices = append(devices, *info)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to enumerate devices: %w", err)
	}

	return devices, nil
}

// String returns a human-readable representation of the port.
func (p Port) String() string {
	switch p {
	case Port1:
		return "Port 1"
	case Port2:
		return "Port 2"
	case Port3:
		return "Port 3"
	case AllPorts:
		return "All Ports"
	default:
		return fmt.Sprintf("Port %d", int(p))
	}
}

// String returns a human-readable representation of the port state.
func (s PortState) String() string {
	if s {
		return "ON"
	}
	return "OFF"
}
