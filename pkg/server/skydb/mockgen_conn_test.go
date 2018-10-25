// Code generated by MockGen. DO NOT EDIT.
// Source: conn.go

// Package skydb is a generated GoMock package.
package skydb

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	time "time"
)

// MockConn is a mock of Conn interface
type MockConn struct {
	ctrl     *gomock.Controller
	recorder *MockConnMockRecorder
}

// MockConnMockRecorder is the mock recorder for MockConn
type MockConnMockRecorder struct {
	mock *MockConn
}

// NewMockConn creates a new mock instance
func NewMockConn(ctrl *gomock.Controller) *MockConn {
	mock := &MockConn{ctrl: ctrl}
	mock.recorder = &MockConnMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockConn) EXPECT() *MockConnMockRecorder {
	return m.recorder
}

// CreateAuth mocks base method
func (m *MockConn) CreateAuth(authinfo *AuthInfo) error {
	ret := m.ctrl.Call(m, "CreateAuth", authinfo)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateAuth indicates an expected call of CreateAuth
func (mr *MockConnMockRecorder) CreateAuth(authinfo interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAuth", reflect.TypeOf((*MockConn)(nil).CreateAuth), authinfo)
}

// GetAuth mocks base method
func (m *MockConn) GetAuth(id string, authinfo *AuthInfo) error {
	ret := m.ctrl.Call(m, "GetAuth", id, authinfo)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetAuth indicates an expected call of GetAuth
func (mr *MockConnMockRecorder) GetAuth(id, authinfo interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAuth", reflect.TypeOf((*MockConn)(nil).GetAuth), id, authinfo)
}

// GetAuthByPrincipalID mocks base method
func (m *MockConn) GetAuthByPrincipalID(principalID string, authinfo *AuthInfo) error {
	ret := m.ctrl.Call(m, "GetAuthByPrincipalID", principalID, authinfo)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetAuthByPrincipalID indicates an expected call of GetAuthByPrincipalID
func (mr *MockConnMockRecorder) GetAuthByPrincipalID(principalID, authinfo interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAuthByPrincipalID", reflect.TypeOf((*MockConn)(nil).GetAuthByPrincipalID), principalID, authinfo)
}

// UpdateAuth mocks base method
func (m *MockConn) UpdateAuth(authinfo *AuthInfo) error {
	ret := m.ctrl.Call(m, "UpdateAuth", authinfo)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateAuth indicates an expected call of UpdateAuth
func (mr *MockConnMockRecorder) UpdateAuth(authinfo interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAuth", reflect.TypeOf((*MockConn)(nil).UpdateAuth), authinfo)
}

// DeleteAuth mocks base method
func (m *MockConn) DeleteAuth(id string) error {
	ret := m.ctrl.Call(m, "DeleteAuth", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAuth indicates an expected call of DeleteAuth
func (mr *MockConnMockRecorder) DeleteAuth(id interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAuth", reflect.TypeOf((*MockConn)(nil).DeleteAuth), id)
}

// GetPasswordHistory mocks base method
func (m *MockConn) GetPasswordHistory(authID string, historySize, historyDays int) ([]PasswordHistory, error) {
	ret := m.ctrl.Call(m, "GetPasswordHistory", authID, historySize, historyDays)
	ret0, _ := ret[0].([]PasswordHistory)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPasswordHistory indicates an expected call of GetPasswordHistory
func (mr *MockConnMockRecorder) GetPasswordHistory(authID, historySize, historyDays interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPasswordHistory", reflect.TypeOf((*MockConn)(nil).GetPasswordHistory), authID, historySize, historyDays)
}

// RemovePasswordHistory mocks base method
func (m *MockConn) RemovePasswordHistory(authID string, historySize, historyDays int) error {
	ret := m.ctrl.Call(m, "RemovePasswordHistory", authID, historySize, historyDays)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemovePasswordHistory indicates an expected call of RemovePasswordHistory
func (mr *MockConnMockRecorder) RemovePasswordHistory(authID, historySize, historyDays interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemovePasswordHistory", reflect.TypeOf((*MockConn)(nil).RemovePasswordHistory), authID, historySize, historyDays)
}

// GetAdminRoles mocks base method
func (m *MockConn) GetAdminRoles() ([]string, error) {
	ret := m.ctrl.Call(m, "GetAdminRoles")
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAdminRoles indicates an expected call of GetAdminRoles
func (mr *MockConnMockRecorder) GetAdminRoles() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAdminRoles", reflect.TypeOf((*MockConn)(nil).GetAdminRoles))
}

// SetAdminRoles mocks base method
func (m *MockConn) SetAdminRoles(roles []string) error {
	ret := m.ctrl.Call(m, "SetAdminRoles", roles)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetAdminRoles indicates an expected call of SetAdminRoles
func (mr *MockConnMockRecorder) SetAdminRoles(roles interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetAdminRoles", reflect.TypeOf((*MockConn)(nil).SetAdminRoles), roles)
}

// GetDefaultRoles mocks base method
func (m *MockConn) GetDefaultRoles() ([]string, error) {
	ret := m.ctrl.Call(m, "GetDefaultRoles")
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDefaultRoles indicates an expected call of GetDefaultRoles
func (mr *MockConnMockRecorder) GetDefaultRoles() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDefaultRoles", reflect.TypeOf((*MockConn)(nil).GetDefaultRoles))
}

// SetDefaultRoles mocks base method
func (m *MockConn) SetDefaultRoles(roles []string) error {
	ret := m.ctrl.Call(m, "SetDefaultRoles", roles)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetDefaultRoles indicates an expected call of SetDefaultRoles
func (mr *MockConnMockRecorder) SetDefaultRoles(roles interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetDefaultRoles", reflect.TypeOf((*MockConn)(nil).SetDefaultRoles), roles)
}

// AssignRoles mocks base method
func (m *MockConn) AssignRoles(userIDs, roles []string) error {
	ret := m.ctrl.Call(m, "AssignRoles", userIDs, roles)
	ret0, _ := ret[0].(error)
	return ret0
}

// AssignRoles indicates an expected call of AssignRoles
func (mr *MockConnMockRecorder) AssignRoles(userIDs, roles interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AssignRoles", reflect.TypeOf((*MockConn)(nil).AssignRoles), userIDs, roles)
}

// RevokeRoles mocks base method
func (m *MockConn) RevokeRoles(userIDs, roles []string) error {
	ret := m.ctrl.Call(m, "RevokeRoles", userIDs, roles)
	ret0, _ := ret[0].(error)
	return ret0
}

// RevokeRoles indicates an expected call of RevokeRoles
func (mr *MockConnMockRecorder) RevokeRoles(userIDs, roles interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RevokeRoles", reflect.TypeOf((*MockConn)(nil).RevokeRoles), userIDs, roles)
}

// GetRoles mocks base method
func (m *MockConn) GetRoles(userIDs []string) (map[string][]string, error) {
	ret := m.ctrl.Call(m, "GetRoles", userIDs)
	ret0, _ := ret[0].(map[string][]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRoles indicates an expected call of GetRoles
func (mr *MockConnMockRecorder) GetRoles(userIDs interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRoles", reflect.TypeOf((*MockConn)(nil).GetRoles), userIDs)
}

// SetRecordAccess mocks base method
func (m *MockConn) SetRecordAccess(recordType string, acl RecordACL) error {
	ret := m.ctrl.Call(m, "SetRecordAccess", recordType, acl)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetRecordAccess indicates an expected call of SetRecordAccess
func (mr *MockConnMockRecorder) SetRecordAccess(recordType, acl interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetRecordAccess", reflect.TypeOf((*MockConn)(nil).SetRecordAccess), recordType, acl)
}

// SetRecordDefaultAccess mocks base method
func (m *MockConn) SetRecordDefaultAccess(recordType string, acl RecordACL) error {
	ret := m.ctrl.Call(m, "SetRecordDefaultAccess", recordType, acl)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetRecordDefaultAccess indicates an expected call of SetRecordDefaultAccess
func (mr *MockConnMockRecorder) SetRecordDefaultAccess(recordType, acl interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetRecordDefaultAccess", reflect.TypeOf((*MockConn)(nil).SetRecordDefaultAccess), recordType, acl)
}

// GetRecordAccess mocks base method
func (m *MockConn) GetRecordAccess(recordType string) (RecordACL, error) {
	ret := m.ctrl.Call(m, "GetRecordAccess", recordType)
	ret0, _ := ret[0].(RecordACL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRecordAccess indicates an expected call of GetRecordAccess
func (mr *MockConnMockRecorder) GetRecordAccess(recordType interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRecordAccess", reflect.TypeOf((*MockConn)(nil).GetRecordAccess), recordType)
}

// GetRecordDefaultAccess mocks base method
func (m *MockConn) GetRecordDefaultAccess(recordType string) (RecordACL, error) {
	ret := m.ctrl.Call(m, "GetRecordDefaultAccess", recordType)
	ret0, _ := ret[0].(RecordACL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRecordDefaultAccess indicates an expected call of GetRecordDefaultAccess
func (mr *MockConnMockRecorder) GetRecordDefaultAccess(recordType interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRecordDefaultAccess", reflect.TypeOf((*MockConn)(nil).GetRecordDefaultAccess), recordType)
}

// SetRecordFieldAccess mocks base method
func (m *MockConn) SetRecordFieldAccess(acl FieldACL) error {
	ret := m.ctrl.Call(m, "SetRecordFieldAccess", acl)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetRecordFieldAccess indicates an expected call of SetRecordFieldAccess
func (mr *MockConnMockRecorder) SetRecordFieldAccess(acl interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetRecordFieldAccess", reflect.TypeOf((*MockConn)(nil).SetRecordFieldAccess), acl)
}

// GetRecordFieldAccess mocks base method
func (m *MockConn) GetRecordFieldAccess() (FieldACL, error) {
	ret := m.ctrl.Call(m, "GetRecordFieldAccess")
	ret0, _ := ret[0].(FieldACL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRecordFieldAccess indicates an expected call of GetRecordFieldAccess
func (mr *MockConnMockRecorder) GetRecordFieldAccess() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRecordFieldAccess", reflect.TypeOf((*MockConn)(nil).GetRecordFieldAccess))
}

// GetAsset mocks base method
func (m *MockConn) GetAsset(name string, asset *Asset) error {
	ret := m.ctrl.Call(m, "GetAsset", name, asset)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetAsset indicates an expected call of GetAsset
func (mr *MockConnMockRecorder) GetAsset(name, asset interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAsset", reflect.TypeOf((*MockConn)(nil).GetAsset), name, asset)
}

// GetAssets mocks base method
func (m *MockConn) GetAssets(names []string) ([]Asset, error) {
	ret := m.ctrl.Call(m, "GetAssets", names)
	ret0, _ := ret[0].([]Asset)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAssets indicates an expected call of GetAssets
func (mr *MockConnMockRecorder) GetAssets(names interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAssets", reflect.TypeOf((*MockConn)(nil).GetAssets), names)
}

// SaveAsset mocks base method
func (m *MockConn) SaveAsset(asset *Asset) error {
	ret := m.ctrl.Call(m, "SaveAsset", asset)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveAsset indicates an expected call of SaveAsset
func (mr *MockConnMockRecorder) SaveAsset(asset interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveAsset", reflect.TypeOf((*MockConn)(nil).SaveAsset), asset)
}

// QueryRelation mocks base method
func (m *MockConn) QueryRelation(user, name, direction string, config QueryConfig) []AuthInfo {
	ret := m.ctrl.Call(m, "QueryRelation", user, name, direction, config)
	ret0, _ := ret[0].([]AuthInfo)
	return ret0
}

// QueryRelation indicates an expected call of QueryRelation
func (mr *MockConnMockRecorder) QueryRelation(user, name, direction, config interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryRelation", reflect.TypeOf((*MockConn)(nil).QueryRelation), user, name, direction, config)
}

// QueryRelationCount mocks base method
func (m *MockConn) QueryRelationCount(user, name, direction string) (uint64, error) {
	ret := m.ctrl.Call(m, "QueryRelationCount", user, name, direction)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryRelationCount indicates an expected call of QueryRelationCount
func (mr *MockConnMockRecorder) QueryRelationCount(user, name, direction interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryRelationCount", reflect.TypeOf((*MockConn)(nil).QueryRelationCount), user, name, direction)
}

// AddRelation mocks base method
func (m *MockConn) AddRelation(user, name, targetUser string) error {
	ret := m.ctrl.Call(m, "AddRelation", user, name, targetUser)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddRelation indicates an expected call of AddRelation
func (mr *MockConnMockRecorder) AddRelation(user, name, targetUser interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddRelation", reflect.TypeOf((*MockConn)(nil).AddRelation), user, name, targetUser)
}

// RemoveRelation mocks base method
func (m *MockConn) RemoveRelation(user, name, targetUser string) error {
	ret := m.ctrl.Call(m, "RemoveRelation", user, name, targetUser)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveRelation indicates an expected call of RemoveRelation
func (mr *MockConnMockRecorder) RemoveRelation(user, name, targetUser interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveRelation", reflect.TypeOf((*MockConn)(nil).RemoveRelation), user, name, targetUser)
}

// GetDevice mocks base method
func (m *MockConn) GetDevice(id string, device *Device) error {
	ret := m.ctrl.Call(m, "GetDevice", id, device)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetDevice indicates an expected call of GetDevice
func (mr *MockConnMockRecorder) GetDevice(id, device interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDevice", reflect.TypeOf((*MockConn)(nil).GetDevice), id, device)
}

// QueryDevicesByUser mocks base method
func (m *MockConn) QueryDevicesByUser(user string) ([]Device, error) {
	ret := m.ctrl.Call(m, "QueryDevicesByUser", user)
	ret0, _ := ret[0].([]Device)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryDevicesByUser indicates an expected call of QueryDevicesByUser
func (mr *MockConnMockRecorder) QueryDevicesByUser(user interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryDevicesByUser", reflect.TypeOf((*MockConn)(nil).QueryDevicesByUser), user)
}

// QueryDevicesByUserAndTopic mocks base method
func (m *MockConn) QueryDevicesByUserAndTopic(user, topic string) ([]Device, error) {
	ret := m.ctrl.Call(m, "QueryDevicesByUserAndTopic", user, topic)
	ret0, _ := ret[0].([]Device)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryDevicesByUserAndTopic indicates an expected call of QueryDevicesByUserAndTopic
func (mr *MockConnMockRecorder) QueryDevicesByUserAndTopic(user, topic interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryDevicesByUserAndTopic", reflect.TypeOf((*MockConn)(nil).QueryDevicesByUserAndTopic), user, topic)
}

// SaveDevice mocks base method
func (m *MockConn) SaveDevice(device *Device) error {
	ret := m.ctrl.Call(m, "SaveDevice", device)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveDevice indicates an expected call of SaveDevice
func (mr *MockConnMockRecorder) SaveDevice(device interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveDevice", reflect.TypeOf((*MockConn)(nil).SaveDevice), device)
}

// DeleteDevice mocks base method
func (m *MockConn) DeleteDevice(id string) error {
	ret := m.ctrl.Call(m, "DeleteDevice", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteDevice indicates an expected call of DeleteDevice
func (mr *MockConnMockRecorder) DeleteDevice(id interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteDevice", reflect.TypeOf((*MockConn)(nil).DeleteDevice), id)
}

// DeleteDevicesByToken mocks base method
func (m *MockConn) DeleteDevicesByToken(token string, t time.Time) error {
	ret := m.ctrl.Call(m, "DeleteDevicesByToken", token, t)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteDevicesByToken indicates an expected call of DeleteDevicesByToken
func (mr *MockConnMockRecorder) DeleteDevicesByToken(token, t interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteDevicesByToken", reflect.TypeOf((*MockConn)(nil).DeleteDevicesByToken), token, t)
}

// DeleteEmptyDevicesByTime mocks base method
func (m *MockConn) DeleteEmptyDevicesByTime(t time.Time) error {
	ret := m.ctrl.Call(m, "DeleteEmptyDevicesByTime", t)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteEmptyDevicesByTime indicates an expected call of DeleteEmptyDevicesByTime
func (mr *MockConnMockRecorder) DeleteEmptyDevicesByTime(t interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteEmptyDevicesByTime", reflect.TypeOf((*MockConn)(nil).DeleteEmptyDevicesByTime), t)
}

// PublicDB mocks base method
func (m *MockConn) PublicDB() Database {
	ret := m.ctrl.Call(m, "PublicDB")
	ret0, _ := ret[0].(Database)
	return ret0
}

// PublicDB indicates an expected call of PublicDB
func (mr *MockConnMockRecorder) PublicDB() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublicDB", reflect.TypeOf((*MockConn)(nil).PublicDB))
}

// PrivateDB mocks base method
func (m *MockConn) PrivateDB(userKey string) Database {
	ret := m.ctrl.Call(m, "PrivateDB", userKey)
	ret0, _ := ret[0].(Database)
	return ret0
}

// PrivateDB indicates an expected call of PrivateDB
func (mr *MockConnMockRecorder) PrivateDB(userKey interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrivateDB", reflect.TypeOf((*MockConn)(nil).PrivateDB), userKey)
}

// UnionDB mocks base method
func (m *MockConn) UnionDB() Database {
	ret := m.ctrl.Call(m, "UnionDB")
	ret0, _ := ret[0].(Database)
	return ret0
}

// UnionDB indicates an expected call of UnionDB
func (mr *MockConnMockRecorder) UnionDB() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnionDB", reflect.TypeOf((*MockConn)(nil).UnionDB))
}

// Subscribe mocks base method
func (m *MockConn) Subscribe(recordEventChan chan RecordEvent) error {
	ret := m.ctrl.Call(m, "Subscribe", recordEventChan)
	ret0, _ := ret[0].(error)
	return ret0
}

// Subscribe indicates an expected call of Subscribe
func (mr *MockConnMockRecorder) Subscribe(recordEventChan interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockConn)(nil).Subscribe), recordEventChan)
}

// EnsureAuthRecordKeysExist mocks base method
func (m *MockConn) EnsureAuthRecordKeysExist(authRecordKeys [][]string) error {
	ret := m.ctrl.Call(m, "EnsureAuthRecordKeysExist", authRecordKeys)
	ret0, _ := ret[0].(error)
	return ret0
}

// EnsureAuthRecordKeysExist indicates an expected call of EnsureAuthRecordKeysExist
func (mr *MockConnMockRecorder) EnsureAuthRecordKeysExist(authRecordKeys interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnsureAuthRecordKeysExist", reflect.TypeOf((*MockConn)(nil).EnsureAuthRecordKeysExist), authRecordKeys)
}

// EnsureAuthRecordKeysIndexesMatch mocks base method
func (m *MockConn) EnsureAuthRecordKeysIndexesMatch(authRecordKeys [][]string) error {
	ret := m.ctrl.Call(m, "EnsureAuthRecordKeysIndexesMatch", authRecordKeys)
	ret0, _ := ret[0].(error)
	return ret0
}

// EnsureAuthRecordKeysIndexesMatch indicates an expected call of EnsureAuthRecordKeysIndexesMatch
func (mr *MockConnMockRecorder) EnsureAuthRecordKeysIndexesMatch(authRecordKeys interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnsureAuthRecordKeysIndexesMatch", reflect.TypeOf((*MockConn)(nil).EnsureAuthRecordKeysIndexesMatch), authRecordKeys)
}

// CreateOAuthInfo mocks base method
func (m *MockConn) CreateOAuthInfo(oauthinfo *OAuthInfo) error {
	ret := m.ctrl.Call(m, "CreateOAuthInfo", oauthinfo)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateOAuthInfo indicates an expected call of CreateOAuthInfo
func (mr *MockConnMockRecorder) CreateOAuthInfo(oauthinfo interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOAuthInfo", reflect.TypeOf((*MockConn)(nil).CreateOAuthInfo), oauthinfo)
}

// GetOAuthInfo mocks base method
func (m *MockConn) GetOAuthInfo(provider, principalID string, oauthinfo *OAuthInfo) error {
	ret := m.ctrl.Call(m, "GetOAuthInfo", provider, principalID, oauthinfo)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetOAuthInfo indicates an expected call of GetOAuthInfo
func (mr *MockConnMockRecorder) GetOAuthInfo(provider, principalID, oauthinfo interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOAuthInfo", reflect.TypeOf((*MockConn)(nil).GetOAuthInfo), provider, principalID, oauthinfo)
}

// GetOAuthInfoByProviderAndUserID mocks base method
func (m *MockConn) GetOAuthInfoByProviderAndUserID(provider, userID string, oauthinfo *OAuthInfo) error {
	ret := m.ctrl.Call(m, "GetOAuthInfoByProviderAndUserID", provider, userID, oauthinfo)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetOAuthInfoByProviderAndUserID indicates an expected call of GetOAuthInfoByProviderAndUserID
func (mr *MockConnMockRecorder) GetOAuthInfoByProviderAndUserID(provider, userID, oauthinfo interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOAuthInfoByProviderAndUserID", reflect.TypeOf((*MockConn)(nil).GetOAuthInfoByProviderAndUserID), provider, userID, oauthinfo)
}

// UpdateOAuthInfo mocks base method
func (m *MockConn) UpdateOAuthInfo(oauthinfo *OAuthInfo) error {
	ret := m.ctrl.Call(m, "UpdateOAuthInfo", oauthinfo)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateOAuthInfo indicates an expected call of UpdateOAuthInfo
func (mr *MockConnMockRecorder) UpdateOAuthInfo(oauthinfo interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOAuthInfo", reflect.TypeOf((*MockConn)(nil).UpdateOAuthInfo), oauthinfo)
}

// DeleteOAuth mocks base method
func (m *MockConn) DeleteOAuth(provider, principalID string) error {
	ret := m.ctrl.Call(m, "DeleteOAuth", provider, principalID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteOAuth indicates an expected call of DeleteOAuth
func (mr *MockConnMockRecorder) DeleteOAuth(provider, principalID interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteOAuth", reflect.TypeOf((*MockConn)(nil).DeleteOAuth), provider, principalID)
}

// Close mocks base method
func (m *MockConn) Close() error {
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockConnMockRecorder) Close() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockConn)(nil).Close))
}

// GetCustomTokenInfo mocks base method
func (m *MockConn) GetCustomTokenInfo(principalID string, tokenInfo *CustomTokenInfo) error {
	ret := m.ctrl.Call(m, "GetCustomTokenInfo", principalID, tokenInfo)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetCustomTokenInfo indicates an expected call of GetCustomTokenInfo
func (mr *MockConnMockRecorder) GetCustomTokenInfo(principalID, tokenInfo interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCustomTokenInfo", reflect.TypeOf((*MockConn)(nil).GetCustomTokenInfo), principalID, tokenInfo)
}

// CreateCustomTokenInfo mocks base method
func (m *MockConn) CreateCustomTokenInfo(tokenInfo *CustomTokenInfo) error {
	ret := m.ctrl.Call(m, "CreateCustomTokenInfo", tokenInfo)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateCustomTokenInfo indicates an expected call of CreateCustomTokenInfo
func (mr *MockConnMockRecorder) CreateCustomTokenInfo(tokenInfo interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCustomTokenInfo", reflect.TypeOf((*MockConn)(nil).CreateCustomTokenInfo), tokenInfo)
}

// DeleteCustomTokenInfo mocks base method
func (m *MockConn) DeleteCustomTokenInfo(principalID string) error {
	ret := m.ctrl.Call(m, "DeleteCustomTokenInfo", principalID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCustomTokenInfo indicates an expected call of DeleteCustomTokenInfo
func (mr *MockConnMockRecorder) DeleteCustomTokenInfo(principalID interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCustomTokenInfo", reflect.TypeOf((*MockConn)(nil).DeleteCustomTokenInfo), principalID)
}

// MockCustomTokenConn is a mock of CustomTokenConn interface
type MockCustomTokenConn struct {
	ctrl     *gomock.Controller
	recorder *MockCustomTokenConnMockRecorder
}

// MockCustomTokenConnMockRecorder is the mock recorder for MockCustomTokenConn
type MockCustomTokenConnMockRecorder struct {
	mock *MockCustomTokenConn
}

// NewMockCustomTokenConn creates a new mock instance
func NewMockCustomTokenConn(ctrl *gomock.Controller) *MockCustomTokenConn {
	mock := &MockCustomTokenConn{ctrl: ctrl}
	mock.recorder = &MockCustomTokenConnMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCustomTokenConn) EXPECT() *MockCustomTokenConnMockRecorder {
	return m.recorder
}

// GetCustomTokenInfo mocks base method
func (m *MockCustomTokenConn) GetCustomTokenInfo(principalID string, tokenInfo *CustomTokenInfo) error {
	ret := m.ctrl.Call(m, "GetCustomTokenInfo", principalID, tokenInfo)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetCustomTokenInfo indicates an expected call of GetCustomTokenInfo
func (mr *MockCustomTokenConnMockRecorder) GetCustomTokenInfo(principalID, tokenInfo interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCustomTokenInfo", reflect.TypeOf((*MockCustomTokenConn)(nil).GetCustomTokenInfo), principalID, tokenInfo)
}

// CreateCustomTokenInfo mocks base method
func (m *MockCustomTokenConn) CreateCustomTokenInfo(tokenInfo *CustomTokenInfo) error {
	ret := m.ctrl.Call(m, "CreateCustomTokenInfo", tokenInfo)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateCustomTokenInfo indicates an expected call of CreateCustomTokenInfo
func (mr *MockCustomTokenConnMockRecorder) CreateCustomTokenInfo(tokenInfo interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCustomTokenInfo", reflect.TypeOf((*MockCustomTokenConn)(nil).CreateCustomTokenInfo), tokenInfo)
}

// DeleteCustomTokenInfo mocks base method
func (m *MockCustomTokenConn) DeleteCustomTokenInfo(principalID string) error {
	ret := m.ctrl.Call(m, "DeleteCustomTokenInfo", principalID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCustomTokenInfo indicates an expected call of DeleteCustomTokenInfo
func (mr *MockCustomTokenConnMockRecorder) DeleteCustomTokenInfo(principalID interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCustomTokenInfo", reflect.TypeOf((*MockCustomTokenConn)(nil).DeleteCustomTokenInfo), principalID)
}
