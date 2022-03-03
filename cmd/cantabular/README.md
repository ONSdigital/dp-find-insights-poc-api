# NOTES

main.go is an adhoc query tool for the 2011 Cantabular data in 

https://ftb-api-ext.ons.sensiblecode.io/graphql"

To put creds in env and see help run

```
gpg -d ftb-api-ext.env.gpg > ftb-api-ext.env
source ftb-api-ext.env
go run . -help
go run . -variables
```

The default dataset is "Usual-Residents" and there are "helper" scripts to
search across multiple datasets.  Eleanor (slack) says "Usual-Residents and
"People-Households" are important, eg.

```
./vars-query.sh birth   
```
Will search variables across datasets and returns COBG.

To drill down and see catalogies under COBG variable for datasets use:

```
./class-query.sh COBG
```

## Queries "1" & "2"

# Query 1 is based on geocode(s)

```
go run . -query1 -code QS501EW -geo E02004332,E02004333,E02004356,E02004341 -geotype MSOA
```
Note currently the geotype must match all geocode types and response includes a
field with a comma in!

Would be possible to make geotype optional

# Query 2 is based on geotype

```
go run . -query2 -code QS416EW -geotype Region  
```

## Using GraphiQL

See also schema in Cantabular docs (only seen 9.0 what about 9.2?)

There is a web UI at https://ftb-api-ext.ons.sensiblecode.io/graphql

Ctrl-space will auto complete inside it.

Example queries are:


```
{
  dataset(name: "Usual-Residents") {
    # "label": "Male", code=1 like QS104EW0002
    table(variables: ["LA", "SEX"]
      , filters: [{variable: "SEX", codes: ["1"]}, {variable: "LA", codes: ["synE06000001"]} ]) {
      dimensions {
        categories {
          label
          code
        }
      }
      values
      error
    }
  }
}


# case insensitive sub-search match

{
  dataset(name: "Usual-Residents") {
    variables(names: ["LA"]) {
      categorySearch(text: "Poole") {
        edges {
          node {
            label
            code
          }
        }
      }
    }
  }
}

```

Note that queries on 4 or more variables don't seem to be possible & return
null responses at least in some cases.

First query above is roughly equal to:

```
SELECT gm.metric FROM geo_metric gm, nomis_category c, geo g
WHERE gm.geo_id = g.id and gm.category_id=c.id
AND c.long_nomis_code='QS104EW0002' and g.code='E06000001';

```
