package provision

import (
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/service/iot"
	"gopkg.in/yaml.v2"
)

// ThingConfig contains credentials for a thing.
type ThingConfig struct {
	// The ARN of the certificate.
	CertificateArn string

	// The ID of the certificate. AWS IoT issues a default subject name for the
	// certificate (e.g., AWS IoT Certificate).
	CertificateID string

	// The certificate data, in PEM format.
	CertificatePem string

	// The generated key pair.
	KeyPair *iot.KeyPair
}

// NewThingConfig create a new thing configuration using the response
func NewThingConfig(resp *iot.CreateKeysAndCertificateOutput) *ThingConfig {
	return &ThingConfig{
		CertificateArn: *resp.CertificateArn,
		CertificateID:  *resp.CertificateId,
		CertificatePem: *resp.CertificatePem,
		KeyPair:        resp.KeyPair,
	}
}

// Save the configuration to a YAML file with the name used as the filename
func (tc *ThingConfig) Save(name string) error {

	d, err := yaml.Marshal(&tc)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(fmt.Sprintf("%s.yaml", name), d, 0600)
}

// LoadConfig load the configuration from a file
func LoadConfig(name string) (*ThingConfig, error) {

	tc := &ThingConfig{}

	d, err := ioutil.ReadFile(fmt.Sprintf("%s.yaml", name))

	if err != nil {
		return tc, err
	}

	return tc, yaml.Unmarshal(d, &tc)
}
