package models

import (
	"time"

	"github.com/chainreactors/malice-network/proto/client/clientpb"
	"github.com/chainreactors/malice-network/proto/client/commonpb"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// ImplantBuild - Represents an implant
type ImplantBuild struct {
	ID        uuid.UUID `gorm:"primaryKey;->;<-:create;type:uuid;"`
	CreatedAt time.Time `gorm:"->;<-:create;"`

	Name string `gorm:"unique;"`

	// Checksums stores of the implant binary
	MD5    string
	SHA1   string
	SHA256 string

	// Burned indicates whether the implant
	// has been seen on threat intel platforms
	Burned bool

	ImplantConfig ImplantConfig
}

// BeforeCreate - GORM hook
func (ib *ImplantBuild) BeforeCreate(tx *gorm.DB) (err error) {
	ib.ID, err = uuid.NewV4()
	if err != nil {
		return err
	}
	ib.CreatedAt = time.Now()
	return nil
}

// ImplantConfig - An implant build configuration
type ImplantConfig struct {
	ID               uuid.UUID `gorm:"primaryKey;->;<-:create;type:uuid;"`
	ImplantBuildID   uuid.UUID
	ImplantProfileID uuid.UUID

	CreatedAt time.Time `gorm:"->;<-:create;"`

	// Go
	GOOS   string
	GOARCH string

	TemplateName string

	IsBeacon       bool
	BeaconInterval int64
	BeaconJitter   int64

	// ECC
	PeerPublicKey           string
	PeerPublicKeyDigest     string
	PeerPrivateKey          string
	PeerPublicKeySignature  string
	AgeServerPublicKey      string
	MinisignServerPublicKey string

	// MTLS
	MtlsCACert string
	MtlsCert   string
	MtlsKey    string

	Debug               bool
	DebugFile           string
	Evasion             bool
	ObfuscateSymbols    bool
	ReconnectInterval   int64
	PollTimeout         int64
	MaxConnectionErrors uint32
	ConnectionStrategy  string
	SGNEnabled          bool

	// WireGuard
	WGImplantPrivKey  string
	WGServerPubKey    string
	WGPeerTunIP       string
	WGKeyExchangePort uint32
	WGTcpCommsPort    uint32

	C2 []ImplantC2

	MTLSc2Enabled bool
	WGc2Enabled   bool
	HTTPc2Enabled bool
	DNSc2Enabled  bool

	CanaryDomains     []CanaryDomain
	NamePipec2Enabled bool
	TCPPivotc2Enabled bool

	// Limits
	LimitDomainJoined bool
	LimitHostname     string
	LimitUsername     string
	LimitDatetime     string
	LimitFileExists   string
	LimitLocale       string

	// Output Format
	Format clientpb.OutputFormat

	// For 	IsSharedLib bool
	IsSharedLib bool
	IsService   bool
	IsShellcode bool

	RunAtLoad bool

	FileName string

	NetGoEnabled           bool
	TrafficEncodersEnabled bool
	Assets                 []EncoderAsset
}

// BeforeCreate - GORM hook
func (ic *ImplantConfig) BeforeCreate(tx *gorm.DB) (err error) {
	ic.ID, err = uuid.NewV4()
	if err != nil {
		return err
	}
	ic.CreatedAt = time.Now()
	return nil
}

// ToProtobuf - Convert ImplantConfig to protobuf equiv
func (ic *ImplantConfig) ToProtobuf() *clientpb.ImplantConfig {
	config := &clientpb.ImplantConfig{
		ID: ic.ID.String(),

		IsBeacon:       ic.IsBeacon,
		BeaconInterval: ic.BeaconInterval,
		BeaconJitter:   ic.BeaconJitter,

		GOOS:               ic.GOOS,
		GOARCH:             ic.GOARCH,
		AgeServerPublicKey: ic.AgeServerPublicKey,
		PeerPublicKey:      ic.PeerPublicKey,
		PeerPrivateKey:     ic.PeerPrivateKey,
		MtlsCACert:         ic.MtlsCACert,
		MtlsCert:           ic.MtlsCert,
		MtlsKey:            ic.MtlsKey,

		Debug:            ic.Debug,
		DebugFile:        ic.DebugFile,
		Evasion:          ic.Evasion,
		ObfuscateSymbols: ic.ObfuscateSymbols,
		TemplateName:     ic.TemplateName,
		SGNEnabled:       ic.SGNEnabled,

		ReconnectInterval:   ic.ReconnectInterval,
		MaxConnectionErrors: ic.MaxConnectionErrors,
		PollTimeout:         ic.PollTimeout,
		ConnectionStrategy:  ic.ConnectionStrategy,

		LimitDatetime:     ic.LimitDatetime,
		LimitDomainJoined: ic.LimitDomainJoined,
		LimitHostname:     ic.LimitHostname,
		LimitUsername:     ic.LimitUsername,
		LimitFileExists:   ic.LimitFileExists,
		LimitLocale:       ic.LimitLocale,

		IsSharedLib:       ic.IsSharedLib,
		IsService:         ic.IsService,
		IsShellcode:       ic.IsShellcode,
		Format:            ic.Format,
		WGImplantPrivKey:  ic.WGImplantPrivKey,
		WGServerPubKey:    ic.WGServerPubKey,
		WGPeerTunIP:       ic.WGPeerTunIP,
		WGKeyExchangePort: ic.WGKeyExchangePort,
		WGTcpCommsPort:    ic.WGTcpCommsPort,

		FileName:               ic.FileName,
		TrafficEncodersEnabled: ic.TrafficEncodersEnabled,
		NetGoEnabled:           ic.NetGoEnabled,
	}
	// Copy Canary Domains
	config.CanaryDomains = []string{}
	for _, canaryDomain := range ic.CanaryDomains {
		config.CanaryDomains = append(config.CanaryDomains, canaryDomain.Domain)
	}

	// Copy Assets
	config.Assets = []*commonpb.File{}
	for _, asset := range ic.Assets {
		config.Assets = append(config.Assets, asset.ToProtobuf())
	}

	// Copy C2
	config.C2 = []*clientpb.ImplantC2{}
	for _, c2 := range ic.C2 {
		config.C2 = append(config.C2, c2.ToProtobuf())
	}
	return config
}

// CanaryDomainsList - Get string slice of canary domains
func (ic *ImplantConfig) CanaryDomainsList() []string {
	domains := []string{}
	for _, canaryDomain := range ic.CanaryDomains {
		domains = append(domains, canaryDomain.Domain)
	}
	return domains
}

// CanaryDomain - Canary domain, belongs to ImplantConfig
type CanaryDomain struct {
	ID              uuid.UUID `gorm:"primaryKey;->;<-:create;type:uuid;"`
	ImplantConfigID uuid.UUID
	CreatedAt       time.Time `gorm:"->;<-:create;"`

	Domain string
}

// BeforeCreate - GORM hook
func (c *CanaryDomain) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID, err = uuid.NewV4()
	if err != nil {
		return err
	}
	c.CreatedAt = time.Now()
	return nil
}

