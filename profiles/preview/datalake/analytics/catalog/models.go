// +build go1.9

// Copyright 2017 Microsoft Corporation
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

// This code was auto-generated by:
// github.com/Azure/azure-sdk-for-go/tools/profileBuilder
// commit ID: 714052a3359963ba46703fe033cf9dd4838bea11

package catalog

import original "github.com/Azure/azure-sdk-for-go/services/datalake/analytics/2016-11-01-preview/catalog"

const (
	DefaultAdlaCatalogDNSSuffix = original.DefaultAdlaCatalogDNSSuffix
)

type ManagementClient = original.ManagementClient
type GroupClient = original.GroupClient
type FileType = original.FileType

const (
	Assembly FileType = original.Assembly
	Nodeploy FileType = original.Nodeploy
	Resource FileType = original.Resource
)

type DataLakeAnalyticsCatalogCredentialCreateParameters = original.DataLakeAnalyticsCatalogCredentialCreateParameters
type DataLakeAnalyticsCatalogCredentialDeleteParameters = original.DataLakeAnalyticsCatalogCredentialDeleteParameters
type DataLakeAnalyticsCatalogCredentialUpdateParameters = original.DataLakeAnalyticsCatalogCredentialUpdateParameters
type DataLakeAnalyticsCatalogSecretCreateOrUpdateParameters = original.DataLakeAnalyticsCatalogSecretCreateOrUpdateParameters
type DdlName = original.DdlName
type EntityID = original.EntityID
type ExternalTable = original.ExternalTable
type Item = original.Item
type ItemList = original.ItemList
type TypeFieldInfo = original.TypeFieldInfo
type USQLAssembly = original.USQLAssembly
type USQLAssemblyClr = original.USQLAssemblyClr
type USQLAssemblyDependencyInfo = original.USQLAssemblyDependencyInfo
type USQLAssemblyFileInfo = original.USQLAssemblyFileInfo
type USQLAssemblyList = original.USQLAssemblyList
type USQLCredential = original.USQLCredential
type USQLCredentialList = original.USQLCredentialList
type USQLDatabase = original.USQLDatabase
type USQLDatabaseList = original.USQLDatabaseList
type USQLDirectedColumn = original.USQLDirectedColumn
type USQLDistributionInfo = original.USQLDistributionInfo
type USQLExternalDataSource = original.USQLExternalDataSource
type USQLExternalDataSourceList = original.USQLExternalDataSourceList
type USQLIndex = original.USQLIndex
type USQLPackage = original.USQLPackage
type USQLPackageList = original.USQLPackageList
type USQLProcedure = original.USQLProcedure
type USQLProcedureList = original.USQLProcedureList
type USQLSchema = original.USQLSchema
type USQLSchemaList = original.USQLSchemaList
type USQLSecret = original.USQLSecret
type USQLTable = original.USQLTable
type USQLTableColumn = original.USQLTableColumn
type USQLTableList = original.USQLTableList
type USQLTablePartition = original.USQLTablePartition
type USQLTablePartitionList = original.USQLTablePartitionList
type USQLTableStatistics = original.USQLTableStatistics
type USQLTableStatisticsList = original.USQLTableStatisticsList
type USQLTableType = original.USQLTableType
type USQLTableTypeList = original.USQLTableTypeList
type USQLTableValuedFunction = original.USQLTableValuedFunction
type USQLTableValuedFunctionList = original.USQLTableValuedFunctionList
type USQLType = original.USQLType
type USQLTypeList = original.USQLTypeList
type USQLView = original.USQLView
type USQLViewList = original.USQLViewList

func UserAgent() string {
	return original.UserAgent()
}
func Version() string {
	return original.Version()
}
func New() ManagementClient {
	return original.New()
}
func NewWithoutDefaults(adlaCatalogDNSSuffix string) ManagementClient {
	return original.NewWithoutDefaults(adlaCatalogDNSSuffix)
}
func NewGroupClient() GroupClient {
	return original.NewGroupClient()
}
