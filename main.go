// session obrero manage session related functions like session encryption and decryption
package main

import (
	"github.com/lock-free/gopcp"
	"github.com/lock-free/gopcp_stream"
	"github.com/lock-free/obrero/mids"
	"github.com/lock-free/obrero/napool"
	"github.com/lock-free/obrero/stdserv"
	"github.com/lock-free/obrero/utils"
	"github.com/lock-free/session_obrero/session"
	"time"
)

var pcpClient = gopcp.PcpClient{}

type AppConfig struct {
	SESSION_SECRECT_KEY string
	AUTH_WP_NAME        string
	AUTH_METHOD         string
}

func main() {
	var appConfig AppConfig
	stdserv.StartStdWorker(&appConfig, func(naPools *napool.NAPools, workerState *stdserv.WorkerState, s *gopcp_stream.StreamServer) map[string]*gopcp.BoxFunc {
		return map[string]*gopcp.BoxFunc{
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
				err := utils.ParseArgs(args, []interface{}{&text, &timeout}, "(getUidFromSessionText, sessionText, timeout)")
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
		}
	}, stdserv.StdWorkerConfig{
		ServiceName: "session_obrero",
	})
}
