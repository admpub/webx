/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package common

import (
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/admpub/errors"
	"github.com/admpub/log"
	"github.com/admpub/sonyflake"
	"github.com/webx-top/com"
)

var (
	sonyFlakeInstances = map[uint16]*sonyflake.Sonyflake{}
	sonyFlakeLock      = &sync.RWMutex{}
	SonyflakeStartDate = `2018-09-01 08:08:08`
	defaultMachineID   uint16
)

var (
	ErrInvalidIPAddress = errors.New("Invalid IP address")
	ErrNotSet           = errors.New("Not set")
)

func init() {
	var err error
	defaultMachineID, err = ParseMachineIDFromEnvVar()
	if err != nil && !errors.Is(err, ErrNotSet) {
		log.Errorf(`failed to ParseMachineIDFromEnvVar(): %v`, err)
	}
	SonyflakeInit(defaultMachineID)
}

func ParseMachineIDFromEnvVar() (uint16, error) {
	machineIDStr := com.Getenv(`SONY_FLAKE_MACHINE_ID`)
	machineIDStr = strings.TrimSpace(machineIDStr)
	if len(machineIDStr) == 0 {
		return 0, ErrNotSet
	}
	if com.StrIsNumeric(machineIDStr) {
		n, err := strconv.ParseUint(machineIDStr, 10, 16)
		return uint16(n), err
	}
	return IPv4ToMachineID(machineIDStr)
}

func IPv4ToMachineID(ipv4 string) (uint16, error) {
	ip := net.ParseIP(ipv4)
	if ip == nil {
		return 0, ErrInvalidIPAddress
	}
	ipBytes := ip.To4()
	return uint16(ipBytes[2])<<8 + uint16(ipBytes[3]), nil
}

// NewSonyflake 19位
func NewSonyflake(startDate string, machineIDs ...uint16) (*sonyflake.Sonyflake, error) {
	if !strings.Contains(startDate, ` `) {
		startDate += ` 00:00:00`
	}
	startTime, err := time.ParseInLocation(`2006-01-02 15:04:05`, startDate, time.Local)
	if err != nil {
		return nil, err
	}
	machineID := defaultMachineID
	if len(machineIDs) > 0 {
		machineID = machineIDs[0]
	}
	st := sonyflake.Settings{
		StartTime: startTime,
		MachineID: func() (uint16, error) {
			return machineID, nil
		},
		CheckMachineID: func(id uint16) bool {
			return machineID == id
		},
	}
	return sonyflake.NewSonyflake(st), err
}

func SonyflakeInit(machineIDs ...uint16) *sonyflake.Sonyflake {
	sonyFlake, err := SetSonyflake(SonyflakeStartDate, machineIDs...)
	if err != nil {
		panic(err)
	}
	return sonyFlake
}

func SetSonyflake(startDate string, machineIDs ...uint16) (sonyFlake *sonyflake.Sonyflake, err error) {
	sonyFlake, err = NewSonyflake(startDate, machineIDs...)
	if err != nil {
		return nil, err
	}
	machineID := defaultMachineID
	if len(machineIDs) > 0 {
		machineID = machineIDs[0]
	}
	sonyFlakeLock.Lock()
	sonyFlakeInstances[machineID] = sonyFlake
	sonyFlakeLock.Unlock()
	return sonyFlake, err
}

func UniqueID(machineIDs ...uint16) (string, error) {
	id, err := NextID(machineIDs...)
	if err != nil {
		return ``, err
	}
	return strconv.FormatUint(id, 10), nil
}

func NextID(machineIDs ...uint16) (uint64, error) {
	machineID := defaultMachineID
	if len(machineIDs) > 0 {
		machineID = machineIDs[0]
	}
	sonyFlakeLock.RLock()
	sonyFlake, ok := sonyFlakeInstances[machineID]
	sonyFlakeLock.RUnlock()
	if !ok || sonyFlake == nil {
		var err error
		sonyFlake, err = SetSonyflake(SonyflakeStartDate, machineIDs...)
		if err != nil {
			return 0, err
		}
	}
	return sonyFlake.NextID()
}
