package route

import (
	"fmt"
	"github.com/Dreamacro/clash/component/resolver"
	"github.com/Dreamacro/clash/config"
	"github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/hub/executor"
	P "github.com/Dreamacro/clash/listener"
	"github.com/Dreamacro/clash/log"
	"github.com/Dreamacro/clash/tunnel"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	fileMode os.FileMode = 0666
	dirMode  os.FileMode = 0755
)

func configRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", getConfigs)
	r.Put("/", updateConfigs)
	r.Patch("/", patchConfigs)
	return r
}

type configSchema struct {
	Port        *int               `json:"port"`
	SocksPort   *int               `json:"socks-port"`
	RedirPort   *int               `json:"redir-port"`
	TProxyPort  *int               `json:"tproxy-port"`
	MixedPort   *int               `json:"mixed-port"`
	AllowLan    *bool              `json:"allow-lan"`
	BindAddress *string            `json:"bind-address"`
	Mode        *tunnel.TunnelMode `json:"mode"`
	LogLevel    *log.LogLevel      `json:"log-level"`
	IPv6        *bool              `json:"ipv6"`
}

func getConfigs(w http.ResponseWriter, r *http.Request) {
	general := executor.GetGeneral()
	render.JSON(w, r, general)
}

func pointerOrDefault(p *int, def int) int {
	if p != nil {
		return *p
	}
	return def
}

func patchConfigs(w http.ResponseWriter, r *http.Request) {
	general := &configSchema{}
	if err := render.DecodeJSON(r.Body, general); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		return
	}

	if general.AllowLan != nil {
		P.SetAllowLan(*general.AllowLan)
	}

	if general.BindAddress != nil {
		P.SetBindAddress(*general.BindAddress)
	}

	ports := P.GetPorts()

	tcpIn := tunnel.TCPIn()
	udpIn := tunnel.UDPIn()

	P.ReCreateHTTP(pointerOrDefault(general.Port, ports.Port), tcpIn)
	P.ReCreateSocks(pointerOrDefault(general.SocksPort, ports.SocksPort), tcpIn, udpIn)
	P.ReCreateRedir(pointerOrDefault(general.RedirPort, ports.RedirPort), tcpIn, udpIn)
	P.ReCreateTProxy(pointerOrDefault(general.TProxyPort, ports.TProxyPort), tcpIn, udpIn)
	P.ReCreateMixed(pointerOrDefault(general.MixedPort, ports.MixedPort), tcpIn, udpIn)

	if general.Mode != nil {
		tunnel.SetMode(*general.Mode)
	}

	if general.LogLevel != nil {
		log.SetLevel(*general.LogLevel)
	}

	if general.IPv6 != nil {
		resolver.DisableIPv6 = !*general.IPv6
	}

	render.NoContent(w, r)
}

type updateConfigRequest struct {
	Path    string `json:"path"`
	Payload string `json:"payload"`
}

func updateConfigs(w http.ResponseWriter, r *http.Request) {
	req := updateConfigRequest{}
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		return
	}

	force := r.URL.Query().Get("force") == "true"
	var cfg *config.Config
	var err error
	backupCfgBuff := config.GlobalHusk.CurrentCfgBuff
	if req.Payload != "" {
		cfg, err = executor.ParseWithBytes([]byte(req.Payload))
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, newError(err.Error()))
			return
		}
		msg, err := saveConfig(config.GlobalHusk.AllowOverWrite, config.GlobalHusk.IsBackupCfg, []byte(req.Payload), backupCfgBuff)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, newError(err.Error()))
			return
		}
		render.Status(r, http.StatusOK)
		render.JSON(w, r, render.M{
			"message": msg,
		})
		executor.ApplyConfig(cfg, force)

		return
	} else {
		if req.Path == "" {
			req.Path = constant.Path.Config()
		}
		if !filepath.IsAbs(req.Path) {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, newError("path is not a absolute path"))
			return
		}

		cfg, err = executor.ParseWithPath(req.Path)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, newError(err.Error()))
			return
		}
	}
	executor.ApplyConfig(cfg, force)
	render.NoContent(w, r)
}

func saveConfig(allowOverwrite, backupCfg bool, updatedBuff, backupBuff []byte) (string, error) {
	var err error
	overWriteDone := false
	backupDone := false
	configPath := constant.Path.Config()
	backupPath := filepath.Join(constant.Path.HomeDir(), "cfgBackup", fmt.Sprintf("config-%d.yaml", time.Now().Unix()))
	if allowOverwrite {
		err = safeWrite(configPath, updatedBuff)
		if err != nil {
			return "", fmt.Errorf("overwrite config file with %s\n", err.Error())
		}
		overWriteDone = true
	} else {
		return "", fmt.Errorf("no allow to modify config file\n")
	}

	if backupCfg {
		err = safeWrite(backupPath, backupBuff)
		if err != nil {
			return "", fmt.Errorf("fail to backup config file, error: %s\n", err.Error())
		}
		backupDone = true
	}

	if overWriteDone {
		if backupDone {
			return fmt.Sprintf("overwrite config file: %s and backup config to %s\n", configPath, backupPath), nil
		} else {
			return fmt.Sprintf("overwrite config file: %s, no backup config file\n", configPath), nil
		}
	}
	return "", fmt.Errorf("no config file was ovewrite")
}

func safeWrite(path string, buf []byte) error {
	dir := filepath.Dir(path)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, dirMode); err != nil {
			return err
		}
	}
	return ioutil.WriteFile(path, buf, fileMode)
}
