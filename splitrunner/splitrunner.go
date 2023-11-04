package splitrunner

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/tools/go/packages"

	"golangci-test/splitter"
)

var ErrPackageNotFound = fmt.Errorf("package not found")
var ErrSplitContainsNotExistingPackage = fmt.Errorf("split contains not existing package")

// RunSplit runs tests for given split of packages.
// Execution context must be set to the directory containing go.mod file (containing packages from the split)
func RunSplit(split splitter.Split, jsonOut string) error {
	args := append([]string{"test", "-coverprofile=" + jsonOut}, split...)
	if jsonOut != "" {
		args = append(args, "-json", jsonOut)
	}
	cmd := exec.Command("go", args...)
	out, err := cmd.Output()
	fmt.Println(string(out))
	if err != nil {
		fmt.Printf("WError running tests: %s\n", err.Error())
		return err
	}
	return nil
}

// ValidateSplits validates that splits contain all packages from the srcPath
// and that splits do not contain not existing packages
func ValidateSplits(splits []splitter.Split, srcPath string) error {
	packagesInSplits := map[string]bool{}
	for _, split := range splits {
		for _, pkg := range split {
			packagesInSplits[pkg] = true
		}
	}
	packagesInSrc := extractDeduplicatedMainPkgsPaths(getPackagesContainingTests(srcPath))
	for pkg, _ := range packagesInSrc {
		if !packagesInSplits[pkg] {
			return ErrPackageNotFound
		}
	}
	var errs error
	for pkg, _ := range packagesInSplits {
		if !packagesInSrc[pkg] {
			errs = errors.Join(errs, ErrSplitContainsNotExistingPackage)
		}
	}
	return errs
}

func getPackagesContainingTests(srcPath string) []*packages.Package {
	cfg := &packages.Config{
		Mode:  packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles,
		Tests: true, // Set this to true to include test packages
	}

	// Load the packages from the current directory, including tests.
	pkgs, err := packages.Load(cfg, srcPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load packages: %s\n", err)
		os.Exit(1)
	}

	// Filter out packages that do not contain any test files.
	testPkgs := filterTestPackages(pkgs)
	return testPkgs
}

func filterTestPackages(pkgs []*packages.Package) []*packages.Package {
	testPkgs := []*packages.Package{}
	for _, pkg := range pkgs {
		if containsTestFiles(pkg.GoFiles) {
			testPkgs = append(testPkgs, pkg)
		}
	}
	return testPkgs
}

func containsTestFiles(files []string) bool {
	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") {
			return true
		}
	}
	return false
}

func extractDeduplicatedMainPkgsPaths(pkgs []*packages.Package) map[string]bool {
	mainPkgsPaths := map[string]bool{}
	for _, pkg := range pkgs {
		mainPkgPath := strings.TrimSuffix(pkg.PkgPath, "_test")
		mainPkgsPaths[mainPkgPath] = true
	}
	return mainPkgsPaths
}
