#!/bin/bash
jq .features[].properties.lad11cd "Local_Authority_Districts_(December_2011)_Boundaries_EW_BGC.geojson" | grep -v ^\"S|grep -v ^\"N|grep -v ^\"L"|grep -v ^\"M" | sed 's/\"//g'  |sort
