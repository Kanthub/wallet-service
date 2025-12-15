package routes

import (
	"net/http"

	"github.com/ethereum/go-ethereum/log"
	"github.com/go-chi/chi/v5"
)

func (rs *Routes) WalletBalanceApi() {
	r := rs.router
	r.Route("/api/v1/balance", func(r chi.Router) {

		r.Get("/wallet", rs.getWalletBalances)
		r.Get("/wallet-token", rs.getWalletBalanceByTokenChain)
		r.Get("/summary", rs.getWalletBalanceSummary)
	})
}

func (rs *Routes) getWalletBalances(w http.ResponseWriter, r *http.Request) {
	walletUUID := r.URL.Query().Get("wallet_uuid")
	if walletUUID == "" {
		http.Error(w, "wallet_uuid required", http.StatusBadRequest)
		return
	}

	list, err := rs.svc.WalletBalanceService.GetWalletBalances(
		r.Context(), walletUUID,
	)
	if err != nil {
		log.Error("get wallet balances error", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, list, http.StatusOK)
}

func (rs *Routes) getWalletBalanceByTokenChain(w http.ResponseWriter, r *http.Request) {
	walletUUID := r.URL.Query().Get("wallet_uuid")
	tokenID := r.URL.Query().Get("token_id")
	chainID := r.URL.Query().Get("chain_id")

	if walletUUID == "" || tokenID == "" || chainID == "" {
		http.Error(w, "wallet_uuid, token_id and chain_id required", http.StatusBadRequest)
		return
	}

	item, err := rs.svc.WalletBalanceService.GetWalletBalanceByTokenChain(
		r.Context(), walletUUID, tokenID, chainID,
	)
	if err != nil {
		log.Error("get wallet token balance error", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, item, http.StatusOK)
}

func (rs *Routes) getWalletBalanceSummary(w http.ResponseWriter, r *http.Request) {
	walletUUID := r.URL.Query().Get("wallet_uuid")
	if walletUUID == "" {
		http.Error(w, "wallet_uuid required", http.StatusBadRequest)
		return
	}

	summary, err := rs.svc.WalletBalanceService.GetWalletBalanceSummary(
		r.Context(), walletUUID,
	)
	if err != nil {
		log.Error("get wallet balance summary error", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, summary, http.StatusOK)
}
