// session obrero manage session related functions like session encryption and decryption
package main

import (
	"github.com/lock-free/gopcp"
	"github.com/lock-free/gopcp_stream"
	"github.com/lock-free/obrero"
	"github.com/lock-free/obrero/mids"
	"github.com/lock-free/obrero/napool"
	"github.com/lock-free/obrero/utils"
	"github.com/lock-free/session_obrero/session"
	"time"
)

const APP_CONFIG = "/data/app.json"

var pcpClient = gopcp.PcpClient{}

type AppConfig struct {
	SESSION_SECRECT_KEY string
	AUTH_WP_NAME        string
	AUTH_METHOD         string
}

func main() {
	var appConfig AppConfig
	err := utils.ReadJson(APP_CONFIG, &appConfig)
	if err != nil {
		panic(err)
	}

	var naPools napool.NAPools

	naPools = obrero.StartWorker(func(*gopcp_stream.StreamServer) *gopcp.Sandbox {
		return gopcp.GetSandbox(map[string]*gopcp.BoxFunc{
			"getServiceType": gopcp.ToSandboxFun(func(args []interface{}, attachment interface{}, pcpServer *gopcp.PcpServer) (interface{}, error) {
				return "session_obrero", nil
			}),

			// (encryptSession, sessionText)
			"encryptSession": gopcp.ToSandboxFun(mids.LogMid("encryptSession", func(args []interface{}, attachment interface{}, pcpServer *gopcp.PcpServer) (interface{}, error) {
				// parse args
				var (
					value string
				)
				err := utils.ParseArgs(args, []interface{}{&value}, "(encryptSession, sessionText)")
				if err != nil {
					return nil, err
				}

				return session.Encrypt([]byte(appConfig.SESSION_SECRECT_KEY), value) // encrypt value with session key
			})),

			// (getUidFromSessionText, text, timeout)
			"getUidFromSessionText": gopcp.ToSandboxFun(mids.LogMid("getUidFromSessionText", func(args []interface{}, attachment interface{}, pcpServer *gopcp.PcpServer) (interface{}, error) {
				// parse args
				var (
					text    string
					timeout int
				)
				err := utils.ParseArgs(args, []interface{}{&text, &timeout}, "(getUidFromSessionText, sessionText)")
				if err != nil {
					return nil, err
				}

				// decrypt session
				sessionTxt, err := session.Decrypt([]byte(appConfig.SESSION_SECRECT_KEY), text)

				if err != nil {
					return nil, err
				}

				// query user service to get uid
				return naPools.CallProxy(appConfig.AUTH_WP_NAME, pcpClient.Call(appConfig.AUTH_METHOD, sessionTxt), time.Duration(timeout)*time.Second)
			})),
		})
	}, obrero.WorkerStartConf{
		PoolSize:            2,
		Duration:            20 * time.Second,
		RetryDuration:       20 * time.Second,
		NAGetClientMaxRetry: 3,
	})

	utils.RunForever()
}
