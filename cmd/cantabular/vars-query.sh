#!/usr/bin/env bash
# Eleanor: Usual-Residents and People-Households

echo "searching for '$1'"

echo "Usual-Residents"
go run . -ds Usual-Residents -variables | grep -i $1
echo "Household-Ref-Persons"
go run . -ds Household-Ref-Persons -variables | grep -i $1
echo "People-Households"
go run . -ds People-Households -variables | grep -i $1
