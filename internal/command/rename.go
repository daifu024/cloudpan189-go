// Copyright (c) 2020 tickstep.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package command

import (
	"fmt"
	"github.com/tickstep/cloudpan189-api/cloudpan"
	"github.com/tickstep/cloudpan189-api/cloudpan/apierror"
	"github.com/tickstep/cloudpan189-api/cloudpan/apiutil"
	"path"
	"strings"
)

func RunRename(familyId int64, oldName string, newName string) {
	if oldName == "" {
		fmt.Println("请指定命名文件")
		return
	}
	if newName == "" {
		fmt.Println("请指定文件新名称")
		return
	}
	activeUser := GetActiveUser()
	oldName = activeUser.PathJoin(familyId, strings.TrimSpace(oldName))
	newName = activeUser.PathJoin(familyId, strings.TrimSpace(newName))
	if path.Dir(oldName) != path.Dir(newName) {
		fmt.Println("只能命名同一个目录的文件")
		return
	}
	if !apiutil.CheckFileNameValid(path.Base(newName)) {
		fmt.Println("文件名不能包含特殊字符：" + apiutil.FileNameSpecialChars)
		return
	}

	fileId := ""
	r, err := GetActivePanClient().AppFileInfoByPath(familyId, activeUser.PathJoin(familyId, oldName))
	if err != nil {
		fmt.Printf("原文件不存在： %s, %s\n", oldName, err)
		return
	}
	fileId = r.FileId

	var b *cloudpan.AppFileEntity
	var e *apierror.ApiError
	if IsFamilyCloud(familyId) {
		b, e = activeUser.PanClient().AppFamilyRenameFile(familyId, fileId, path.Base(newName))
	} else {
		b, e = activeUser.PanClient().AppRenameFile(fileId, path.Base(newName))
	}
	if e != nil {
		fmt.Println(e.Err)
		return
	}
	if b == nil {
		fmt.Println("重命名文件失败")
		return
	}
	fmt.Printf("重命名文件成功：%s -> %s\n", path.Base(oldName), path.Base(newName))
}
