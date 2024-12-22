/*
 *
 * Copyright 2024 tofuutils authors.
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
 *
 */

package asdfparser

import (
	"bufio"
	"errors"
	"io/fs"
	"os"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/tofuutils/tenv/v3/config"
	"github.com/tofuutils/tenv/v3/config/cmdconst"
	"github.com/tofuutils/tenv/v3/pkg/loghelper"
	"github.com/tofuutils/tenv/v3/versionmanager/semantic/types"
)

const ToolFileName = ".tool-versions"

func RetrieveTofuVersion(filePath string, conf *config.Config) (string, error) {
	return retrieveVersionFromToolFile(filePath, cmdconst.TofuName, conf)
}

func RetrieveTfVersion(filePath string, conf *config.Config) (string, error) {
	return retrieveVersionFromToolFile(filePath, cmdconst.TerraformName, conf)
}

func RetrieveTgVersion(filePath string, conf *config.Config) (string, error) {
	return retrieveVersionFromToolFile(filePath, cmdconst.TerragruntName, conf)
}

func RetrieveAtmosVersion(filePath string, conf *config.Config) (string, error) {
	return retrieveVersionFromToolFile(filePath, cmdconst.AtmosName, conf)
}

func retrieveVersionFromToolFile(filePath, toolName string, conf *config.Config) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		conf.Displayer.Log(loghelper.LevelWarnOrDebug(errors.Is(err, fs.ErrNotExist)), "Failed to open tool file", loghelper.Error, err)

		return "", nil
	}
	defer file.Close()

	resolvedVersion := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		trimmedLine := strings.TrimSpace(scanner.Text())

		if trimmedLine == "" || trimmedLine[0] == '#' {
			continue
		}

		name, remain, found := strings.Cut(trimmedLine, " ")
		if found && name == toolName { // no need to trim name (already done (left by TrimSpace, right by Cut)
			version, _, _ := strings.Cut(remain, "#")
			resolvedVersion = strings.TrimSpace(version)

			break
		}
	}

	if err := scanner.Err(); err != nil {
		conf.Displayer.Log(hclog.Warn, "Failed to parse tool file", loghelper.Error, err)

		return "", nil
	}

	if resolvedVersion == "" {
		return "", nil
	}

	return types.DisplayDetectionInfo(conf.Displayer, resolvedVersion, filePath), nil
}