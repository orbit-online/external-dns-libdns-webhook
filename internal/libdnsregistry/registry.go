// This file is auto generated, DO NOT EDIT.
package libdnsregistry

import (
	libdnsautodns "github.com/libdns/autodns"
	libdnsbunny "github.com/libdns/bunny"
	libdnsdesec "github.com/libdns/desec"
	libdnsdirectadmin "github.com/libdns/directadmin"
	libdnsdnsimple "github.com/libdns/dnsimple"
	libdnsduckdns "github.com/libdns/duckdns"
	libdnsdynu "github.com/libdns/dynu"
	libdnsdynv6 "github.com/libdns/dynv6"
	libdnseasydns "github.com/libdns/easydns"
	libdnshe "github.com/libdns/he"
	libdnsinfomaniak "github.com/libdns/infomaniak"
	libdnsinwx "github.com/libdns/inwx"
	libdnsloopia "github.com/libdns/loopia"
	libdnsluadns "github.com/libdns/luadns"
	libdnsmailinabox "github.com/libdns/mailinabox"
	libdnsmetaname "github.com/libdns/metaname"
	libdnsmythicbeasts "github.com/libdns/mythicbeasts"
	libdnsnamesilo "github.com/libdns/namesilo"
	libdnsnetlify "github.com/libdns/netlify"
	libdnsnfsn "github.com/libdns/nfsn"
	libdnsnjalla "github.com/libdns/njalla"
	libdnsporkbun "github.com/libdns/porkbun"
	libdnsrfc2136 "github.com/libdns/rfc2136"
	libdnstransip "github.com/libdns/transip"
)

var registry = RegistryStore{
	"rfc2136": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsrfc2136.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://tools.ietf.org/html/rfc2136",
				Description:   "RFC2136 is a standard for dynamic DNS updates.",
				Configuration: configurationDetails[libdnsrfc2136.Provider](),
			}
		},
	},
	"inwx": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsinwx.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://www.inwx.com/",
				Description:   "INWX is a domain registrar and DNS hosting provider.",
				Configuration: configurationDetails[libdnsinwx.Provider](),
			}
		},
	},
	"autodns": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsautodns.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://www.internetx.com/",
				Description:   "AutoDNS is a DNS management service by InterNetX.",
				Configuration: configurationDetails[libdnsautodns.Provider](),
			}
		},
	},
	"he": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnshe.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://dns.he.net/",
				Description:   "Hurricane Electric provides DNS hosting services.",
				Configuration: configurationDetails[libdnshe.Provider](),
			}
		},
	},
	"porkbun": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsporkbun.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://porkbun.com/",
				Description:   "Porkbun is a domain registrar and DNS hosting provider.",
				Configuration: configurationDetails[libdnsporkbun.Provider](),
			}
		},
	},
	"luadns": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsluadns.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://www.luadns.com/",
				Description:   "LuaDNS provides DNS hosting services.",
				Configuration: configurationDetails[libdnsluadns.Provider](),
			}
		},
	},
	"dynu": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsdynu.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://www.dynu.com/",
				Description:   "Dynu provides dynamic DNS services.",
				Configuration: configurationDetails[libdnsdynu.Provider](),
			}
		},
	},
	"easydns": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnseasydns.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://www.easydns.com/",
				Description:   "EasyDNS provides DNS hosting services.",
				Configuration: configurationDetails[libdnseasydns.Provider](),
			}
		},
	},
	"transip": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnstransip.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://www.transip.eu/",
				Description:   "TransIP provides web hosting and domain services.",
				Configuration: configurationDetails[libdnstransip.Provider](),
			}
		},
	},
	"directadmin": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsdirectadmin.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://www.directadmin.com/",
				Description:   "DirectAdmin is a web hosting control panel.",
				Configuration: configurationDetails[libdnsdirectadmin.Provider](),
			}
		},
	},
	"infomaniak": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsinfomaniak.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://www.infomaniak.com/",
				Description:   "Infomaniak provides web hosting and domain services.",
				Configuration: configurationDetails[libdnsinfomaniak.Provider](),
			}
		},
	},
	"desec": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsdesec.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://desec.io/",
				Description:   "deSEC is a free DNS hosting service.",
				Configuration: configurationDetails[libdnsdesec.Provider](),
			}
		},
	},
	"dnsimple": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsdnsimple.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://dnsimple.com/",
				Description:   "DNSimple provides DNS hosting services.",
				Configuration: configurationDetails[libdnsdnsimple.Provider](),
			}
		},
	},
	"netlify": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsnetlify.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://www.netlify.com/",
				Description:   "Netlify provides web hosting and domain services.",
				Configuration: configurationDetails[libdnsnetlify.Provider](),
			}
		},
	},
	"nfsn": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsnfsn.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://www.nearlyfreespeech.net/",
				Description:   "NearlyFreeSpeech.NET provides web hosting and domain services.",
				Configuration: configurationDetails[libdnsnfsn.Provider](),
			}
		},
	},
	"namesilo": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsnamesilo.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://www.namesilo.com/",
				Description:   "NameSilo is a domain registrar and DNS hosting provider.",
				Configuration: configurationDetails[libdnsnamesilo.Provider](),
			}
		},
	},
	"bunny": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsbunny.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://bunny.net/",
				Description:   "BunnyCDN provides content delivery network services.",
				Configuration: configurationDetails[libdnsbunny.Provider](),
			}
		},
	},
	"loopia": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsloopia.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://www.loopia.se/",
				Description:   "Loopia provides web hosting and domain services.",
				Configuration: configurationDetails[libdnsloopia.Provider](),
			}
		},
	},
	"mailinabox": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsmailinabox.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://mailinabox.email/",
				Description:   "Mail-in-a-Box is an easy-to-deploy mail server.",
				Configuration: configurationDetails[libdnsmailinabox.Provider](),
			}
		},
	},
	"duckdns": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsduckdns.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://www.duckdns.org/",
				Description:   "DuckDNS provides free dynamic DNS services.",
				Configuration: configurationDetails[libdnsduckdns.Provider](),
			}
		},
	},
	"njalla": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsnjalla.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://njal.la/",
				Description:   "Njalla provides privacy-focused domain registration and DNS services.",
				Configuration: configurationDetails[libdnsnjalla.Provider](),
			}
		},
	},
	"mythicbeasts": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsmythicbeasts.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://www.mythic-beasts.com/",
				Description:   "Mythic Beasts provides web hosting and domain services.",
				Configuration: configurationDetails[libdnsmythicbeasts.Provider](),
			}
		},
	},
	"dynv6": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsdynv6.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://dynv6.com/",
				Description:   "Dynv6 provides free dynamic DNS services.",
				Configuration: configurationDetails[libdnsdynv6.Provider](),
			}
		},
	},
	"metaname": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsmetaname.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://metaname.net/",
				Description:   "Metaname provides domain registration and DNS services.",
				Configuration: configurationDetails[libdnsmetaname.Provider](),
			}
		},
	},
}
