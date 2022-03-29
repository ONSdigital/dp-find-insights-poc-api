# Generate static file geoLookup.json for the frontend

https://github.com/ONSdigital/dp-census-atlas/blob/develop/src/data/geoLookup.json

```
go run ./cmd/geobb   
```

A large (~800K) static file which contains LAD and MSOA.

Sample record

```
[
 {
    "en": "Merthyr Tydfil",
    "cy": "",
    "geoType": "LAD",
    "geoCode": "W06000024",
    "bbox": [
      -3.45437,
      51.64512,
      -3.27517,
      51.83552
    ]
  }
]
```

Currently there is no Welsh content.

`geobb` isn't a great package name.  It's possible this becomes a REST endpoint.
