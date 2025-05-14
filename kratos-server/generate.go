//go:build !wireinject

package dc1

//go:generate kratos proto client api/gateway/v1

//go:generate kratos proto client app/dc1server/internal/conf/conf.proto

//go:generate kratos proto client app/gateway/internal/conf/conf.proto

//go:generate wire .\app\dc1server\cmd

//go:generate wire .\app\gateway\cmd
