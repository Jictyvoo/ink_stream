package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/Jictyvoo/ink_stream/internal/utils"
)

const kccScriptPath = "/opt/kcc/kcc-c2e.py"

func normalizeOutputName(name string) string {
	var builder strings.Builder
	for _, character := range name {
		var (
			isUpper  = unicode.IsUpper(character)
			isLower  = unicode.IsLower(character)
			isNumber = unicode.IsNumber(character)
		)
		if unicode.IsSpace(character) {
			builder.WriteRune('_')
		} else if isUpper || isLower || isNumber {
			builder.WriteRune(character)
		}
	}

	return builder.String()
}

func runKccScript(extractDir, outputDir string) {
	lastFolderName := filepath.Base(extractDir)
	outputFolder := filepath.Join(outputDir, normalizeOutputName(lastFolderName))
	if err := utils.CreateDirIfNotExist(outputFolder); err != nil {
		log.Fatalf("Failed to create kcc script directory %s: %s", outputDir, err)
	}

	cmd := exec.Command(
		kccScriptPath,
		"--profile", string(KDeviceKindle11),
		"--manga-style", "-q", "--upscale",
		"--format", `"`+string(FormatMOBIEPUB)+`"`,
		"--batchsplit", "2",
		"--output", outputFolder,
		extractDir,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Printf("Running command: %s\n\n", cmd.String())
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to run kcc-c2e.py on %s: %v\n", extractDir, err)
	}
}

func main() {
	var sourceDir, outputDir string
	flag.StringVar(&sourceDir, "src", "", "source directory")
	flag.StringVar(&outputDir, "out", "./output", "output directory")
	flag.Parse()
	if sourceDir == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Run the kcc-c2e.py command on the extracted directory
	runKccScript(sourceDir, outputDir)
}
