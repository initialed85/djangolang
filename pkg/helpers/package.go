package helpers

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

func GetModulePathAndPackageName() (string, string) {
	modulePath := GetEnvironmentVariable("DJANGOLANG_MODULE_PATH")

	packageName := GetEnvironmentVariable("DJANGOLANG_PACKAGE_NAME")

	func() {
		if modulePath != "" {
			return
		}

		wd, err := os.Getwd()
		if err != nil {
			return
		}

		b, err := os.ReadFile(path.Join(wd, "go.mod"))
		if err != nil {
			return
		}

		parts := strings.Split(strings.Split(strings.TrimSpace(string(b)), "\n")[0], "module ")

		if len(parts) < 2 {
			log.Printf("err: %v", fmt.Errorf("not enough parts in %v", parts))
			return
		}

		modulePath = parts[1]
		log.Printf("DJANGOLANG_MODULE_PATH empty or unset; defaulted to %v (found at %v)", modulePath, path.Join(wd, "go.mod"))

		modulePathParts := strings.Split(modulePath, "/")

		packageName = modulePathParts[len(modulePathParts)-1]
		log.Printf("DJANGOLANG_PACKAGE_NAME empty or unset; defaulted to %v (found at %v)", packageName, path.Join(wd, "go.mod"))
	}()

	if modulePath == "" {
		modulePath = "github.com/initialed85/djangolang"
		log.Printf("DJANGOLANG_MODULE_PATH empty or unset; defaulted to %v", modulePath)
	}

	if packageName == "" {
		packageName = "djangolang"
		log.Printf("DJANGOLANG_PACKAGE_NAME empty or unset; defaulted to %v", packageName)
	}

	return modulePath, packageName
}
