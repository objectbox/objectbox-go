package generator

import (
	"io/ioutil"
	"os"
)

func copyFile(sourceFile, targetFile string) error {
	data, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		return err
	}

	// copy permissions either from the existing target file or from the source file
	var perm os.FileMode
	if info, _ := os.Stat(targetFile); info != nil {
		perm = info.Mode()
	} else if info, err := os.Stat(sourceFile); info != nil {
		perm = info.Mode()
	} else {
		return err
	}

	err = ioutil.WriteFile(targetFile, data, perm)
	if err != nil {
		return err
	}

	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
