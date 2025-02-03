package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-vgo/robotgo"
)

const (
	MicrolabStartAddress = 16 * 16
	MicrolabNext         = ","
	MicrolabSave         = "."
	MicrolabByteMode     = "0"
	MicrolabRunMode      = "2"
	KeyPressDelay        = 110 * time.Millisecond
)

func main() {
	microlabPID, err := findMicrolabPID()

	if err != nil {
		fmt.Println(err.Error())
		time.Sleep(time.Second * 5)

		return
	}

	listingInput, err := inputEmuListing()
	if err != nil {
		fmt.Println(err.Error())

		return
	}

	hexValues := parseEmuListing(listingInput)

	emulateKeyPressing(hexValues, microlabPID)
}

func inputEmuListing() (string, error) {
	p := tea.NewProgram(initTextarea("Input here...", "Input emu8086 code listing"))
	m, err := p.Run()

	if err != nil {
		fmt.Println(err.Error())
	}
	if finalModel, ok := m.(textAreaModel); ok {
		return finalModel.content, nil
	}

	return "", errors.New("invalid model field")
}

func findMicrolabPID() (int, error) {
	pids, err := robotgo.FindIds("microlab")

	if err != nil {
		panic("Error while trying to find microlab window")
	}

	pids = filterPids(pids)

	if len(pids) < 1 {
		return -1, errors.New("cant find microlab window, please open microlab window")
	}

	return pids[0], nil
}

func emulateKeyPressing(hexValues []string, microlabPID int) {
	microlabInput := generateMicrolabInput(hexValues)

	m := progressBarModel{
		progress:         progress.New(progress.WithDefaultGradient()),
		incrementPercent: (4.0 / float64(len(microlabInput))),
		delay:            KeyPressDelay * 4,
	}

	go func() {
		if _, err := tea.NewProgram(m).Run(); err != nil {
			fmt.Println("Oh no!", err)
			os.Exit(1)
		}
	}()

	for _, char := range microlabInput {
		key := string(char)
		if key == MicrolabSave {
			key = robotgo.Enter
		}
		time.Sleep(KeyPressDelay)
		robotgo.KeyTap(key, microlabPID)
	}
}

func generateMicrolabInput(hexValues []string) string {
	startAddress := fmt.Sprintf("%04X", MicrolabStartAddress)
	endAddress := fmt.Sprintf("%04X", MicrolabStartAddress+len(hexValues))

	var result strings.Builder

	result.WriteString(MicrolabByteMode)
	result.WriteString(startAddress)
	result.WriteString(MicrolabNext)
	result.WriteString(strings.Join(hexValues, MicrolabNext))
	result.WriteString(MicrolabSave)

	result.WriteString(MicrolabRunMode)
	result.WriteString(startAddress)
	result.WriteString(MicrolabNext)
	result.WriteString(endAddress)
	result.WriteString(MicrolabSave)

	return result.String()
}

func filterPids(pids []int) []int {
	result := []int{}
	for _, pid := range pids {
		if pid != robotgo.GetPid() {
			result = append(result, pid)
		}
	}

	return result
}
