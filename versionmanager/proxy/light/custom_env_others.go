//go:build !windows

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

package lightproxy

import (
	"os"
	"os/exec"

	"github.com/tofuutils/tenv/v4/config/envname"
	"github.com/tofuutils/tenv/v4/pkg/tty"
)

const changeDefaultDetach = envname.TenvDetachedProxyDefault + "=true"

func updateDefaultDetachInCmdEnv(cmd *exec.Cmd) bool {
	if tty.Detect() {
		cmd.Env = append(os.Environ(), changeDefaultDetach)

		return true
	}

	return false
}