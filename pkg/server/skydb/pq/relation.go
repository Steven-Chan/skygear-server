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

package pq

import (
	"fmt"

	sq "github.com/lann/squirrel"
	"github.com/skygeario/skygear-server/pkg/server/skydb"
	"github.com/skygeario/skygear-server/pkg/server/skydb/pq/builder"
)

func (c *conn) QueryRelation(user string, name string, direction string, config skydb.QueryConfig) []skydb.AuthInfo {
	log.Debugf("Query Relation: %v, %v", user, name)
	var selectBuilder sq.SelectBuilder

	if direction == "outward" {
		selectBuilder = psql.Select("u.id").
			From(c.tableName("_auth")+" AS u").
			Join(c.tableName(name)+" AS relation ON relation.right_id = u.id").
			Where("relation.left_id = ?", user)
	} else if direction == "inward" {
		selectBuilder = psql.Select("u.id").
			From(c.tableName("_auth")+" AS u").
			Join(c.tableName(name)+" AS relation ON relation.left_id = u.id").
			Where("relation.right_id = ?", user)
	} else {
		selectBuilder = psql.Select("u.id").
			From(c.tableName("_auth")+" AS u").
			Join(c.tableName(name)+" AS inward_relation ON inward_relation.left_id = u.id").
			Join(c.tableName(name)+" AS outward_relation ON outward_relation.right_id = u.id").
			Where("inward_relation.right_id = ?", user).
			Where("outward_relation.left_id = ?", user)
	}

	selectBuilder = selectBuilder.OrderBy("u.id").
		Offset(config.Offset)
	if config.Limit != 0 {
		selectBuilder = selectBuilder.Limit(config.Limit)
	}

	rows, err := c.QueryWith(selectBuilder)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	results := []skydb.AuthInfo{}
	for rows.Next() {
		var (
			id string
		)
		if err := rows.Scan(&id); err != nil {
			panic(err)
		}
		authInfo := skydb.AuthInfo{
			ID: id,
		}
		results = append(results, authInfo)
	}
	return results
}

func (c *conn) QueryRelationCount(user string, name string, direction string) (uint64, error) {
	log.Debugf("Query Relation Count: %v, %v, %v", user, name, direction)
	query := psql.Select("COUNT(*)").From(c.tableName(name) + "AS _primary")
	if direction == "outward" {
		query = query.Where("_primary.left_id = ?", user)
	} else if direction == "inward" {
		query = query.Where("_primary.right_id = ?", user)
	} else {
		query = query.
			Join(c.tableName(name)+" AS _secondary ON _secondary.left_id = _primary.right_id").
			Where("_primary.left_id = ?", user).
			Where("_secondary.right_id = ?", user)
	}
	var count uint64
	err := c.GetWith(&count, query)
	if err != nil {
		panic(err)
	}
	return count, err
}

func (c *conn) AddRelation(user string, name string, targetUser string) error {
	ralationPair := map[string]interface{}{
		"left_id":  user,
		"right_id": targetUser,
	}

	upsert := builder.UpsertQuery(c.tableName(name), ralationPair, nil)
	_, err := c.ExecWith(upsert)
	if err != nil {
		if isForeignKeyViolated(err) {
			return fmt.Errorf("userID not exist")
		}
	}

	return err
}

func (c *conn) RemoveRelation(user string, name string, targetUser string) error {
	builder := psql.Delete(c.tableName(name)).
		Where("left_id = ? AND right_id = ?", user, targetUser)
	result, err := c.ExecWith(builder)

	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%v relation not exist {%v} => {%v}",
			name, user, targetUser)
	} else if rowsAffected > 1 {
		panic(fmt.Errorf("want 1 rows updated, got %v", rowsAffected))
	}
	return nil
}
