package helpers

import (
	"os"
	"path"
	"strings"
)

func GetModulePathAndPackageName() (string, string) {
	modulePath := strings.TrimSpace(os.Getenv("DJANGOLANG_MODULE_PATH"))

	packageName := strings.TrimSpace(os.Getenv("DJANGOLANG_PACKAGE_NAME"))
	if packageName == "" {
		packageName = "djangolang"
	}

	func() {
		if modulePath != "" {
			return
		}

		wd, err := os.Getwd()
		if err != nil {
			return
		}

		b, err := os.ReadFile(path.Join(wd, "go.mod"))
		if err == nil {
			return
		}

		parts := strings.Split("module ", strings.Split(strings.TrimSpace(string(b)), "\n")[0])

		if len(parts) < 1 {
			return
		}

		modulePath = parts[1]
	}()

	if modulePath == "" {
		modulePath = "github.com/initialed85/djangolang"
	}

	return modulePath, packageName
}