// ImplantC2 - C2 struct
type ImplantC2 struct {
	// gorm.Model

	ID              uuid.UUID `gorm:"primaryKey;->;<-:create;type:uuid;"`
	ImplantConfigID uuid.UUID
	CreatedAt       time.Time `gorm:"->;<-:create;"`

	Priority uint32
	URL      string
	Options  string
}

// BeforeCreate - GORM hook
func (c2 *ImplantC2) BeforeCreate(tx *gorm.DB) (err error) {
	c2.ID, err = uuid.NewV4()
	if err != nil {
		return err
	}
	c2.CreatedAt = time.Now()
	return nil
}

// ToProtobuf - Convert to protobuf version
func (c2 *ImplantC2) ToProtobuf() *clientpb.ImplantC2 {
	return &clientpb.ImplantC2{
		Priority: c2.Priority,
		URL:      c2.URL,
		Options:  c2.Options,
	}
}

func (c2 *ImplantC2) String() string {
	return c2.URL
}

// ImplantProfile - An implant build configuration
type ImplantProfile struct {
	// gorm.Model

	ID        uuid.UUID `gorm:"primaryKey;->;<-:create;type:uuid;"`
	CreatedAt time.Time `gorm:"->;<-:create;"`

	Name          string `gorm:"unique;"`
	ImplantConfig *ImplantConfig
}

// BeforeCreate - GORM hook
func (ip *ImplantProfile) BeforeCreate(tx *gorm.DB) (err error) {
	ip.ID, err = uuid.NewV4()
	if err != nil {
		return err
	}
	ip.CreatedAt = time.Now()
	return nil
}

// EncoderAsset - Tracks which assets were embedded into the implant
// but we currently don't keep a copy of the actual data
type EncoderAsset struct {
	ID              uuid.UUID `gorm:"primaryKey;->;<-:create;type:uuid;"`
	ImplantConfigID uuid.UUID

	Name string
}

func (t *EncoderAsset) ToProtobuf() *commonpb.File {
	return &commonpb.File{Name: t.Name}
}
