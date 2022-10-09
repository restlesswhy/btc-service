package ledger

import (
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"path"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"github.com/restlesswhy/btc-service/config"
	"github.com/restlesswhy/btc-service/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type instance struct {
	id       int
	grpcConn *grpc.ClientConn
	gateway  *client.Gateway
	network  *client.Network
	contract *client.Contract
	log      logger.Logger

	done chan int
}

func newInstance(cfg *config.Config, id int, done chan int, log logger.Logger) (*instance, error) {
	clientConnection := newGrpcConnection(cfg)

	clientid := newIdentity(cfg)
	sign := newSign(cfg)

	// Create a Gateway connection for a specific client identity
	gateway, err := client.Connect(
		clientid,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		// Default timeouts for different gRPC calls
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		return nil, err
	}

	network := gateway.GetNetwork(cfg.Ledger.ChannelName)
	contract := network.GetContract(cfg.Ledger.ChaincodeName)

	i := &instance{
		id:       id,
		grpcConn: clientConnection,
		gateway:  gateway,
		network:  network,
		contract: contract,
		log:      log,
		done:     done,
	}

	return i, nil
}

func newGrpcConnection(cfg *config.Config) *grpc.ClientConn {
	certificate, err := loadCertificate(cfg.Ledger.TlsCertPath)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, cfg.Ledger.GatewayPeer)

	connection, err := grpc.Dial(cfg.Ledger.PeerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		panic(fmt.Errorf("failed to create gRPC connection: %w", err))
	}

	return connection
}

// newIdentity creates a client identity for this Gateway connection using an X.509 certificate.
func newIdentity(cfg *config.Config) *identity.X509Identity {
	certificate, err := loadCertificate(cfg.Ledger.CertPath)
	if err != nil {
		panic(err)
	}

	id, err := identity.NewX509Identity(cfg.Ledger.MspID, certificate)
	if err != nil {
		panic(err)
	}

	return id
}

func loadCertificate(filename string) (*x509.Certificate, error) {
	certificatePEM, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}
	return identity.CertificateFromPEM(certificatePEM)
}

// newSign creates a function that generates a digital signature from a message digest using a private key.
func newSign(cfg *config.Config) identity.Sign {
	files, err := ioutil.ReadDir(cfg.Ledger.KeyPath)
	if err != nil {
		panic(fmt.Errorf("failed to read private key directory: %w", err))
	}
	privateKeyPEM, err := ioutil.ReadFile(path.Join(cfg.Ledger.KeyPath, files[0].Name()))

	if err != nil {
		panic(fmt.Errorf("failed to read private key file: %w", err))
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		panic(err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(err)
	}

	return sign
}

func (i *instance) Close() error {
	i.grpcConn.Close()
	i.gateway.Close()
	i.log.Info("connections closed, send id")
	i.done <- i.id

	return nil
}
