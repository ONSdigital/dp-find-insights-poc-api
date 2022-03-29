#!/bin/bash
jq .features[].properties.lad17cd "Local_Authority_Districts_(December_2017)_Boundaries_in_the_UK_(WGS84).geojson" | grep -v ^\"S|grep -v ^\"N|grep -v ^\"L"|grep -v ^\"M" | sed 's/\"//g'  |sort
