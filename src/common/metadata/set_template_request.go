/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017,-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */
package metadata

import "configcenter/src/framework/core/errors"

type CreateSetTemplateOption struct {
	Name               string  `field:"name" json:"name" bson:"name"`
	ServiceTemplateIDs []int64 `field:"service_template_ids" json:"service_template_ids" bson:"service_template_ids"`
}

type UpdateSetTemplateOption struct {
	Name               string  `field:"name" json:"name" bson:"name"`
	ServiceTemplateIDs []int64 `field:"service_template_ids" json:"service_template_ids" bson:"service_template_ids"`
}

func (option UpdateSetTemplateOption) Validate() (string, error) {
	if len(option.Name) == 0 && option.ServiceTemplateIDs == nil {
		return "", errors.New("at least one update field not empty")
	}
	return "", nil
}

type SetTemplateResult struct {
	BaseResp
	Data SetTemplate `field:"data" json:"data" bson:"data"`
}
