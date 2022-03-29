package geobb

import (
	"testing"

	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/database"
	"github.com/cockroachdb/copyist"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	copyist.Register("postgres")
}

func TestAsJSON(t *testing.T) {

	defer copyist.Open(t).Close()

	dsn := database.GetDSN() + "?sslmode=disable"

	// use recorded DB response
	// "go test -v ../ladbb -run TestAsJSON -record" to create new
	gdb, err := gorm.Open(postgres.New(postgres.Config{DriverName: "copyist_postgres", DSN: dsn}), &gorm.Config{})
	if err != nil {
		t.Errorf(err.Error())
	}

	g := GeoBB{Gdb: gdb}

	got := g.AsJSON(Params{Pretty: true, Geos: []string{"LAD", "LAD"}})

	if got != exp() {
		t.Fail()
	}

}
func exp() string {

	return `[
  {
   "en": "Hartlepool",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000001",
   "bbox": [
    -1.38534,
    54.62239,
    -1.15941,
    54.72727
   ]
  },
  {
   "en": "Middlesbrough",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000002",
   "bbox": [
    -1.28702,
    54.50295,
    -1.13897,
    54.5908
   ]
  },
  {
   "en": "Redcar and Cleveland",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000003",
   "bbox": [
    -1.2025,
    54.48801,
    -0.79008,
    54.6477
   ]
  },
  {
   "en": "Stockton-on-Tees",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000004",
   "bbox": [
    -1.45249,
    54.46429,
    -1.16219,
    54.64535
   ]
  },
  {
   "en": "Darlington",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000005",
   "bbox": [
    -1.71082,
    54.45151,
    -1.40887,
    54.61948
   ]
  },
  {
   "en": "Halton",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000006",
   "bbox": [
    -2.83382,
    53.30629,
    -2.59661,
    53.40233
   ]
  },
  {
   "en": "Warrington",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000007",
   "bbox": [
    -2.6991,
    53.3227,
    -2.428,
    53.48116
   ]
  },
  {
   "en": "Blackburn with Darwen",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000008",
   "bbox": [
    -2.566,
    53.61684,
    -2.36407,
    53.78124
   ]
  },
  {
   "en": "Blackpool",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000009",
   "bbox": [
    -3.06237,
    53.7733,
    -2.98598,
    53.87591
   ]
  },
  {
   "en": "Kingston upon Hull, City of",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000010",
   "bbox": [
    -0.42401,
    53.71981,
    -0.24312,
    53.81349
   ]
  },
  {
   "en": "East Riding of Yorkshire",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000011",
   "bbox": [
    -1.10515,
    53.57172,
    0.14534,
    54.1767
   ]
  },
  {
   "en": "North East Lincolnshire",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000012",
   "bbox": [
    -0.29379,
    53.43383,
    0.01557,
    53.63814
   ]
  },
  {
   "en": "North Lincolnshire",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000013",
   "bbox": [
    -0.9516,
    53.45559,
    -0.20611,
    53.71483
   ]
  },
  {
   "en": "York",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000014",
   "bbox": [
    -1.22368,
    53.87479,
    -0.92129,
    54.05679
   ]
  },
  {
   "en": "Derby",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000015",
   "bbox": [
    -1.55835,
    52.86162,
    -1.38546,
    52.96842
   ]
  },
  {
   "en": "Leicester",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000016",
   "bbox": [
    -1.2172,
    52.58138,
    -1.04852,
    52.69186
   ]
  },
  {
   "en": "Rutland",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000017",
   "bbox": [
    -0.82334,
    52.52533,
    -0.43208,
    52.76005
   ]
  },
  {
   "en": "Nottingham",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000018",
   "bbox": [
    -1.24838,
    52.88934,
    -1.08767,
    53.01887
   ]
  },
  {
   "en": "Herefordshire, County of",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000019",
   "bbox": [
    -3.14251,
    51.82655,
    -2.33968,
    52.39565
   ]
  },
  {
   "en": "Telford and Wrekin",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000020",
   "bbox": [
    -2.66872,
    52.61494,
    -2.31361,
    52.82846
   ]
  },
  {
   "en": "Stoke-on-Trent",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000021",
   "bbox": [
    -2.24018,
    52.9465,
    -2.0809,
    53.093
   ]
  },
  {
   "en": "Bath and North East Somerset",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000022",
   "bbox": [
    -2.70919,
    51.2736,
    -2.2799,
    51.44002
   ]
  },
  {
   "en": "Bristol, City of",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000023",
   "bbox": [
    -2.71966,
    51.39803,
    -2.51233,
    51.5449
   ]
  },
  {
   "en": "North Somerset",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000024",
   "bbox": [
    -3.11631,
    51.29111,
    -2.58859,
    51.50315
   ]
  },
  {
   "en": "South Gloucestershire",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000025",
   "bbox": [
    -2.67512,
    51.41642,
    -2.25376,
    51.6777
   ]
  },
  {
   "en": "Plymouth",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000026",
   "bbox": [
    -4.20643,
    50.33368,
    -4.02073,
    50.44474
   ]
  },
  {
   "en": "Torbay",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000027",
   "bbox": [
    -3.6292,
    50.37409,
    -3.482,
    50.51827
   ]
  },
  {
   "en": "Bournemouth",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000028",
   "bbox": [
    -1.93766,
    50.71001,
    -1.74252,
    50.78067
   ]
  },
  {
   "en": "Poole",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000029",
   "bbox": [
    -2.04159,
    50.68327,
    -1.89252,
    50.79938
   ]
  },
  {
   "en": "Swindon",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000030",
   "bbox": [
    -1.86655,
    51.48294,
    -1.60426,
    51.6925
   ]
  },
  {
   "en": "Peterborough",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000031",
   "bbox": [
    -0.49927,
    52.50649,
    -0.01452,
    52.67558
   ]
  },
  {
   "en": "Luton",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000032",
   "bbox": [
    -0.50753,
    51.85512,
    -0.35152,
    51.9282
   ]
  },
  {
   "en": "Southend-on-Sea",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000033",
   "bbox": [
    0.62124,
    51.52235,
    0.81951,
    51.57731
   ]
  },
  {
   "en": "Thurrock",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000034",
   "bbox": [
    0.20888,
    51.45266,
    0.53956,
    51.56833
   ]
  },
  {
   "en": "Medway",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000035",
   "bbox": [
    0.39823,
    51.32844,
    0.72152,
    51.48768
   ]
  },
  {
   "en": "Bracknell Forest",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000036",
   "bbox": [
    -0.83889,
    51.33248,
    -0.63212,
    51.46919
   ]
  },
  {
   "en": "West Berkshire",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000037",
   "bbox": [
    -1.58952,
    51.32989,
    -0.98322,
    51.5642
   ]
  },
  {
   "en": "Reading",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000038",
   "bbox": [
    -1.05449,
    51.41029,
    -0.93001,
    51.49359
   ]
  },
  {
   "en": "Slough",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000039",
   "bbox": [
    -0.66171,
    51.46878,
    -0.49161,
    51.5385
   ]
  },
  {
   "en": "Windsor and Maidenhead",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000040",
   "bbox": [
    -0.85545,
    51.38406,
    -0.52434,
    51.57773
   ]
  },
  {
   "en": "Wokingham",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000041",
   "bbox": [
    -1.01341,
    51.35279,
    -0.79038,
    51.56262
   ]
  },
  {
   "en": "Milton Keynes",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000042",
   "bbox": [
    -0.88858,
    51.96967,
    -0.59341,
    52.19674
   ]
  },
  {
   "en": "Brighton and Hove",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000043",
   "bbox": [
    -0.24657,
    50.7998,
    -0.01821,
    50.89292
   ]
  },
  {
   "en": "Portsmouth",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000044",
   "bbox": [
    -1.11753,
    50.75075,
    -1.02212,
    50.85916
   ]
  },
  {
   "en": "Southampton",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000045",
   "bbox": [
    -1.48016,
    50.88061,
    -1.32344,
    50.95669
   ]
  },
  {
   "en": "Isle of Wight",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000046",
   "bbox": [
    -1.59313,
    50.57628,
    -1.07131,
    50.76781
   ]
  },
  {
   "en": "County Durham",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000047",
   "bbox": [
    -2.35694,
    54.45165,
    -1.24353,
    54.91878
   ]
  },
  {
   "en": "Northumberland",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000048",
   "bbox": [
    -2.69124,
    54.78247,
    -1.46336,
    55.81105
   ]
  },
  {
   "en": "Cheshire East",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000049",
   "bbox": [
    -2.75429,
    52.94746,
    -1.97632,
    53.3877
   ]
  },
  {
   "en": "Cheshire West and Chester",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000050",
   "bbox": [
    -3.11091,
    52.98356,
    -2.34777,
    53.34354
   ]
  },
  {
   "en": "Shropshire",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000051",
   "bbox": [
    -3.23682,
    52.30665,
    -2.23431,
    52.9987
   ]
  },
  {
   "en": "Cornwall",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000052",
   "bbox": [
    -5.74777,
    49.9596,
    -4.16772,
    50.93179
   ]
  },
  {
   "en": "Isles of Scilly",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000053",
   "bbox": [
    -6.4194,
    49.86541,
    -6.24616,
    49.98064
   ]
  },
  {
   "en": "Wiltshire",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000054",
   "bbox": [
    -2.36587,
    50.94584,
    -1.48717,
    51.7036
   ]
  },
  {
   "en": "Bedford",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000055",
   "bbox": [
    -0.66973,
    52.05546,
    -0.24237,
    52.32335
   ]
  },
  {
   "en": "Central Bedfordshire",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E06000056",
   "bbox": [
    -0.70374,
    51.80556,
    -0.14605,
    52.19134
   ]
  },
  {
   "en": "Aylesbury Vale",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000004",
   "bbox": [
    -1.1422,
    51.71905,
    -0.53923,
    52.08195
   ]
  },
  {
   "en": "Chiltern",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000005",
   "bbox": [
    -0.76367,
    51.58988,
    -0.5067,
    51.76894
   ]
  },
  {
   "en": "South Bucks",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000006",
   "bbox": [
    -0.70472,
    51.48599,
    -0.47858,
    51.63519
   ]
  },
  {
   "en": "Wycombe",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000007",
   "bbox": [
    -0.95227,
    51.54525,
    -0.67103,
    51.78291
   ]
  },
  {
   "en": "Cambridge",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000008",
   "bbox": [
    0.06811,
    52.15837,
    0.18151,
    52.23765
   ]
  },
  {
   "en": "East Cambridgeshire",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000009",
   "bbox": [
    0.03273,
    52.15685,
    0.51274,
    52.51401
   ]
  },
  {
   "en": "Fenland",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000010",
   "bbox": [
    -0.19953,
    52.38614,
    0.23461,
    52.73966
   ]
  },
  {
   "en": "Huntingdonshire",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000011",
   "bbox": [
    -0.5011,
    52.15916,
    0.04954,
    52.58336
   ]
  },
  {
   "en": "South Cambridgeshire",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000012",
   "bbox": [
    -0.23645,
    52.00623,
    0.41916,
    52.35318
   ]
  },
  {
   "en": "Allerdale",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000026",
   "bbox": [
    -3.58128,
    54.45422,
    -2.98439,
    54.95393
   ]
  },
  {
   "en": "Barrow-in-Furness",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000027",
   "bbox": [
    -3.27802,
    54.04442,
    -3.14394,
    54.21862
   ]
  },
  {
   "en": "Carlisle",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000028",
   "bbox": [
    -3.13251,
    54.77624,
    -2.48442,
    55.18902
   ]
  },
  {
   "en": "Copeland",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000029",
   "bbox": [
    -3.6402,
    54.1887,
    -3.1163,
    54.60744
   ]
  },
  {
   "en": "Eden",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000030",
   "bbox": [
    -3.09437,
    54.35508,
    -2.16085,
    54.85651
   ]
  },
  {
   "en": "South Lakeland",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000031",
   "bbox": [
    -3.24313,
    54.09756,
    -2.3113,
    54.49995
   ]
  },
  {
   "en": "Amber Valley",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000032",
   "bbox": [
    -1.61155,
    52.92451,
    -1.30898,
    53.13425
   ]
  },
  {
   "en": "Bolsover",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000033",
   "bbox": [
    -1.3798,
    53.08128,
    -1.16805,
    53.31278
   ]
  },
  {
   "en": "Chesterfield",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000034",
   "bbox": [
    -1.48568,
    53.21271,
    -1.30284,
    53.29931
   ]
  },
  {
   "en": "Derbyshire Dales",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000035",
   "bbox": [
    -1.92003,
    52.869,
    -1.49688,
    53.39216
   ]
  },
  {
   "en": "Erewash",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000036",
   "bbox": [
    -1.48134,
    52.87277,
    -1.24052,
    53.00395
   ]
  },
  {
   "en": "High Peak",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000037",
   "bbox": [
    -2.03555,
    53.19503,
    -1.65532,
    53.54066
   ]
  },
  {
   "en": "North East Derbyshire",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000038",
   "bbox": [
    -1.6006,
    53.10365,
    -1.28356,
    53.34222
   ]
  },
  {
   "en": "South Derbyshire",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000039",
   "bbox": [
    -1.75097,
    52.69687,
    -1.32084,
    52.94846
   ]
  },
  {
   "en": "East Devon",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000040",
   "bbox": [
    -3.58665,
    50.6073,
    -2.8883,
    50.90881
   ]
  },
  {
   "en": "Exeter",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000041",
   "bbox": [
    -3.57139,
    50.67546,
    -3.45283,
    50.76201
   ]
  },
  {
   "en": "Mid Devon",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000042",
   "bbox": [
    -3.92706,
    50.70375,
    -3.14366,
    51.03436
   ]
  },
  {
   "en": "North Devon",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000043",
   "bbox": [
    -4.26302,
    50.87734,
    -3.59565,
    51.24688
   ]
  },
  {
   "en": "South Hams",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000044",
   "bbox": [
    -4.1909,
    50.20253,
    -3.5088,
    50.54311
   ]
  },
  {
   "en": "Teignbridge",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000045",
   "bbox": [
    -3.88395,
    50.46175,
    -3.42573,
    50.76453
   ]
  },
  {
   "en": "Torridge",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000046",
   "bbox": [
    -4.68161,
    50.647,
    -3.88544,
    51.20296
   ]
  },
  {
   "en": "West Devon",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000047",
   "bbox": [
    -4.3362,
    50.43486,
    -3.73344,
    50.87515
   ]
  },
  {
   "en": "Christchurch",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000048",
   "bbox": [
    -1.87484,
    50.72212,
    -1.68323,
    50.81018
   ]
  },
  {
   "en": "East Dorset",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000049",
   "bbox": [
    -2.16524,
    50.74343,
    -1.79201,
    50.99259
   ]
  },
  {
   "en": "North Dorset",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000050",
   "bbox": [
    -2.48151,
    50.75929,
    -2.04128,
    51.08021
   ]
  },
  {
   "en": "Purbeck",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000051",
   "bbox": [
    -2.32971,
    50.57709,
    -1.92514,
    50.79357
   ]
  },
  {
   "en": "West Dorset",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000052",
   "bbox": [
    -2.9626,
    50.59169,
    -2.26282,
    51.00008
   ]
  },
  {
   "en": "Weymouth and Portland",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000053",
   "bbox": [
    -2.5053,
    50.5137,
    -2.40669,
    50.67921
   ]
  },
  {
   "en": "Eastbourne",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000061",
   "bbox": [
    0.20774,
    50.73555,
    0.33756,
    50.81355
   ]
  },
  {
   "en": "Hastings",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000062",
   "bbox": [
    0.50142,
    50.84373,
    0.65687,
    50.89542
   ]
  },
  {
   "en": "Lewes",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000063",
   "bbox": [
    -0.13689,
    50.75637,
    0.1528,
    51.00309
   ]
  },
  {
   "en": "Rother",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000064",
   "bbox": [
    0.31992,
    50.8209,
    0.86616,
    51.08367
   ]
  },
  {
   "en": "Wealden",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000065",
   "bbox": [
    -0.03782,
    50.73901,
    0.44848,
    51.14757
   ]
  },
  {
   "en": "Basildon",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000066",
   "bbox": [
    0.37343,
    51.53145,
    0.56671,
    51.65149
   ]
  },
  {
   "en": "Braintree",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000067",
   "bbox": [
    0.3779,
    51.7509,
    0.77971,
    52.08748
   ]
  },
  {
   "en": "Brentwood",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000068",
   "bbox": [
    0.17403,
    51.56501,
    0.41086,
    51.7175
   ]
  },
  {
   "en": "Castle Point",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000069",
   "bbox": [
    0.51629,
    51.50835,
    0.63514,
    51.58875
   ]
  },
  {
   "en": "Chelmsford",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000070",
   "bbox": [
    0.3305,
    51.61763,
    0.64497,
    51.85773
   ]
  },
  {
   "en": "Colchester",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000071",
   "bbox": [
    0.69785,
    51.7723,
    1.02292,
    51.97739
   ]
  },
  {
   "en": "Epping Forest",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000072",
   "bbox": [
    -0.02138,
    51.60497,
    0.34575,
    51.82263
   ]
  },
  {
   "en": "Harlow",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000073",
   "bbox": [
    0.0524,
    51.73517,
    0.16634,
    51.79672
   ]
  },
  {
   "en": "Maldon",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000074",
   "bbox": [
    0.58777,
    51.62174,
    0.94892,
    51.8296
   ]
  },
  {
   "en": "Rochford",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000075",
   "bbox": [
    0.54807,
    51.54065,
    0.95575,
    51.63738
   ]
  },
  {
   "en": "Tendring",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000076",
   "bbox": [
    0.92974,
    51.7704,
    1.29481,
    51.95979
   ]
  },
  {
   "en": "Uttlesford",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000077",
   "bbox": [
    0.06645,
    51.76808,
    0.51754,
    52.09311
   ]
  },
  {
   "en": "Cheltenham",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000078",
   "bbox": [
    -2.14381,
    51.85877,
    -2.01163,
    51.9393
   ]
  },
  {
   "en": "Cotswold",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000079",
   "bbox": [
    -2.32496,
    51.57806,
    -1.61665,
    52.11299
   ]
  },
  {
   "en": "Forest of Dean",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000080",
   "bbox": [
    -2.68841,
    51.61021,
    -2.26893,
    52.02413
   ]
  },
  {
   "en": "Gloucester",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000081",
   "bbox": [
    -2.30247,
    51.80802,
    -2.17865,
    51.88548
   ]
  },
  {
   "en": "Stroud",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000082",
   "bbox": [
    -2.53608,
    51.5903,
    -2.07002,
    51.85053
   ]
  },
  {
   "en": "Tewkesbury",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000083",
   "bbox": [
    -2.35411,
    51.82,
    -1.80312,
    52.05061
   ]
  },
  {
   "en": "Basingstoke and Deane",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000084",
   "bbox": [
    -1.46161,
    51.13408,
    -0.97631,
    51.38443
   ]
  },
  {
   "en": "East Hampshire",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000085",
   "bbox": [
    -1.1333,
    50.86783,
    -0.74622,
    51.21327
   ]
  },
  {
   "en": "Eastleigh",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000086",
   "bbox": [
    -1.4011,
    50.84954,
    -1.25962,
    51.00492
   ]
  },
  {
   "en": "Fareham",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000087",
   "bbox": [
    -1.30946,
    50.80988,
    -1.11442,
    50.89832
   ]
  },
  {
   "en": "Gosport",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000088",
   "bbox": [
    -1.21562,
    50.77442,
    -1.11373,
    50.83987
   ]
  },
  {
   "en": "Hart",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000089",
   "bbox": [
    -1.00312,
    51.18654,
    -0.76445,
    51.36683
   ]
  },
  {
   "en": "Havant",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000090",
   "bbox": [
    -1.05611,
    50.77772,
    -0.92796,
    50.91004
   ]
  },
  {
   "en": "New Forest",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000091",
   "bbox": [
    -1.95821,
    50.70658,
    -1.30837,
    51.00997
   ]
  },
  {
   "en": "Rushmoor",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000092",
   "bbox": [
    -0.80906,
    51.23099,
    -0.73204,
    51.32061
   ]
  },
  {
   "en": "Test Valley",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000093",
   "bbox": [
    -1.69549,
    50.92887,
    -1.3103,
    51.33987
   ]
  },
  {
   "en": "Winchester",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000094",
   "bbox": [
    -1.45795,
    50.85569,
    -1.03561,
    51.19728
   ]
  },
  {
   "en": "Broxbourne",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000095",
   "bbox": [
    -0.11572,
    51.68137,
    0.01237,
    51.78099
   ]
  },
  {
   "en": "Dacorum",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000096",
   "bbox": [
    -0.74725,
    51.68031,
    -0.40647,
    51.85846
   ]
  },
  {
   "en": "East Hertfordshire",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000097",
   "bbox": [
    -0.18524,
    51.73518,
    0.19391,
    51.99762
   ]
  },
  {
   "en": "Hertsmere",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000098",
   "bbox": [
    -0.38818,
    51.62399,
    -0.16475,
    51.73862
   ]
  },
  {
   "en": "North Hertfordshire",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000099",
   "bbox": [
    -0.40744,
    51.8337,
    0.07264,
    52.08098
   ]
  },
  {
   "en": "St Albans",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000100",
   "bbox": [
    -0.44219,
    51.68807,
    -0.24362,
    51.85007
   ]
  },
  {
   "en": "Stevenage",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000101",
   "bbox": [
    -0.23601,
    51.86973,
    -0.15007,
    51.93287
   ]
  },
  {
   "en": "Three Rivers",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000102",
   "bbox": [
    -0.54085,
    51.60018,
    -0.36423,
    51.74131
   ]
  },
  {
   "en": "Watford",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000103",
   "bbox": [
    -0.44114,
    51.63821,
    -0.37038,
    51.70242
   ]
  },
  {
   "en": "Welwyn Hatfield",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000104",
   "bbox": [
    -0.27917,
    51.686,
    -0.09379,
    51.86091
   ]
  },
  {
   "en": "Ashford",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000105",
   "bbox": [
    0.58867,
    50.99008,
    1.02749,
    51.2708
   ]
  },
  {
   "en": "Canterbury",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000106",
   "bbox": [
    0.94849,
    51.17774,
    1.24746,
    51.3804
   ]
  },
  {
   "en": "Dartford",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000107",
   "bbox": [
    0.1475,
    51.38708,
    0.34256,
    51.48003
   ]
  },
  {
   "en": "Dover",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000108",
   "bbox": [
    1.14197,
    51.09777,
    1.40387,
    51.33227
   ]
  },
  {
   "en": "Gravesham",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000109",
   "bbox": [
    0.309,
    51.32555,
    0.48949,
    51.46547
   ]
  },
  {
   "en": "Maidstone",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000110",
   "bbox": [
    0.37698,
    51.13451,
    0.79474,
    51.33928
   ]
  },
  {
   "en": "Sevenoaks",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000111",
   "bbox": [
    0.03196,
    51.13235,
    0.34304,
    51.41748
   ]
  },
  {
   "en": "Shepway",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000112",
   "bbox": [
    0.77531,
    50.9129,
    1.21936,
    51.20557
   ]
  },
  {
   "en": "Swale",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000113",
   "bbox": [
    0.59927,
    51.22499,
    1.01523,
    51.44782
   ]
  },
  {
   "en": "Thanet",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000114",
   "bbox": [
    1.21109,
    51.30854,
    1.44778,
    51.39433
   ]
  },
  {
   "en": "Tonbridge and Malling",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000115",
   "bbox": [
    0.19888,
    51.17648,
    0.52507,
    51.36859
   ]
  },
  {
   "en": "Tunbridge Wells",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000116",
   "bbox": [
    0.14841,
    51.00404,
    0.64497,
    51.20217
   ]
  },
  {
   "en": "Burnley",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000117",
   "bbox": [
    -2.34374,
    53.72322,
    -2.11375,
    53.82539
   ]
  },
  {
   "en": "Chorley",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000118",
   "bbox": [
    -2.82543,
    53.5939,
    -2.51274,
    53.75147
   ]
  },
  {
   "en": "Fylde",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000119",
   "bbox": [
    -3.05821,
    53.73165,
    -2.78189,
    53.86459
   ]
  },
  {
   "en": "Hyndburn",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000120",
   "bbox": [
    -2.46723,
    53.70534,
    -2.31517,
    53.81681
   ]
  },
  {
   "en": "Lancaster",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000121",
   "bbox": [
    -2.92581,
    53.91825,
    -2.46103,
    54.23971
   ]
  },
  {
   "en": "Pendle",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000122",
   "bbox": [
    -2.33509,
    53.80591,
    -2.04761,
    53.95244
   ]
  },
  {
   "en": "Preston",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000123",
   "bbox": [
    -2.82738,
    53.74887,
    -2.59685,
    53.89635
   ]
  },
  {
   "en": "Ribble Valley",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000124",
   "bbox": [
    -2.65224,
    53.75658,
    -2.18598,
    54.04926
   ]
  },
  {
   "en": "Rossendale",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000125",
   "bbox": [
    -2.41242,
    53.61473,
    -2.14671,
    53.75525
   ]
  },
  {
   "en": "South Ribble",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000126",
   "bbox": [
    -2.85621,
    53.67162,
    -2.54696,
    53.78299
   ]
  },
  {
   "en": "West Lancashire",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000127",
   "bbox": [
    -3.04805,
    53.48301,
    -2.6907,
    53.73412
   ]
  },
  {
   "en": "Wyre",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000128",
   "bbox": [
    -3.05154,
    53.82017,
    -2.61481,
    53.98041
   ]
  },
  {
   "en": "Blaby",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000129",
   "bbox": [
    -1.3394,
    52.49376,
    -1.06094,
    52.66474
   ]
  },
  {
   "en": "Charnwood",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000130",
   "bbox": [
    -1.3363,
    52.65433,
    -0.94881,
    52.82505
   ]
  },
  {
   "en": "Harborough",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000131",
   "bbox": [
    -1.30747,
    52.39255,
    -0.71525,
    52.68532
   ]
  },
  {
   "en": "Hinckley and Bosworth",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000132",
   "bbox": [
    -1.57249,
    52.5014,
    -1.20122,
    52.71498
   ]
  },
  {
   "en": "Melton",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000133",
   "bbox": [
    -1.04818,
    52.64647,
    -0.66571,
    52.97723
   ]
  },
  {
   "en": "North West Leicestershire",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000134",
   "bbox": [
    -1.59903,
    52.6643,
    -1.24868,
    52.87739
   ]
  },
  {
   "en": "Oadby and Wigston",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000135",
   "bbox": [
    -1.14516,
    52.55693,
    -1.04373,
    52.61841
   ]
  },
  {
   "en": "Boston",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000136",
   "bbox": [
    -0.24837,
    52.86323,
    0.19786,
    53.087
   ]
  },
  {
   "en": "East Lindsey",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000137",
   "bbox": [
    -0.31912,
    53.00143,
    0.35389,
    53.52697
   ]
  },
  {
   "en": "Lincoln",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000138",
   "bbox": [
    -0.62424,
    53.18661,
    -0.49676,
    53.25484
   ]
  },
  {
   "en": "North Kesteven",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000139",
   "bbox": [
    -0.76574,
    52.88732,
    -0.19156,
    53.25841
   ]
  },
  {
   "en": "South Holland",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000140",
   "bbox": [
    -0.30846,
    52.65187,
    0.27049,
    52.92915
   ]
  },
  {
   "en": "South Kesteven",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000141",
   "bbox": [
    -0.80587,
    52.64059,
    -0.21416,
    53.06022
   ]
  },
  {
   "en": "West Lindsey",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000142",
   "bbox": [
    -0.82157,
    53.17966,
    -0.13357,
    53.61663
   ]
  },
  {
   "en": "Breckland",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000143",
   "bbox": [
    0.52877,
    52.37052,
    1.1065,
    52.81148
   ]
  },
  {
   "en": "Broadland",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000144",
   "bbox": [
    0.98856,
    52.55591,
    1.67493,
    52.83076
   ]
  },
  {
   "en": "Great Yarmouth",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000145",
   "bbox": [
    1.5436,
    52.52652,
    1.74391,
    52.74295
   ]
  },
  {
   "en": "King's Lynn and West Norfolk",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000146",
   "bbox": [
    0.15236,
    52.43682,
    0.81978,
    52.98873
   ]
  },
  {
   "en": "North Norfolk",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000147",
   "bbox": [
    0.70616,
    52.68028,
    1.67298,
    52.98036
   ]
  },
  {
   "en": "Norwich",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000148",
   "bbox": [
    1.20206,
    52.59656,
    1.3404,
    52.68531
   ]
  },
  {
   "en": "South Norfolk",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000149",
   "bbox": [
    0.94523,
    52.35583,
    1.68018,
    52.67844
   ]
  },
  {
   "en": "Corby",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000150",
   "bbox": [
    -0.80951,
    52.45485,
    -0.598,
    52.559
   ]
  },
  {
   "en": "Daventry",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000151",
   "bbox": [
    -1.28616,
    52.14362,
    -0.78736,
    52.47758
   ]
  },
  {
   "en": "East Northamptonshire",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000152",
   "bbox": [
    -0.67646,
    52.25387,
    -0.34323,
    52.64273
   ]
  },
  {
   "en": "Kettering",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000153",
   "bbox": [
    -0.90782,
    52.3476,
    -0.61538,
    52.52869
   ]
  },
  {
   "en": "Northampton",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000154",
   "bbox": [
    -0.97258,
    52.18651,
    -0.79299,
    52.28291
   ]
  },
  {
   "en": "South Northamptonshire",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000155",
   "bbox": [
    -1.33341,
    51.97786,
    -0.70704,
    52.25987
   ]
  },
  {
   "en": "Wellingborough",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000156",
   "bbox": [
    -0.81953,
    52.19198,
    -0.61222,
    52.36462
   ]
  },
  {
   "en": "Craven",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000163",
   "bbox": [
    -2.56484,
    53.85033,
    -1.81824,
    54.2573
   ]
  },
  {
   "en": "Hambleton",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000164",
   "bbox": [
    -1.70679,
    53.98937,
    -0.97371,
    54.51107
   ]
  },
  {
   "en": "Harrogate",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000165",
   "bbox": [
    -2.006,
    53.89163,
    -1.17709,
    54.26175
   ]
  },
  {
   "en": "Richmondshire",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000166",
   "bbox": [
    -2.36917,
    54.17295,
    -1.46918,
    54.54242
   ]
  },
  {
   "en": "Ryedale",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000167",
   "bbox": [
    -1.26722,
    53.98174,
    -0.41978,
    54.41549
   ]
  },
  {
   "en": "Scarborough",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000168",
   "bbox": [
    -1.06631,
    54.13262,
    -0.21424,
    54.55848
   ]
  },
  {
   "en": "Selby",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000169",
   "bbox": [
    -1.35453,
    53.62134,
    -0.90554,
    53.93579
   ]
  },
  {
   "en": "Ashfield",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000170",
   "bbox": [
    -1.34564,
    53.0084,
    -1.16581,
    53.17174
   ]
  },
  {
   "en": "Bassetlaw",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000171",
   "bbox": [
    -1.21401,
    53.19695,
    -0.74948,
    53.50276
   ]
  },
  {
   "en": "Broxtowe",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000172",
   "bbox": [
    -1.3374,
    52.8926,
    -1.18385,
    53.05477
   ]
  },
  {
   "en": "Gedling",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000173",
   "bbox": [
    -1.23283,
    52.94888,
    -1.00866,
    53.10425
   ]
  },
  {
   "en": "Mansfield",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000174",
   "bbox": [
    -1.26161,
    53.1151,
    -1.09682,
    53.23588
   ]
  },
  {
   "en": "Newark and Sherwood",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000175",
   "bbox": [
    -1.17944,
    52.95777,
    -0.6682,
    53.26112
   ]
  },
  {
   "en": "Rushcliffe",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000176",
   "bbox": [
    -1.27682,
    52.79049,
    -0.81675,
    53.03623
   ]
  },
  {
   "en": "Cherwell",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000177",
   "bbox": [
    -1.52431,
    51.78142,
    -1.04866,
    52.16889
   ]
  },
  {
   "en": "Oxford",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000178",
   "bbox": [
    -1.30517,
    51.71145,
    -1.17731,
    51.79672
   ]
  },
  {
   "en": "South Oxfordshire",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000179",
   "bbox": [
    -1.29172,
    51.46002,
    -0.8716,
    51.81423
   ]
  },
  {
   "en": "Vale of White Horse",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000180",
   "bbox": [
    -1.70164,
    51.51878,
    -1.20378,
    51.79011
   ]
  },
  {
   "en": "West Oxfordshire",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000181",
   "bbox": [
    -1.72094,
    51.68431,
    -1.28764,
    51.99727
   ]
  },
  {
   "en": "Mendip",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000187",
   "bbox": [
    -2.84397,
    51.06384,
    -2.24579,
    51.32621
   ]
  },
  {
   "en": "Sedgemoor",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000188",
   "bbox": [
    -3.2177,
    51.0414,
    -2.71277,
    51.32904
   ]
  },
  {
   "en": "South Somerset",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000189",
   "bbox": [
    -3.09363,
    50.82173,
    -2.3272,
    51.14824
   ]
  },
  {
   "en": "Taunton Deane",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000190",
   "bbox": [
    -3.4169,
    50.89194,
    -2.88511,
    51.11858
   ]
  },
  {
   "en": "West Somerset",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000191",
   "bbox": [
    -3.84097,
    51.00387,
    -3.05285,
    51.23375
   ]
  },
  {
   "en": "Cannock Chase",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000192",
   "bbox": [
    -2.05607,
    52.64074,
    -1.91226,
    52.77408
   ]
  },
  {
   "en": "East Staffordshire",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000193",
   "bbox": [
    -2.04409,
    52.73018,
    -1.58947,
    53.04513
   ]
  },
  {
   "en": "Lichfield",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000194",
   "bbox": [
    -1.96467,
    52.58529,
    -1.58778,
    52.80774
   ]
  },
  {
   "en": "Newcastle-under-Lyme",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000195",
   "bbox": [
    -2.47223,
    52.87431,
    -2.18264,
    53.1161
   ]
  },
  {
   "en": "South Staffordshire",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000196",
   "bbox": [
    -2.32687,
    52.42362,
    -1.98518,
    52.78652
   ]
  },
  {
   "en": "Stafford",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000197",
   "bbox": [
    -2.41774,
    52.71549,
    -1.93503,
    52.98087
   ]
  },
  {
   "en": "Staffordshire Moorlands",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000198",
   "bbox": [
    -2.21274,
    52.91778,
    -1.77713,
    53.22651
   ]
  },
  {
   "en": "Tamworth",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000199",
   "bbox": [
    -1.73724,
    52.58969,
    -1.63573,
    52.65785
   ]
  },
  {
   "en": "Babergh",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000200",
   "bbox": [
    0.62517,
    51.94892,
    1.27727,
    52.18151
   ]
  },
  {
   "en": "Forest Heath",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000201",
   "bbox": [
    0.33828,
    52.20324,
    0.71702,
    52.4629
   ]
  },
  {
   "en": "Ipswich",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000202",
   "bbox": [
    1.10548,
    52.02123,
    1.2218,
    52.09492
   ]
  },
  {
   "en": "Mid Suffolk",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000203",
   "bbox": [
    0.79835,
    52.06584,
    1.40766,
    52.40482
   ]
  },
  {
   "en": "St Edmundsbury",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000204",
   "bbox": [
    0.38008,
    52.05429,
    0.96567,
    52.40113
   ]
  },
  {
   "en": "Suffolk Coastal",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000205",
   "bbox": [
    1.15603,
    51.93286,
    1.66531,
    52.36926
   ]
  },
  {
   "en": "Waveney",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000206",
   "bbox": [
    1.34447,
    52.31209,
    1.76168,
    52.55005
   ]
  },
  {
   "en": "Elmbridge",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000207",
   "bbox": [
    -0.48274,
    51.29542,
    -0.30991,
    51.41242
   ]
  },
  {
   "en": "Epsom and Ewell",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000208",
   "bbox": [
    -0.30779,
    51.29287,
    -0.21888,
    51.38056
   ]
  },
  {
   "en": "Guildford",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000209",
   "bbox": [
    -0.74988,
    51.17373,
    -0.38873,
    51.33203
   ]
  },
  {
   "en": "Mole Valley",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000210",
   "bbox": [
    -0.44002,
    51.10576,
    -0.17832,
    51.33561
   ]
  },
  {
   "en": "Reigate and Banstead",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000211",
   "bbox": [
    -0.27384,
    51.16038,
    -0.12592,
    51.34412
   ]
  },
  {
   "en": "Runnymede",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000212",
   "bbox": [
    -0.62006,
    51.33984,
    -0.45957,
    51.45152
   ]
  },
  {
   "en": "Spelthorne",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000213",
   "bbox": [
    -0.54219,
    51.37896,
    -0.38494,
    51.47204
   ]
  },
  {
   "en": "Surrey Heath",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000214",
   "bbox": [
    -0.77701,
    51.27978,
    -0.55012,
    51.3929
   ]
  },
  {
   "en": "Tandridge",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000215",
   "bbox": [
    -0.15594,
    51.13732,
    0.05596,
    51.33918
   ]
  },
  {
   "en": "Waverley",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000216",
   "bbox": [
    -0.85044,
    51.07205,
    -0.41482,
    51.2456
   ]
  },
  {
   "en": "Woking",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000217",
   "bbox": [
    -0.64945,
    51.26604,
    -0.46393,
    51.34988
   ]
  },
  {
   "en": "North Warwickshire",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000218",
   "bbox": [
    -1.78955,
    52.43584,
    -1.4615,
    52.68759
   ]
  },
  {
   "en": "Nuneaton and Bedworth",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000219",
   "bbox": [
    -1.55671,
    52.4514,
    -1.4054,
    52.55194
   ]
  },
  {
   "en": "Rugby",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000220",
   "bbox": [
    -1.46648,
    52.25365,
    -1.17366,
    52.53481
   ]
  },
  {
   "en": "Stratford-on-Avon",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000221",
   "bbox": [
    -1.9631,
    51.95583,
    -1.23452,
    52.36797
   ]
  },
  {
   "en": "Warwick",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000222",
   "bbox": [
    -1.78068,
    52.2142,
    -1.40814,
    52.38967
   ]
  },
  {
   "en": "Adur",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000223",
   "bbox": [
    -0.37309,
    50.81804,
    -0.21763,
    50.87511
   ]
  },
  {
   "en": "Arun",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000224",
   "bbox": [
    -0.76516,
    50.75925,
    -0.36507,
    50.90626
   ]
  },
  {
   "en": "Chichester",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000225",
   "bbox": [
    -0.95909,
    50.72296,
    -0.47477,
    51.09504
   ]
  },
  {
   "en": "Crawley",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000226",
   "bbox": [
    -0.25677,
    51.08671,
    -0.13458,
    51.1674
   ]
  },
  {
   "en": "Horsham",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000227",
   "bbox": [
    -0.57047,
    50.86292,
    -0.20161,
    51.14367
   ]
  },
  {
   "en": "Mid Sussex",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000228",
   "bbox": [
    -0.25068,
    50.86847,
    0.04292,
    51.14319
   ]
  },
  {
   "en": "Worthing",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000229",
   "bbox": [
    -0.4481,
    50.80373,
    -0.33428,
    50.86381
   ]
  },
  {
   "en": "Bromsgrove",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000234",
   "bbox": [
    -2.16912,
    52.27848,
    -1.84705,
    52.44564
   ]
  },
  {
   "en": "Malvern Hills",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000235",
   "bbox": [
    -2.66455,
    51.96698,
    -2.1502,
    52.36862
   ]
  },
  {
   "en": "Redditch",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000236",
   "bbox": [
    -2.02033,
    52.23566,
    -1.87665,
    52.32455
   ]
  },
  {
   "en": "Worcester",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000237",
   "bbox": [
    -2.26327,
    52.16239,
    -2.15919,
    52.23176
   ]
  },
  {
   "en": "Wychavon",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000238",
   "bbox": [
    -2.27628,
    52.00012,
    -1.75929,
    52.36144
   ]
  },
  {
   "en": "Wyre Forest",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E07000239",
   "bbox": [
    -2.43789,
    52.31522,
    -2.12063,
    52.45567
   ]
  },
  {
   "en": "Bolton",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E08000001",
   "bbox": [
    -2.62891,
    53.52317,
    -2.33963,
    53.64626
   ]
  },
  {
   "en": "Bury",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E08000002",
   "bbox": [
    -2.38489,
    53.51225,
    -2.23723,
    53.66729
   ]
  },
  {
   "en": "Manchester",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E08000003",
   "bbox": [
    -2.32135,
    53.34047,
    -2.15056,
    53.54438
   ]
  },
  {
   "en": "Oldham",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E08000004",
   "bbox": [
    -2.18746,
    53.49214,
    -1.9111,
    53.62439
   ]
  },
  {
   "en": "Rochdale",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E08000005",
   "bbox": [
    -2.28388,
    53.52929,
    -2.02829,
    53.68594
   ]
  },
  {
   "en": "Salford",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E08000006",
   "bbox": [
    -2.49112,
    53.41613,
    -2.24667,
    53.54221
   ]
  },
  {
   "en": "Stockport",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E08000007",
   "bbox": [
    -2.24826,
    53.32835,
    -1.9938,
    53.45517
   ]
  },
  {
   "en": "Tameside",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E08000008",
   "bbox": [
    -2.16899,
    53.42707,
    -1.96486,
    53.5309
   ]
  },
  {
   "en": "Trafford",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E08000009",
   "bbox": [
    -2.47986,
    53.35767,
    -2.25513,
    53.48022
   ]
  },
  {
   "en": "Wigan",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E08000010",
   "bbox": [
    -2.7319,
    53.44629,
    -2.41677,
    53.6085
   ]
  },
  {
   "en": "Knowsley",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E08000011",
   "bbox": [
    -2.92397,
    53.34757,
    -2.74477,
    53.50405
   ]
  },
  {
   "en": "Liverpool",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E08000012",
   "bbox": [
    -3.0101,
    53.32726,
    -2.82014,
    53.47521
   ]
  },
  {
   "en": "St. Helens",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E08000013",
   "bbox": [
    -2.82633,
    53.38563,
    -2.57814,
    53.53165
   ]
  },
  {
   "en": "Sefton",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E08000014",
   "bbox": [
    -3.10681,
    53.43861,
    -2.88257,
    53.6984
   ]
  },
  {
   "en": "Wirral",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E08000015",
   "bbox": [
    -3.22959,
    53.29653,
    -2.92992,
    53.44312
   ]
  },
  {
   "en": "Barnsley",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E08000016",
   "bbox": [
    -1.82372,
    53.43856,
    -1.27728,
    53.61297
   ]
  },
  {
   "en": "Doncaster",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E08000017",
   "bbox": [
    -1.35028,
    53.40614,
    -0.86694,
    53.66141
   ]
  },
  {
   "en": "Rotherham",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E08000018",
   "bbox": [
    -1.45675,
    53.30191,
    -1.11715,
    53.51546
   ]
  },
  {
   "en": "Sheffield",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E08000019",
   "bbox": [
    -1.80296,
    53.30501,
    -1.32621,
    53.50336
   ]
  },
  {
   "en": "Gateshead",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E08000020",
   "bbox": [
    -1.85426,
    54.87785,
    -1.51279,
    54.98459
   ]
  },
  {
   "en": "Newcastle upon Tyne",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E08000021",
   "bbox": [
    -1.77723,
    54.96008,
    -1.53242,
    55.07945
   ]
  },
  {
   "en": "North Tyneside",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E08000022",
   "bbox": [
    -1.64065,
    54.98339,
    -1.40439,
    55.07455
   ]
  },
  {
   "en": "South Tyneside",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E08000023",
   "bbox": [
    -1.53575,
    54.92845,
    -1.35481,
    55.01139
   ]
  },
  {
   "en": "Sunderland",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E08000024",
   "bbox": [
    -1.57046,
    54.80009,
    -1.34904,
    54.9442
   ]
  },
  {
   "en": "Birmingham",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E08000025",
   "bbox": [
    -2.03364,
    52.38143,
    -1.73032,
    52.60903
   ]
  },
  {
   "en": "Coventry",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E08000026",
   "bbox": [
    -1.61593,
    52.36541,
    -1.42565,
    52.46516
   ]
  },
  {
   "en": "Dudley",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E08000027",
   "bbox": [
    -2.19316,
    52.42641,
    -2.01447,
    52.55855
   ]
  },
  {
   "en": "Sandwell",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E08000028",
   "bbox": [
    -2.0985,
    52.4611,
    -1.9196,
    52.56941
   ]
  },
  {
   "en": "Solihull",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E08000029",
   "bbox": [
    -1.87348,
    52.34834,
    -1.59671,
    52.5146
   ]
  },
  {
   "en": "Walsall",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E08000030",
   "bbox": [
    -2.07922,
    52.54765,
    -1.87402,
    52.66297
   ]
  },
  {
   "en": "Wolverhampton",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E08000031",
   "bbox": [
    -2.20829,
    52.54447,
    -2.04946,
    52.63804
   ]
  },
  {
   "en": "Bradford",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E08000032",
   "bbox": [
    -2.06272,
    53.72457,
    -1.64214,
    53.96332
   ]
  },
  {
   "en": "Calderdale",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E08000033",
   "bbox": [
    -2.17475,
    53.61606,
    -1.72873,
    53.82583
   ]
  },
  {
   "en": "Kirklees",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E08000034",
   "bbox": [
    -2.01094,
    53.52014,
    -1.57265,
    53.76506
   ]
  },
  {
   "en": "Leeds",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E08000035",
   "bbox": [
    -1.80193,
    53.6992,
    -1.29352,
    53.94608
   ]
  },
  {
   "en": "Wakefield",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E08000036",
   "bbox": [
    -1.62561,
    53.5756,
    -1.20039,
    53.74194
   ]
  },
  {
   "en": "City of London",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E09000001",
   "bbox": [
    -0.11543,
    51.50835,
    -0.07468,
    51.5221
   ]
  },
  {
   "en": "Barking and Dagenham",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E09000002",
   "bbox": [
    0.06664,
    51.5127,
    0.18854,
    51.59946
   ]
  },
  {
   "en": "Barnet",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E09000003",
   "bbox": [
    -0.30607,
    51.55569,
    -0.13076,
    51.67066
   ]
  },
  {
   "en": "Bexley",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E09000004",
   "bbox": [
    0.07372,
    51.40901,
    0.21597,
    51.51381
   ]
  },
  {
   "en": "Brent",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E09000005",
   "bbox": [
    -0.33717,
    51.52816,
    -0.19308,
    51.60087
   ]
  },
  {
   "en": "Bromley",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E09000006",
   "bbox": [
    -0.0826,
    51.28989,
    0.15965,
    51.44484
   ]
  },
  {
   "en": "Camden",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E09000007",
   "bbox": [
    -0.2151,
    51.51394,
    -0.10696,
    51.57293
   ]
  },
  {
   "en": "Croydon",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E09000008",
   "bbox": [
    -0.1635,
    51.2873,
    0.00167,
    51.42376
   ]
  },
  {
   "en": "Ealing",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E09000009",
   "bbox": [
    -0.42011,
    51.49098,
    -0.24668,
    51.56018
   ]
  },
  {
   "en": "Enfield",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E09000010",
   "bbox": [
    -0.1875,
    51.60613,
    -0.01148,
    51.69236
   ]
  },
  {
   "en": "Greenwich",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E09000011",
   "bbox": [
    -0.02646,
    51.42473,
    0.12253,
    51.51196
   ]
  },
  {
   "en": "Hackney",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E09000012",
   "bbox": [
    -0.10612,
    51.52084,
    -0.01819,
    51.57829
   ]
  },
  {
   "en": "Hammersmith and Fulham",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E09000013",
   "bbox": [
    -0.25668,
    51.46572,
    -0.18102,
    51.53326
   ]
  },
  {
   "en": "Haringey",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E09000014",
   "bbox": [
    -0.17289,
    51.56528,
    -0.04307,
    51.61171
   ]
  },
  {
   "en": "Harrow",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E09000015",
   "bbox": [
    -0.40565,
    51.55356,
    -0.26875,
    51.64102
   ]
  },
  {
   "en": "Havering",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E09000016",
   "bbox": [
    0.13652,
    51.48838,
    0.33221,
    51.63224
   ]
  },
  {
   "en": "Hillingdon",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E09000017",
   "bbox": [
    -0.51128,
    51.45378,
    -0.37796,
    51.63218
   ]
  },
  {
   "en": "Hounslow",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E09000018",
   "bbox": [
    -0.46291,
    51.42179,
    -0.24601,
    51.50336
   ]
  },
  {
   "en": "Islington",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E09000019",
   "bbox": [
    -0.14402,
    51.51905,
    -0.07845,
    51.57524
   ]
  },
  {
   "en": "Kensington and Chelsea",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E09000020",
   "bbox": [
    -0.2301,
    51.47793,
    -0.15161,
    51.53086
   ]
  },
  {
   "en": "Kingston upon Thames",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E09000021",
   "bbox": [
    -0.33226,
    51.32725,
    -0.2413,
    51.43781
   ]
  },
  {
   "en": "Lambeth",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E09000022",
   "bbox": [
    -0.15222,
    51.41284,
    -0.07992,
    51.50895
   ]
  },
  {
   "en": "Lewisham",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E09000023",
   "bbox": [
    -0.07554,
    51.41501,
    0.03741,
    51.49357
   ]
  },
  {
   "en": "Merton",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E09000024",
   "bbox": [
    -0.25584,
    51.38067,
    -0.12702,
    51.44199
   ]
  },
  {
   "en": "Newham",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E09000025",
   "bbox": [
    -0.02279,
    51.49882,
    0.09513,
    51.5645
   ]
  },
  {
   "en": "Redbridge",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E09000026",
   "bbox": [
    0.00705,
    51.54457,
    0.14731,
    51.62933
   ]
  },
  {
   "en": "Richmond upon Thames",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E09000027",
   "bbox": [
    -0.39294,
    51.39196,
    -0.22502,
    51.48946
   ]
  },
  {
   "en": "Southwark",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E09000028",
   "bbox": [
    -0.11257,
    51.42113,
    -0.03403,
    51.50895
   ]
  },
  {
   "en": "Sutton",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E09000029",
   "bbox": [
    -0.24702,
    51.32204,
    -0.11891,
    51.39394
   ]
  },
  {
   "en": "Tower Hamlets",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E09000030",
   "bbox": [
    -0.08099,
    51.48666,
    0.00713,
    51.54519
   ]
  },
  {
   "en": "Waltham Forest",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E09000031",
   "bbox": [
    -0.0628,
    51.55043,
    0.02357,
    51.64672
   ]
  },
  {
   "en": "Wandsworth",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E09000032",
   "bbox": [
    -0.2607,
    51.41827,
    -0.13006,
    51.48556
   ]
  },
  {
   "en": "Westminster",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "E09000033",
   "bbox": [
    -0.21763,
    51.48533,
    -0.11318,
    51.5403
   ]
  },
  {
   "en": "Isle of Anglesey",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "W06000001",
   "bbox": [
    -4.70157,
    53.12694,
    -4.02214,
    53.43592
   ]
  },
  {
   "en": "Gwynedd",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "W06000002",
   "bbox": [
    -4.80528,
    52.54153,
    -3.43806,
    53.24862
   ]
  },
  {
   "en": "Conwy",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "W06000003",
   "bbox": [
    -4.03233,
    52.94965,
    -3.45625,
    53.34299
   ]
  },
  {
   "en": "Denbighshire",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "W06000004",
   "bbox": [
    -3.60288,
    52.86223,
    -3.09161,
    53.35229
   ]
  },
  {
   "en": "Flintshire",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "W06000005",
   "bbox": [
    -3.40189,
    53.07264,
    -2.92166,
    53.35693
   ]
  },
  {
   "en": "Wrexham",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "W06000006",
   "bbox": [
    -3.37629,
    52.8669,
    -2.72554,
    53.13458
   ]
  },
  {
   "en": "Ceredigion",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "W06000008",
   "bbox": [
    -4.69707,
    52.02714,
    -3.65939,
    52.56116
   ]
  },
  {
   "en": "Pembrokeshire",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "W06000009",
   "bbox": [
    -5.67109,
    51.59654,
    -4.48671,
    52.11826
   ]
  },
  {
   "en": "Carmarthenshire",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "W06000010",
   "bbox": [
    -4.72417,
    51.65524,
    -3.64834,
    52.14257
   ]
  },
  {
   "en": "Swansea",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "W06000011",
   "bbox": [
    -4.33422,
    51.53622,
    -3.84426,
    51.77457
   ]
  },
  {
   "en": "Neath Port Talbot",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "W06000012",
   "bbox": [
    -3.93766,
    51.52785,
    -3.56452,
    51.81078
   ]
  },
  {
   "en": "Bridgend",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "W06000013",
   "bbox": [
    -3.76386,
    51.47066,
    -3.46424,
    51.64593
   ]
  },
  {
   "en": "Vale of Glamorgan",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "W06000014",
   "bbox": [
    -3.64305,
    51.38167,
    -3.16682,
    51.51548
   ]
  },
  {
   "en": "Cardiff",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "W06000015",
   "bbox": [
    -3.345,
    51.37556,
    -3.07017,
    51.56096
   ]
  },
  {
   "en": "Rhondda Cynon Taf",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "W06000016",
   "bbox": [
    -3.59492,
    51.49943,
    -3.23804,
    51.83055
   ]
  },
  {
   "en": "Caerphilly",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "W06000018",
   "bbox": [
    -3.33571,
    51.5461,
    -3.06497,
    51.7994
   ]
  },
  {
   "en": "Blaenau Gwent",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "W06000019",
   "bbox": [
    -3.31134,
    51.68237,
    -3.10728,
    51.82591
   ]
  },
  {
   "en": "Torfaen",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "W06000020",
   "bbox": [
    -3.14519,
    51.60713,
    -2.96019,
    51.79591
   ]
  },
  {
   "en": "Monmouthshire",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "W06000021",
   "bbox": [
    -3.15861,
    51.52577,
    -2.65149,
    51.98355
   ]
  },
  {
   "en": "Newport",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "W06000022",
   "bbox": [
    -3.1249,
    51.50229,
    -2.80441,
    51.64992
   ]
  },
  {
   "en": "Powys",
   "cy": "",
   "geoType": "LAD",
   "geoCode": "W06000023",
   "bbox": [
    -3.92988,
    51.75381,
    -2.95093,
    52.90186
   ]
  },
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
`
}
