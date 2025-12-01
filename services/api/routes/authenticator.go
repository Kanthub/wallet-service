package routes

import (
	"encoding/json"
	"net/http"

	"github.com/ethereum/go-ethereum/log"

	"github.com/roothash-pay/wallet-services/services/api/models"
)

// GenerateTOTPHandler godoc
// @Summary 生成TOTP密钥
// @Description 生成60秒周期的TOTP密钥
// @Tags 身份认证
// @Accept json
// @Produce json
// @Param request body models.GenerateTOTPRequest true "生成TOTP请求"
// @Success 200 {object} models.GenerateTOTPResponse "生成成功"
// @Failure 400 {object} string "请求参数错误"
// @Failure 500 {object} string "服务器内部错误"
// @Router /authenticator/generate [post]
func (rs *Routes) GenerateTOTPHandler(w http.ResponseWriter, r *http.Request) {
	var req models.GenerateTOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	resp, err := rs.svc.GenerateTOTP(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = jsonResponse(w, resp, http.StatusOK)
	if err != nil {
		log.Error("Error writing response", "err", err.Error())
	}
}

// VerifyTOTPHandler godoc
// @Summary 验证TOTP代码
// @Description 验证60秒周期的TOTP验证码
// @Tags 身份认证
// @Accept json
// @Produce json
// @Param request body models.VerifyTOTPRequest true "验证TOTP请求"
// @Success 200 {object} models.VerifyTOTPResponse "验证成功"
// @Failure 400 {object} string "请求参数错误"
// @Failure 500 {object} string "服务器内部错误"
// @Router /authenticator/verify [post]
func (rs *Routes) VerifyTOTPHandler(w http.ResponseWriter, r *http.Request) {
	var req models.VerifyTOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	resp, err := rs.svc.VerifyTOTP(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = jsonResponse(w, resp, http.StatusOK)
	if err != nil {
		log.Error("Error writing response", "err", err.Error())
	}
}
