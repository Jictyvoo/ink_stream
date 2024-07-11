package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

const kccScriptPath = "/opt/kcc/kcc-c2e.py"

func runKccScript(extractDir, outputDir string) {
	lastFolderName := filepath.Base(extractDir)
	outputFolder := filepath.Join(outputDir, lastFolderName)

	cmd := exec.Command(
		kccScriptPath,
		"--profile", string(KDeviceKindle11),
		"--manga-style", "-q", "--upscale",
		"--format", `"`+string(FormatMOBIEPUB)+`"`,
		"--batchsplit", "2",
		"--output", outputFolder,
		extractDir,
	)

	fmt.Printf("Running command: %s\n\n", cmd.String())
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Printf("Failed to run kcc-c2e.py on %s: %v\n%s", extractDir, err, string(output))
	} else {
		fmt.Printf("Successfully processed %s\n", extractDir)
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
