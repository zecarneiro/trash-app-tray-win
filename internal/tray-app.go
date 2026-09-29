package internal

import (
	"fmt"
	"golangutils/pkg/conv"
	"golangutils/pkg/enums"
	"golangutils/pkg/exe"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"
	"golangutils/pkg/timec"
	"syscall"

	"github.com/energye/systray"
)

const (
	APP_NAME     = "Recycle Bin"
	MONITOR_TIME = 2
)

var (
	isSystrayCreated bool
)

func isRecycleBinFull() bool {
	cmdStr := "(New-Object -ComObject Shell.Application).Namespace(0xa).Items().Count"
	output, err := exe.Exec(models.Command{Cmd: cmdStr, Verbose: false, UseShell: true, ShellToUse: enums.PowerShell})
	if err != nil {
		return false
	}
	count, err := conv.StringToInt(output)
	if err != nil {
		return false
	}
	return count > 0
}

func updateState() {
	full := isRecycleBinFull()
	systray.SetIcon(generateIconData(full))
	if full {
		systray.SetTooltip("Recycle Bin (Full)")
	} else {
		systray.SetTooltip("Recycle Bin (Empty)")
	}
}

func showMenu(menu systray.IMenu) {
	if menu != nil {
		menu.ShowMenu()
	}
}

func buildTrayApp() {
	systray.AddMenuItem("Open Recycle Bin", "Open Recycle Bin").Click(func() {
		exe.ExecRealTime(models.Command{Cmd: "explorer.exe shell:RecycleBinFolder", Verbose: false, UseShell: true, IsAsync: true})
	})
	systray.AddMenuItem("Empty Recycle Bin", "Empty Recycle Bin").Click(func() {
		// Call Windows API (Shell32.dll) to clean Recycle Bin
		shell32 := syscall.NewLazyDLL("shell32.dll")
		shEmptyRecycleBin := shell32.NewProc("SHEmptyRecycleBinW")
		_, _, _ = shEmptyRecycleBin.Call(0, 0, 0)
		updateState()
	})
	systray.AddSeparator()
	buildThemeMenu()
	systray.AddMenuItem("Exit", "Exit of the application").Click(func() {
		systray.Quit()
	})
	if !isSystrayCreated {
		systray.SetTitle(APP_NAME)
		systray.SetTooltip(APP_NAME)
		systray.SetOnClick(showMenu)
		systray.SetOnRClick(showMenu)
		isSystrayCreated = true
	}
	// Periodic monitoring every by Monitor time
	go func() {
		for {
			updateState()
			timec.Sleep(MONITOR_TIME)
		}
	}()
}

func Start() {
	if !platform.IsWindows() {
		logic.ProcessError(fmt.Errorf("Invalid Platform."))
	}
	loadConfigurations()
	systray.Run(buildTrayApp, nil)
}
