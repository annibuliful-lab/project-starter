package graphql_enum

import (
	"fmt"
)

type PermissionAbility int

const (
	CREATE PermissionAbility = iota
	UPDATE
	DELETE
	READ
	EXECUTE
)

var permissionAbilityStates = [...]string{"CREATE", "UPDATE", "DELETE", "READ", "EXECUTE"}

func GetPermissionAbility(str string) PermissionAbility {

	for i, st := range permissionAbilityStates {
		if st == str {
			return PermissionAbility(i)
		}
	}

	panic("invalid value for enum State: " + str)

}

func (s PermissionAbility) String() string { return permissionAbilityStates[s] }

func (s *PermissionAbility) Deserialize(str string) {
	var found bool
	for i, st := range permissionAbilityStates {
		if st == str {
			found = true
			(*s) = PermissionAbility(i)
		}
	}
	if !found {
		panic("invalid value for enum State: " + str)
	}
}

func (PermissionAbility) ImplementsGraphQLType(name string) bool {
	return name == "PermissionAbility"
}

func (s *PermissionAbility) UnmarshalGraphQL(input interface{}) error {
	var err error
	switch input := input.(type) {
	case string:
		s.Deserialize(input)
	default:
		err = fmt.Errorf("wrong type for State: %T", input)
	}
	return err
}
