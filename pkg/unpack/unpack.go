package unpack

import (
	"crypto/x509"
	"fmt"

	"github.com/varshaprasad96/kubectl-rukpak/pkg/unpack/source"
	"github.com/varshaprasad96/kubectl-rukpak/pkg/util"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cluster"
)

const (
	coreServiceName       = "core"
	defaultBundleCacheDir = "/var/cache/bundles"
)

type unpackOptions struct {
	httpBindAddr         string
	httpExternalAddr     string
	bundleCAFile         string
	systemNamespace      string
	namespace            string
	unpackImage          string
	baseUploadManagerURL string
	storageDirectory     string
}

func (u *unpackOptions) Complete() {
	if u.httpBindAddr == "" {
		u.httpBindAddr = ":8080"
	}

	if u.httpExternalAddr == "" {
		u.httpExternalAddr = "http://localhost:8080"
	}

	if u.bundleCAFile == "" {
		u.bundleCAFile = "/etc/pki/tls/ca.crt"
	}

	if u.systemNamespace == "" {
		u.systemNamespace = "rukpak-system"
	}
	u.namespace = util.PodNamespace(u.systemNamespace)

	if u.unpackImage == "" {
		u.unpackImage = "quay.io/operator-framework/rukpak:latest"
	}

	u.baseUploadManagerURL = fmt.Sprintf("https://%s.%s.svc", coreServiceName, u.namespace)

	if u.storageDirectory == "" {
		u.storageDirectory = defaultBundleCacheDir
	}
}

func (u *unpackOptions) Run() error {
	cfg := ctrl.GetConfigOrDie()
	scheme := runtime.NewScheme()

	systemNsCluster, err := cluster.New(cfg, func(opts *cluster.Options) {
		opts.Scheme = scheme
		opts.Namespace = u.systemNamespace
	})
	if err != nil {
		return err
	}

	var rootCAs *x509.CertPool

	fmt.Println("rootCAs before", rootCAs, "bundleCAFile", u.bundleCAFile)

	if u.bundleCAFile != "" {
		var err error
		if rootCAs, err = util.LoadCertPool(u.bundleCAFile); err != nil {
			return fmt.Errorf("error loading the certificate %v", err)
		}
	}

	_, err = source.NewDefaultUnpacker(systemNsCluster, u.namespace, u.unpackImage, u.baseUploadManagerURL, rootCAs)
	if err != nil {
		return err
	}

	return nil
}
