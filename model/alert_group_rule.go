// pmm-ruled
// Copyright (C) 2019 gywndi@gmail.com in kakaoBank
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package model

import (
	"fmt"
	"pmm-ruled/common"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// AlertGroupRule Alerts for group
type AlertGroupRule struct {
	GroupID   int       `form:"group_id" json:"group_id" xorm:"group_id pk not null "`
	RuleID    int       `form:"rule_id" json:"rule_id" xorm:"rule_id pk not null index(01)"`
	Val       *string   `form:"val" json:"val" xorm:"val varchar(30) not null default '' "`
	CreatedAt time.Time `json:"created_at" xorm:"datetime not null created"`
	UpdatedAt time.Time `json:"updated_at" xorm:"datetime not null updated"`
}

// Exist check exists
func (o *AlertGroupRule) Exist() bool {
	boolean, _ := orm.Exist(o)
	return boolean
}

// GetFirst get first one
func (o *AlertGroupRule) GetFirst() (AlertGroupRule, error) {
	var err error

	ret := *o
	boolean, err := orm.Get(&ret)
	if err != nil {
		return ret, err
	}

	if !boolean {
		return ret, fmt.Errorf(common.MSG["err.row_not_found"])
	}

	return ret, err
}

// GetList get rows
func (o *AlertGroupRule) GetList(sort ...string) ([]AlertGroupRule, error) {
	var err error
	var arr []AlertGroupRule
	var order string

	for i, s := range sort {
		if i > 0 {
			order += ","
		}
		order += s
	}
	err = orm.OrderBy(order).Find(&arr, o)
	common.Log.Info(reflect.TypeOf(o), len(arr), " rows selected")
	return arr, err
}

// Insert new row
func (o *AlertGroupRule) Insert() error {
	var err error
	var affected int64

	session := orm.NewSession()
	defer session.Close()

	o.rewriteCols()

	if err = o.InsertCheck(); err != nil {
		return err
	}

	if affected, err = session.Insert(o); err != nil {
		return err
	}
	common.Log.Info(reflect.TypeOf(o), affected, "rows inserted!")

	return err
}

// Update update row (partitial column)
func (o *AlertGroupRule) Update(to *AlertGroupRule) (int64, error) {
	var err error
	var affected int64

	session := orm.NewSession()
	defer session.Close()

	o.rewriteCols()

	if err = to.UpdateCheck(); err != nil {
		return affected, err
	}

	if affected, err = session.Update(to, o); err != nil {
		return affected, err
	}

	common.Log.Info(reflect.TypeOf(o), affected, "rows updated!")
	return affected, err
}

// Delete delete row
func (o *AlertGroupRule) Delete() (int64, error) {
	var err error
	var affected int64

	session := orm.NewSession()
	defer session.Close()

	if err = o.DeleteCheck(); err != nil {
		return affected, err
	}

	if affected, err = session.Delete(o); err != nil {
		return affected, err
	}

	common.Log.Info(reflect.TypeOf(o), affected, "rows deleted!")
	return affected, err
}

// InsertCheck validation check
func (o *AlertGroupRule) InsertCheck() error {
	var err error

	// Check group  exists
	if !(&AlertGroup{ID: o.GroupID}).Exist() {
		return fmt.Errorf(common.MSG["err.group_not_exists"])
	}

	// Check Alert exists
	if !(&AlertRule{ID: o.RuleID}).Exist() {
		return fmt.Errorf(common.MSG["err.rule_not_exists"])
	}

	if o.Val == nil {
		o.Val = new(string)
	}

	// Val digit check
	if o.Val != nil && *o.Val != "" {
		if _, err := strconv.ParseFloat(*o.Val, 64); err != nil {
			return fmt.Errorf(common.MSG["err.val_not_digit"])
		}
	}

	return err
}

// UpdateCheck validation check
func (o *AlertGroupRule) UpdateCheck() error {
	var err error

	// check exists and get first row
	if _, err := o.GetFirst(); err != nil {
		return fmt.Errorf(common.MSG["err.rule_not_exists"])
	}

	if o.Val == nil {
		return fmt.Errorf(common.MSG["err.val_empty"])
	}

	// Val digit check
	if o.Val != nil && *o.Val != "" {
		if _, err := strconv.ParseFloat(*o.Val, 64); err != nil {
			return fmt.Errorf(common.MSG["err.val_not_digit"])
		}
	}

	return err
}

// DeleteCheck validation check
func (o *AlertGroupRule) DeleteCheck() error {
	var err error

	return err
}

// rewriteCols rewrite column value
func (o *AlertGroupRule) rewriteCols() {
	if o.Val != nil {
		*o.Val = strings.TrimSpace(*o.Val)
	}
}
