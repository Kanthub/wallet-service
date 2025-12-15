package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/log"
	"github.com/go-chi/chi/v5"
)

func (rs *Routes) SysLogApi() {
	r := rs.router
	r.Route("/api/v1/syslog", func(r chi.Router) {
		r.Get("/list", rs.listSysLogs)
		r.Get("/info", rs.getSysLog)
	})
}

func (rs *Routes) listSysLogs(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	filters := map[string]interface{}{
		"action":       r.URL.Query().Get("action"),
		"remark":       r.URL.Query().Get("remark"),
		"admin":        r.URL.Query().Get("admin"),
		"cate":         r.URL.Query().Get("cate"),
		"status":       r.URL.Query().Get("status"),
		"asset":        r.URL.Query().Get("asset"),
		"order_number": r.URL.Query().Get("order_number"),
	}

	list, total, err := rs.svc.SysLogService.ListLogs(
		r.Context(),
		page,
		pageSize,
		filters,
	)
	if err != nil {
		log.Error("list syslog failed", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"list":      list,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (rs *Routes) getSysLog(w http.ResponseWriter, r *http.Request) {
	guid := r.URL.Query().Get("guid")
	if guid == "" {
		http.Error(w, "guid is required", http.StatusBadRequest)
		return
	}

	item, err := rs.svc.SysLogService.GetLog(r.Context(), guid)
	if err != nil {
		log.Error("get syslog failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}
