/*
 * Copyright 2019 ObjectBox Ltd. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package generator

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func copyFile(sourceFile, targetFile string, permsOverride os.FileMode) error {
	data, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		return err
	}

	// copy permissions either from the existing target file or from the source file
	var perm os.FileMode = permsOverride
	if perm == 0 {
		if info, _ := os.Stat(targetFile); info != nil {
			perm = info.Mode()
		} else if info, err := os.Stat(sourceFile); info != nil {
			perm = info.Mode()
		} else {
			return err
		}
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

func copyDirectory(sourceDir, targetDir string, dirPerms, filePerms os.FileMode) error {
	if err := os.MkdirAll(targetDir, dirPerms); err != nil {
		return err
	}

	entries, err := ioutil.ReadDir(sourceDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		sourcePath := filepath.Join(sourceDir, entry.Name())
		targetPath := filepath.Join(targetDir, entry.Name())

		info, err := os.Stat(sourcePath)
		if err != nil {
			return err
		}

		if info.IsDir() {
			if err := copyDirectory(sourcePath, targetPath, dirPerms, filePerms); err != nil {
				return err
			}
		} else if info.Mode().IsRegular() {
			if err := copyFile(sourcePath, targetPath, filePerms); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("not a regular file or directory: %s", sourcePath)
		}
	}
	return nil
}
