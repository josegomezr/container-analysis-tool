package scanner

import (
	"context"
	"fmt"
	"github.com/josegomezr/container-analysis-tools/pkg/exec_factories"
	// couldn't really use trivy from within :sweat:
	// "github.com/aquasecurity/trivy/pkg/flag"
	// "github.com/aquasecurity/trivy/pkg/types"
	// "github.com/aquasecurity/trivy/pkg/commands/artifact"
)

func ScanImage(ctx context.Context, input string) (string, error) {
	fmt.Printf("Scanning \n")
	cmd := exec_factories.DefaultCmdFactory("trivy",
		"image",
		"--input",
		input,
		"--format",
		"json",
	)

	fmt.Printf("debug: Trivy invokation: %s %v \n", cmd.Path, cmd.Args)
	scan, err := cmd.Output()
	if err != nil {
		fmt.Printf("debug: error on trivy invocation: %v \n", err)
		return "", err
	}
	return string(scan), nil
}
