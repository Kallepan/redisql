/**
* This file exists only to list the possible types of a Result
* Ignore this file :)
**/
package app

import "time"

type Result interface {
	int64 | int32 | int16 | int8 | string | float64 | float32 | bool | time.Duration | time.Time
}
