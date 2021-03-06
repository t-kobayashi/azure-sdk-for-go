package automanageapi

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/services/preview/automanage/mgmt/2020-06-30-preview/automanage"
	"github.com/Azure/go-autorest/autorest"
)

// AccountsClientAPI contains the set of methods on the AccountsClient type.
type AccountsClientAPI interface {
	CreateOrUpdate(ctx context.Context, accountName string, resourceGroupName string, parameters automanage.Account) (result automanage.Account, err error)
	Delete(ctx context.Context, resourceGroupName string, accountName string) (result autorest.Response, err error)
	Get(ctx context.Context, accountName string, resourceGroupName string) (result automanage.Account, err error)
	ListByResourceGroup(ctx context.Context, resourceGroupName string) (result automanage.AccountList, err error)
	ListBySubscription(ctx context.Context) (result automanage.AccountList, err error)
	Update(ctx context.Context, accountName string, resourceGroupName string, parameters automanage.Account) (result automanage.Account, err error)
}

var _ AccountsClientAPI = (*automanage.AccountsClient)(nil)

// ConfigurationProfileAssignmentsClientAPI contains the set of methods on the ConfigurationProfileAssignmentsClient type.
type ConfigurationProfileAssignmentsClientAPI interface {
	CreateOrUpdate(ctx context.Context, configurationProfileAssignmentName string, parameters automanage.ConfigurationProfileAssignment, resourceGroupName string, VMName string) (result automanage.ConfigurationProfileAssignmentsCreateOrUpdateFuture, err error)
	Delete(ctx context.Context, resourceGroupName string, configurationProfileAssignmentName string, VMName string) (result autorest.Response, err error)
	Get(ctx context.Context, resourceGroupName string, configurationProfileAssignmentName string, VMName string) (result automanage.ConfigurationProfileAssignment, err error)
	List(ctx context.Context, resourceGroupName string) (result automanage.ConfigurationProfileAssignmentList, err error)
	ListBySubscription(ctx context.Context) (result automanage.ConfigurationProfileAssignmentList, err error)
}

var _ ConfigurationProfileAssignmentsClientAPI = (*automanage.ConfigurationProfileAssignmentsClient)(nil)

// ConfigurationProfilePreferencesClientAPI contains the set of methods on the ConfigurationProfilePreferencesClient type.
type ConfigurationProfilePreferencesClientAPI interface {
	CreateOrUpdate(ctx context.Context, configurationProfilePreferenceName string, resourceGroupName string, parameters automanage.ConfigurationProfilePreference) (result automanage.ConfigurationProfilePreference, err error)
	Delete(ctx context.Context, resourceGroupName string, configurationProfilePreferenceName string) (result autorest.Response, err error)
	Get(ctx context.Context, configurationProfilePreferenceName string, resourceGroupName string) (result automanage.ConfigurationProfilePreference, err error)
	ListByResourceGroup(ctx context.Context, resourceGroupName string) (result automanage.ConfigurationProfilePreferenceList, err error)
	ListBySubscription(ctx context.Context) (result automanage.ConfigurationProfilePreferenceList, err error)
	Update(ctx context.Context, configurationProfilePreferenceName string, resourceGroupName string, parameters automanage.ConfigurationProfilePreference) (result automanage.ConfigurationProfilePreference, err error)
}

var _ ConfigurationProfilePreferencesClientAPI = (*automanage.ConfigurationProfilePreferencesClient)(nil)

// OperationsClientAPI contains the set of methods on the OperationsClient type.
type OperationsClientAPI interface {
	List(ctx context.Context) (result automanage.OperationList, err error)
}

var _ OperationsClientAPI = (*automanage.OperationsClient)(nil)
