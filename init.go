package main

import (
	"crypto/rand"
	"fmt"
	"os"
	"strconv"
	"strings"

	pb "github.com/PretendoNetwork/grpc/go/account"
	"github.com/PretendoNetwork/plogger-go"
	"github.com/joho/godotenv"
	"github.com/silver-volt4/swapdoodle/database"
	"github.com/silver-volt4/swapdoodle/globals"

	"github.com/PretendoNetwork/nex-go/v2"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func verifyPort(port string) (int, error) {
	port = strings.TrimSpace(port)
	if port == "" {
		return 0, fmt.Errorf("Environment variable not set.")
	}
	result, err := strconv.Atoi(port)
	if err != nil {
		return 0, fmt.Errorf("Invalid port. Expected 0-65535, got %s", port)
	} else if result < 0 || result > 65535 {
		return 0, fmt.Errorf("Invalid port. Expected 0-65535, got %s", port)
	}
	return result, nil
}

func init() {
	globals.Logger = plogger.NewLogger()

	var err error

	err = godotenv.Load()
	if err != nil {
		globals.Logger.Warning("Error loading .env file")
	}

	s3Endpoint := os.Getenv("PN_SD_CONFIG_S3_ENDPOINT")
	s3AccessKey := os.Getenv("PN_SD_CONFIG_S3_ACCESS_KEY")
	s3AccessSecret := os.Getenv("PN_SD_CONFIG_S3_ACCESS_SECRET")
	s3SecureEnv := os.Getenv("PN_SD_CONFIG_S3_SECURE")
	s3Bucket := os.Getenv("PN_SD_CONFIG_S3_BUCKET")

	postgresURI := os.Getenv("PN_SD_POSTGRES_URI")
	hppServerPort := os.Getenv("PN_SD_HPP_SERVER_PORT")
	grpcServerPort := os.Getenv("PN_SD_GRPC_SERVER_PORT")
	grpcApiKey := os.Getenv("PN_SD_CONFIG_GRPC_API_KEY")
	accountGRPCHost := os.Getenv("PN_SD_ACCOUNT_GRPC_HOST")
	accountGRPCPort := os.Getenv("PN_SD_ACCOUNT_GRPC_PORT")
	accountGRPCAPIKey := os.Getenv("PN_SD_ACCOUNT_GRPC_API_KEY")

	if strings.TrimSpace(postgresURI) == "" {
		globals.Logger.Error("PN_SD_POSTGRES_URI environment variable not set")
		os.Exit(0)
	}

	kerberosPassword := make([]byte, 0x10)
	_, err = rand.Read(kerberosPassword)
	if err != nil {
		globals.Logger.Error("Error generating Kerberos password")
		os.Exit(0)
	}

	globals.KerberosPassword = string(kerberosPassword)

	globals.HppServerAccount = nex.NewAccount(2, "Quazal Rendez-Vous", globals.KerberosPassword)

	globals.HPPServerPort, err = verifyPort(hppServerPort)
	if err != nil {
		globals.Logger.Errorf("Error in environment variable PN_SD_HPP_SERVER_PORT: %s", err.Error())
		os.Exit(0)
	}

	globals.GRPCServerPort, err = verifyPort(grpcServerPort)
	if err != nil {
		globals.Logger.Errorf("Error in environment variable PN_SD_GRPC_SERVER_PORT: %s", err.Error())
		os.Exit(0)
	}

	globals.GRPCApiKey = strings.TrimSpace(grpcApiKey)
	if globals.GRPCApiKey == "" {
		globals.Logger.Warning("PN_SD_CONFIG_GRPC_API_KEY is not set. Your gRPC server will be open.")
	}

	if strings.TrimSpace(accountGRPCAPIKey) == "" {
		globals.Logger.Warning("Insecure gRPC server detected. PN_SD_ACCOUNT_GRPC_API_KEY environment variable not set")
	}

	globals.GRPCAccountClientConnection, err = grpc.NewClient(fmt.Sprintf("%s:%s", accountGRPCHost, accountGRPCPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		globals.Logger.Criticalf("Failed to connect to account gRPC server: %v", err)
		os.Exit(0)
	}

	globals.GRPCAccountClient = pb.NewAccountClient(globals.GRPCAccountClientConnection)
	globals.GRPCAccountCommonMetadata = metadata.Pairs(
		"X-API-Key", accountGRPCAPIKey,
	)

	staticCredentials := credentials.NewStaticV4(s3AccessKey, s3AccessSecret, "")

	s3Secure, err := strconv.ParseBool(s3SecureEnv)
	if err != nil {
		globals.Logger.Warningf("PN_SD_CONFIG_S3_SECURE environment variable not set. Using default value: %t", true)
		s3Secure = true
	}

	minIOClient, err := minio.New(s3Endpoint, &minio.Options{
		Creds:  staticCredentials,
		Secure: s3Secure,
	})
	if err != nil {
		globals.Logger.Criticalf("Failed to connect to initialize minIOClient: %v", err)
		os.Exit(0)
	}

	globals.MinIOClient = minIOClient

	globals.S3BucketName = s3Bucket

	database.ConnectPostgres(postgresURI)
}
