package postgres

import "errors"

var SubscriptionAlreadyExist = errors.New("subscription already exists")
var SubscriptionNotFound = errors.New("subscription not found")
var SubscriptionDateError = errors.New("subscription date error")
