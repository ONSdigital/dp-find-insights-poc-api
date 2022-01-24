// Package Swagger provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/pseudo-su/oapi-ui-codegen DO NOT EDIT.
package Swagger

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"strings"
)

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+xYbW/jNhL+KwPeAdfibOslyW7ib4u02OZ2uym2AQ5osQgoaSSxS5EqX+wYC//3AynJ",
	"cSzacbPpIUCbT7FIzjwznHlmhl9ILptWChRGk/kXovMaG+r/vaQGK6kY+l/MYOP/+afCkszJP6L7g1F/",
	"KrpRrOVoyHpCzKpFMidUKbpyv79XSip3vlWyRWV6sTh8LlDnirWGSUHm3WdoUGtaIZkQvKNNy53AXFpe",
	"gJAGNF1BjZxLstGmjWKiIuv1hCj83TKFBZn/2iv5tNkms98w9yh/QMpNPYaV15h/Pt7uTsylO4QqZL02",
	"VJlbwxoc23pTI/h1KKhBoKIAtxFkCW9+ugJlhXBGbTshjdN4Gr+aJslNksxPL+ZpMjtL44s0/WXsDK/d",
	"WL1Xs7HaKTM1OoVOkbCN89v1OzIh/33z8cPVh7dkQi4/Xt1cXb55v+XJex223W9dt9Yb9MCQk9Oz5FUI",
	"8gKV9gJ2byazjBcHPOnXx57cMi7oxdR7Mf53nMzjOASoYuY2l03DTFhvxQx061BTXe/T+TpPS8yyMs3O",
	"k/Pk9VmSpKevz4vTssxokSEm2auz0/LVSQgCp6KyLh+CAFolK0WbhokKhp1gNRZgJDCnvkFhRoAqeUjV",
	"7dY9jFX2i4OtT0aQzJLT2ckjYXBQ/a7MZBbP4iAv7FDAei8pDNk8ikBOtbn1BOHYJQTM7YDaSwG/MRyP",
	"eGdQCcpBo1qwHA+n+Fk8OzmJ4/OLX8IXps1tSRm3Cg+Acjuw+HpsycU0vpimqcd2Pj9LZrH/S/aD0zbP",
	"UesD4PodpeX/Z+cNdSYIrV/cIcqD6hspKllkwDR4Ch0pFHQffbkVp2NXfpdH2WqHoXtNQUY+nvVDxvzh",
	"EhDKpB/R0IIaGiiwsvAe2OuasTncVsEFQzOOj7cm3a6DMD+ibqXQeHTd39gXKPleYcDwB13VIeFb/dd6",
	"8qwOk4byo1u5kMNuNi4/ri30rgj5qFfy54VHgPHXE8JEKceJ8f1di4q5GkU50LblLKduCUqpoMAFctm6",
	"2nZ5FV1+By1rkTPhkoWzHPuw6fCR6xYFvJULVMLXvPduR46wOPFlySpO5qQ2pp1H0XK5nAmviHKq8pot",
	"UM8quZjZz1Eh80i2KKbVRtaUd7KivvxFJ5F3LTOee4p2WjJRTJnQrKqNnrYyn9KWka1i2pfH9YQ42W5x",
	"Tk76itlSU/sriPLPDVKhoy8rpGrtPlW42/xoMicf0VglNFDgTJuBUxp6xxrbwIJyixqYAKR5DZ1QyLnV",
	"BtXMK1S0QYNKk/mvu1fi/Go1OATEXRuZe4BkCATSr9x3+0ZZnPRjzFY8MGGwcr35evKlk/O7RbW6F5RT",
	"QwLn7uMofKxC6Xc+4ehnchDnJ2dVx0f+PtI47jJDGBT+HrZCNPpNd13SvcBNZpZSNdS40JDW5eEmI4Rt",
	"stC44jJk5xq6UNiq0HwFOeW55dRg4USk8WkgNoQEhdpyo6GUVvidp3/QkEPM0k2VAcQN09rlqlSQ0YKv",
	"XA43WAATrTV9UBJ/qqSWmz8fkLvYTY1FBdhvnBBtm4aqlY/B3qMwONxlPVCo2AIF9KVj5XuhCl2z3dYr",
	"6OPP0MplEGltxllOPjnRQw4rZ8pfJ5GTp6RjTk36NwP8zQAvjAFc7+HS0KcwZGiWiALMUj4kBYYavnGh",
	"D5H7lH7rDw7McRxZ1JvHsDE/bNFDP5uNHo48LXUPZyAFFNiiKFzz008VmjwhnoXlfHKk5/vHvLHnf94O",
	"WtWbcf1uMGHJhlHTmxMA7mI2vXi2ENkAHSPtNcKSKv+OYlvwzWelaIEFfEMNcHTDshTYj8dMQD+gua3D",
	"hNYb9+2LC/EhjN78dPWvnWDaCsxCDlHZbM2QB+NS+ycmuJEty/3EAZdDwfzOnR9VqBAll4wbVH5C0iF2",
	"zqTkSMUzsPMxc+VmIA34+vrdi7vat2hgMxLvIRnv73Er8kKaByWXD699U0PHU/TDirmvq5D8eQVmmbx7",
	"rjbla7Fw2cXYU/AoWjCryVMuqZV8VT1NaxdJ/sXo4PHHk9vgnYlyvehK1PqlpaK3Gzpz96WiXtKq6l64",
	"g8SqW8yfWrNHoO+Hi//8fP3BN2Nd78A0lIzjDv5e94Db/3yA2rKvxO0vsDYNfwTwRh/8cPPjew/8OKzr",
	"9f8CAAD//4NCTahdHQAA",
}

// GetOpenAPISpec returns the Swagger specification corresponding to the generated code
// in this file.
func GetOpenAPISpec() (*openapi3.T, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	swagger, err := openapi3.NewLoader().LoadFromData(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("error loading Swagger: %s", err)
	}
	return swagger, nil
}
