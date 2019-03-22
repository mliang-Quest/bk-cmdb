/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package extensions

import (
	"configcenter/src/common"
	"configcenter/src/common/condition"
	"context"
	"fmt"
	"net/http"
	
	"configcenter/src/auth/meta"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (am *AuthManager) collectAssociationTypesByIDs(ctx context.Context, header http.Header, ids ...int64) ([]metadata.AssociationKind, error) {
	// get model by objID
	cond := condition.CreateCondition().Field(common.BKFieldID).In(ids)
	queryCond := &metadata.QueryCondition{Condition: cond.ToMapStr()}
	resp, err := am.clientSet.CoreService().Instance().ReadInstance(ctx, header, common.BKTableNameAsstDes, queryCond)
	if err != nil {
		return nil, fmt.Errorf("get association types by id: %+v failed, err: %+v", ids, err)
	}
	if len(resp.Data.Info) == 0 {
		return nil, fmt.Errorf("get association types by id: %+v failed, not found", ids)
	}
	if len(resp.Data.Info) != len(ids) {
		return nil, fmt.Errorf("get association types by id: %+v failed, get %d, expect %d", ids, len(resp.Data.Info), len(ids))
	}

	aks := make([]metadata.AssociationKind, 0)
	for _, item := range resp.Data.Info {
		ak := metadata.AssociationKind{}
		ak.Parse(item)
		aks = append(aks, ak)
	}
	return aks, nil
}
func (am *AuthManager) makeResourceByAssociationType(ctx context.Context, header http.Header, action meta.Action, aks ...metadata.AssociationKind) ([]meta.ResourceAttribute, error) {
	resources := make([]meta.ResourceAttribute, 0)
	for _, ak := range aks {
		resource := meta.ResourceAttribute{
			Basic: meta.Basic{
				Type:       meta.AssociationType,
				Name:       ak.AssociationKindID,
				InstanceID: ak.ID,
			},
			SupplierAccount: util.GetOwnerID(header),
		}
		resources = append(resources, resource)
	}
	return resources, nil
}

func (am *AuthManager) RegisterAssociationType(ctx context.Context, header http.Header, aks ...metadata.AssociationKind) error {
	resources, err := am.makeResourceByAssociationType(ctx, header, meta.EmptyAction, aks...)
	if err != nil {
		return fmt.Errorf("make auth resource from association type failed, err: %+v", err)
	}

	return am.Authorize.RegisterResource(ctx, resources...)
}

func (am *AuthManager) RegisterAssociationTypeByID(ctx context.Context, header http.Header, ids ...int64) error {
	aks, err := am.collectAssociationTypesByIDs(ctx, header, ids...)
	if err != nil {
		return fmt.Errorf("get asssociation type by id failed, err: %+v", err)
	}

	return am.RegisterAssociationType(ctx, header, aks...)
}

func (am *AuthManager) UpdateAssociationTypeByID(ctx context.Context, header http.Header, ids ...int64) error {
	aks, err := am.collectAssociationTypesByIDs(ctx, header, ids...)
	if err != nil {
		return fmt.Errorf("get asssociation type by id failed, err: %+v", err)
	}

	resources, err := am.makeResourceByAssociationType(ctx, header, meta.EmptyAction, aks...)
	if err != nil {
		return fmt.Errorf("make auth resource from association type failed, err: %+v", err)
	}

	return am.updateResources(ctx, resources...)
}

func (am *AuthManager) DeregisterAssociationTypeByIDs(ctx context.Context, header http.Header, ids ...int64) error {
	aks, err := am.collectAssociationTypesByIDs(ctx, header, ids...)
	if err != nil {
		return fmt.Errorf("get asssociation type by id failed, err: %+v", err)
	}

	resources, err := am.makeResourceByAssociationType(ctx, header, meta.EmptyAction, aks...)
	if err != nil {
		return fmt.Errorf("make auth resource from association type failed, err: %+v", err)
	}

	return am.Authorize.DeregisterResource(ctx, resources...)
}
