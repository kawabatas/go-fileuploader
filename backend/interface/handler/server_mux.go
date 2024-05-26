package handler

import (
	"net/http"

	"github.com/kawabatas/go-fileuploader/application"
	"github.com/kawabatas/go-fileuploader/infrastructure/localfs"
	"github.com/kawabatas/go-fileuploader/infrastructure/mysqldb"
)

type ServerMux struct {
	router *http.ServeMux

	fileUseCase *application.FileUseCase
}

func NewServeMux(
	stage,
	bucketBaseURL, bucketBaseDir, publicDir,
	dbUser, dbPass, dbHost, dbName string,
) (*ServerMux, error) {
	dbClient, err := mysqldb.Initialize(dbUser, dbPass, dbHost, dbName)
	if err != nil {
		return nil, err
	}
	// 接続確認
	if err := dbClient.Ping(); err != nil {
		return nil, err
	}

	fileMetadataRepo := mysqldb.NewFileMetadataRepository(dbClient)
	userRepo := mysqldb.NewUserRepository(dbClient)

	fileBucketRepo := localfs.NewFileBucketRepository(bucketBaseDir)

	fileUseCase := application.NewFileUseCase(fileMetadataRepo, fileBucketRepo, userRepo, bucketBaseURL)

	router := http.NewServeMux()
	mux := &ServerMux{
		router:      router,
		fileUseCase: fileUseCase,
	}
	mux.routingPatterns(publicDir)
	return mux, nil
}

func (mux *ServerMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux.router.ServeHTTP(w, r)
}
