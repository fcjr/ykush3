# YKUSH3 Go Library

[![Go Reference](https://pkg.go.dev/badge/github.com/fcjr/ykush3.svg)](https://pkg.go.dev/github.com/fcjr/ykush3)
[![Go Report Card](https://goreportcard.com/badge/github.com/fcjr/ykush3)](https://goreportcard.com/report/github.com/fcjr/ykush3)

A Go library for controlling YKUSH3 USB switching devices. The YKUSH3 is a 3-port USB switch that allows you to programmatically turn USB ports on and off, making it perfect for power cycling USB devices, managing USB device connections, or automating hardware testing scenarios.

## Features

- 🔌 Control individual USB ports (Port 1, 2, 3)
- ⚡ Bulk operations (turn all ports on/off simultaneously)
- 🔍 Query port states and device information
- 📡 Support for multiple devices via serial number

## Installation

```bash
go get github.com/fcjr/ykush3
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/fcjr/ykush3"
)

func main() {
    // Connect to the first available YKUSH3 device
    ykush, err := ykush3.New()
    if err != nil {
        log.Fatal(err)
    }
    defer ykush.Close()

    // Turn on port 1
    err = ykush.PortUp(ykush3.Port1)
    if err != nil {
        log.Fatal(err)
    }

    // Check port state
    state, err := ykush.GetPortState(ykush3.Port1)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Port 1 is %s\n", state) // Output: Port 1 is ON
}
```

## Usage Examples

### Basic Port Control

```go
// Turn ports on
ykush.PortUp(ykush3.Port1)    // Turn on port 1
ykush.PortUp(ykush3.Port2)    // Turn on port 2
ykush.AllPortsUp()            // Turn on all ports

// Turn ports off
ykush.PortDown(ykush3.Port1)  // Turn off port 1
ykush.AllPortsDown()          // Turn off all ports

// Use SetPortState for conditional control
ykush.SetPortState(ykush3.Port1, ykush3.PortOn)  // Turn on
ykush.SetPortState(ykush3.Port1, ykush3.PortOff) // Turn off
```

### Checking Port States

```go
// Check individual port state
state, err := ykush.GetPortState(ykush3.Port1)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Port 1 is %s\n", state)

// Get all port states at once
states, err := ykush.GetAllPortsState()
if err != nil {
    log.Fatal(err)
}

for port, state := range states {
    fmt.Printf("%s: %s\n", port, state)
}
```

### Device Management

```go
// List all connected YKUSH3 devices
devices, err := ykush3.ListDevices()
if err != nil {
    log.Fatal(err)
}

for i, device := range devices {
    fmt.Printf("Device %d: %s (Serial: %s)\n", 
        i+1, device.ProductStr, device.SerialNbr)
}

// Connect to a specific device by serial number
ykush, err := ykush3.NewWithSerial("YK12345")
if err != nil {
    log.Fatal(err)
}
defer ykush.Close()

// Get the serial number of connected device
serial, err := ykush.GetSerial()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Connected to device: %s\n", serial)
```

## API Reference

### Types

| Type | Description |
|------|-------------|
| `YKUSH3` | Represents a connection to a YKUSH3 device |
| `Port` | USB port number (Port1, Port2, Port3, AllPorts) |
| `PortState` | Port state (PortOn, PortOff) |

### Functions

#### Connection Management

| Function | Description |
|----------|-------------|
| `New() (*YKUSH3, error)` | Connect to first available device |
| `NewWithSerial(serial string) (*YKUSH3, error)` | Connect to device with specific serial |
| `(*YKUSH3) Close() error` | Close connection and release resources |
| `(*YKUSH3) GetSerial() (string, error)` | Get device serial number |
| `ListDevices() ([]hid.DeviceInfo, error)` | List all connected devices |

#### Port Control

| Function | Description |
|----------|-------------|
| `(*YKUSH3) PortUp(port Port) error` | Turn on specified port |
| `(*YKUSH3) PortDown(port Port) error` | Turn off specified port |
| `(*YKUSH3) SetPortState(port Port, state PortState) error` | Set port to specific state |
| `(*YKUSH3) AllPortsUp() error` | Turn on all ports |
| `(*YKUSH3) AllPortsDown() error` | Turn off all ports |

#### State Queries

| Function | Description |
|----------|-------------|
| `(*YKUSH3) GetPortState(port Port) (PortState, error)` | Get state of specific port |
| `(*YKUSH3) GetAllPortsState() (map[Port]PortState, error)` | Get state of all ports |

### Constants

| Constant | Value | Description |
|----------|-------|-------------|
| `Port1` | 1 | First USB port |
| `Port2` | 2 | Second USB port |
| `Port3` | 3 | Third USB port |
| `AllPorts` | 10 | All ports (for bulk operations) |
| `PortOn` | true | Port is turned on |
| `PortOff` | false | Port is turned off |

## Hardware Requirements

- YKUSH3 USB switching device
- USB connection to host computer
- Appropriate drivers (typically plug-and-play on most systems)

### Supported Operating Systems

- Linux
- macOS
- Windows

## Troubleshooting

### Device Not Found

```bash
# Check if device is connected and recognized
lsusb | grep "04d8:f11b"  # Linux
system_profiler SPUSBDataType | grep -A5 YKUSH  # macOS
```

### Permission Issues (Linux)

Add udev rule to allow non-root access:

```bash
# Create udev rule file
sudo tee /etc/udev/rules.d/99-ykush3.rules << EOF
SUBSYSTEM=="usb", ATTR{idVendor}=="04d8", ATTR{idProduct}=="f11b", MODE="0666"
EOF

# Reload udev rules
sudo udevadm control --reload-rules
sudo udevadm trigger
```

### Multiple Devices

When using multiple YKUSH3 devices, always use `NewWithSerial()` to connect to specific devices:

```go
// List devices first to get serial numbers
devices, err := ykush3.ListDevices()
if err != nil {
    log.Fatal(err)
}

// Connect to specific device
ykush, err := ykush3.NewWithSerial(devices[0].SerialNbr)
if err != nil {
    log.Fatal(err)
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Related Projects

- [YKUSH Command Line Tool](https://github.com/Yepkit/ykush) - Official command line utility
- [YKUSH Python Library](https://github.com/Yepkit/pykush) - Python implementation

---

**Made with ❤️ at the [@recursecenter](https://www.recurse.com/)**
