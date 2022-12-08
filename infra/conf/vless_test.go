package conf_test

import (
	"testing"

	"github.com/dharak36/xray-core/common/net"
	"github.com/dharak36/xray-core/common/protocol"
	"github.com/dharak36/xray-core/common/serial"
	. "github.com/dharak36/xray-core/infra/conf"
	"github.com/dharak36/xray-core/proxy/vless"
	"github.com/dharak36/xray-core/proxy/vless/inbound"
	"github.com/dharak36/xray-core/proxy/vless/outbound"
)

func TestVLessOutbound(t *testing.T) {
	creator := func() Buildable {
		return new(VLessOutboundConfig)
	}

	runMultiTestCase(t, []TestCase{
		{
			Input: `{
				"vnext": [{
					"address": "example.com",
					"port": 443,
					"users": [
						{
							"id": "27848739-7e62-4138-9fd3-098a63964b6b",
							"flow": "xtls-rprx-direct-udp443",
							"encryption": "none",
							"level": 0
						}
					]
				}]
			}`,
			Parser: loadJSON(creator),
			Output: &outbound.Config{
				Vnext: []*protocol.ServerEndpoint{
					{
						Address: &net.IPOrDomain{
							Address: &net.IPOrDomain_Domain{
								Domain: "example.com",
							},
						},
						Port: 443,
						User: []*protocol.User{
							{
								Account: serial.ToTypedMessage(&vless.Account{
									Id:         "27848739-7e62-4138-9fd3-098a63964b6b",
									Flow:       "xtls-rprx-direct-udp443",
									Encryption: "none",
								}),
								Level: 0,
							},
						},
					},
				},
			},
		},
	})
}

func TestVLessInbound(t *testing.T) {
	creator := func() Buildable {
		return new(VLessInboundConfig)
	}

	runMultiTestCase(t, []TestCase{
		{
			Input: `{
				"clients": [
					{
						"id": "27848739-7e62-4138-9fd3-098a63964b6b",
						"flow": "xtls-rprx-direct",
						"level": 0,
						"email": "love@example.com"
					}
				],
				"decryption": "none",
				"fallbacks": [
					{
						"dest": 80
					},
					{
						"alpn": "h2",
						"dest": "@/dev/shm/domain.socket",
						"xver": 2
					},
					{
						"path": "/innerws",
						"dest": "serve-ws-none"
					}
				]
			}`,
			Parser: loadJSON(creator),
			Output: &inbound.Config{
				Clients: []*protocol.User{
					{
						Account: serial.ToTypedMessage(&vless.Account{
							Id:   "27848739-7e62-4138-9fd3-098a63964b6b",
							Flow: "xtls-rprx-direct",
						}),
						Level: 0,
						Email: "love@example.com",
					},
				},
				Decryption: "none",
				Fallbacks: []*inbound.Fallback{
					{
						Alpn: "",
						Path: "",
						Type: "tcp",
						Dest: "127.0.0.1:80",
						Xver: 0,
					},
					{
						Alpn: "h2",
						Path: "",
						Type: "unix",
						Dest: "@/dev/shm/domain.socket",
						Xver: 2,
					},
					{
						Alpn: "",
						Path: "/innerws",
						Type: "serve",
						Dest: "serve-ws-none",
						Xver: 0,
					},
				},
			},
		},
	})
}