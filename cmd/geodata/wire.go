//+build wireinject

package main

import (
	"github.com/ONSdigital/dp-find-insights-poc-api/metadata"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/aws"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/database"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/geodata"
	"github.com/google/wire"
)

func InitService(maxmetrics geodata.MetricCount) (*Service, error) {
	wire.Build(
		aws.ProvideAWS,
		database.ProvidePassword,
		database.ProvideDSN,
		database.ProvideDatabase,
		geodata.ProvideGeodata,
		metadata.ProvideMetadata,
		wire.Struct(new(Service), "*"),
	)
	return &Service{}, nil
}
