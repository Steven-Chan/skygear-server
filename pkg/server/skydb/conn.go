// Copyright 2015-present Oursky Ltd.
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

//go:generate mockgen -package skydb -source=conn.go -destination=mock_conn_test.go

package skydb

import (
	"errors"
	"time"
)

// ErrUserDuplicated is returned by Conn.CreateAuth when
// the AuthInfo to be created has the same ID/username in the current container
var ErrUserDuplicated = errors.New("skydb: duplicated UserID or Username")

// ErrUserNotFound is returned by Conn.GetAuth, Conn.UpdateAuth and
// Conn.DeleteAuth when the AuthInfo's ID is not found
// in the current container
var ErrUserNotFound = errors.New("skydb: AuthInfo ID not found")

var ErrRoleUpdatesFailed = errors.New("skydb: Update of user roles failed")

// ErrDeviceNotFound is returned by Conn.GetDevice, Conn.DeleteDevice,
// Conn.DeleteDevicesByToken and Conn.DeleteEmptyDevicesByTime, if the desired Device
// cannot be found in the current container
var ErrDeviceNotFound = errors.New("skydb: Specific device not found")

// ErrDatabaseIsReadOnly is returned by skydb.Database if the requested
// operation modifies the database and the database is readonly.
var ErrDatabaseIsReadOnly = errors.New("skydb: database is read only")

// ZeroTime represent a zero time.Time. It is used in DeleteDevicesByToken and
// DeleteEmptyDevicesByTime to signify a Delete without time constraint.
var ZeroTime = time.Time{}

// DBHookFunc specifies the interface of a database hook function
type DBHookFunc func(Database, *Record, RecordHookEvent)

// QueryConfig provides optional parameters for queries.
// result is unlimited if Limit=0
type QueryConfig struct {
	Limit  uint64
	Offset uint64
}

// Conn encapsulates the interface of an Skygear Server connection to a container.
//go:generate mockgen -destination=mock_skydb/mock_conn.go github.com/skygeario/skygear-server/pkg/server/skydb Conn
type Conn interface {
	// CRUD of AuthInfo, smell like a bad design to attach these onto
	// a Conn, but looks very convenient to user.

	// CreateAuth creates a new AuthInfo in the container
	// this Conn associated to.
	CreateAuth(authinfo *AuthInfo) error

	// GetAuth fetches the AuthInfo with supplied ID in the container and
	// fills in the supplied AuthInfo with the result.
	//
	// GetAuth returns ErrUserNotFound if no AuthInfo exists
	// for the supplied ID.
	GetAuth(id string, authinfo *AuthInfo) error

	// GetAuthByPrincipalID fetches the AuthInfo with supplied principal ID in the
	// container and fills in the supplied AuthInfo with the result.
	//
	// Principal ID is an ID of an authenticated principal with such
	// authentication provided by AuthProvider.
	//
	// GetAuthByPrincipalID returns ErrUserNotFound if no AuthInfo exists
	// for the supplied principal ID.
	GetAuthByPrincipalID(principalID string, authinfo *AuthInfo) error

	// UpdateAuth updates an existing AuthInfo matched by the ID field.
	//
	// UpdateAuth returns ErrUserNotFound if such AuthInfo does not
	// exist in the container.
	UpdateAuth(authinfo *AuthInfo) error

	// DeleteAuth removes AuthInfo with the supplied ID in the container.
	//
	// DeleteAuth returns ErrUserNotFound if such AuthInfo does not
	// exist in the container.
	DeleteAuth(id string) error

	// GetAdminRoles return the current admine roles
	GetAdminRoles() ([]string, error)

	// SetAdminRoles accepts array of role, the order will be
	SetAdminRoles(roles []string) error

	// GetDefaultRoles return the current default roles
	GetDefaultRoles() ([]string, error)

	// SetDefaultRoles accepts array of roles, the supplied roles will assigned
	// to newly created user CreateAuth
	SetDefaultRoles(roles []string) error

	// AssignRoles accepts array of roles and userID, the supplied roles will
	// be assigned to all passed in users
	AssignRoles(userIDs []string, roles []string) error

	// RevokeRoles accepts array of roles and userID, the supplied roles will
	// be revoked from all passed in users
	RevokeRoles(userIDs []string, roles []string) error

	// SetRecordAccess sets default record access of a specific type
	SetRecordAccess(recordType string, acl RecordACL) error

	// SetRecordDefaultAccess sets default record access of a specific type
	SetRecordDefaultAccess(recordType string, acl RecordACL) error

	// GetRecordAccess returns the record creation access of a specific type
	GetRecordAccess(recordType string) (RecordACL, error)

	// GetRecordDefaultAccess returns default record access of a specific type
	GetRecordDefaultAccess(recordType string) (RecordACL, error)

	// SetRecordFieldAccess replace field ACL setting
	SetRecordFieldAccess(acl FieldACL) (err error)

	// GetRecordFieldAccess retrieve field ACL setting
	GetRecordFieldAccess() (FieldACL, error)

	// GetAsset retrieves Asset information by its name
	GetAsset(name string, asset *Asset) error

	GetAssets(names []string) ([]Asset, error)

	// SaveAsset saves an Asset information into a container to
	// be referenced by records.
	SaveAsset(asset *Asset) error

	QueryRelation(user string, name string, direction string, config QueryConfig) []AuthInfo
	QueryRelationCount(user string, name string, direction string) (uint64, error)
	AddRelation(user string, name string, targetUser string) error
	RemoveRelation(user string, name string, targetUser string) error

	GetDevice(id string, device *Device) error

	// QueryDevicesByUser queries the Device database which are registered
	// by the specified user.
	QueryDevicesByUser(user string) ([]Device, error)
	QueryDevicesByUserAndTopic(user, topic string) ([]Device, error)
	SaveDevice(device *Device) error
	DeleteDevice(id string) error

	// DeleteDevicesByToken deletes device where its Token == token and
	// LastRegisteredAt < t. If t == ZeroTime, LastRegisteredAt is not considered.
	//
	// If such device does not exist, ErrDeviceNotFound is returned.
	DeleteDevicesByToken(token string, t time.Time) error

	// DeleteEmptyDevicesByTime deletes device where Token is empty and
	// LastRegisteredAt < t. If t == ZeroTime, LastRegisteredAt is not considered.
	//
	// If such device does not exist, ErrDeviceNotFound is returned.
	DeleteEmptyDevicesByTime(t time.Time) error

	PublicDB() Database
	PrivateDB(userKey string) Database
	UnionDB() Database

	// Subscribe registers the specified recordEventChan to receive
	// RecordEvent from the Conn implementation
	Subscribe(recordEventChan chan RecordEvent) error

	Close() error
}

// AccessModel indicates the type of access control model while db query.
//go:generate stringer -type=AccessModel
type AccessModel int

// RoleBasedAccess is tranditional Role based Access Control
// RelationBasedAccess is Access Control determine by the user-user relation
// between creator and accessor
const (
	RoleBasedAccess AccessModel = iota + 1
	RelationBasedAccess
)

// RecordHookEvent indicates the type of record event that triggered
// the hook
type RecordHookEvent int

// See the definition of RecordHookEvent
const (
	RecordCreated RecordHookEvent = iota + 1
	RecordUpdated
	RecordDeleted
)

// RecordEvent describes a change event on Record which is either
// Created, Updated or Deleted.
//
// For RecordCreated or RecordUpdated event, Record is the newly
// created / updated Record. For RecordDeleted, Record is the Record
// being deleted.
type RecordEvent struct {
	Record *Record
	Event  RecordHookEvent
}
