package handler

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/muhadif/e-wallet/core/entity"
	"github.com/muhadif/e-wallet/core/module"
	"github.com/muhadif/e-wallet/pkg"
	"io/ioutil"
	"net/http"
)

type WalletHandlerAPI interface {
	DepositWalletByWalletID(w http.ResponseWriter, r *http.Request)
	GetDepositByWalletID(w http.ResponseWriter, r *http.Request)
}

type walletHandler struct {
	walletModule module.WalletModule
}

func NewWalletHandlerAPI(walletModule module.WalletModule) WalletHandlerAPI {
	return &walletHandler{walletModule: walletModule}
}

func (h walletHandler) DepositWalletByWalletID(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var req *entity.DepositByWalletID
	if err := json.Unmarshal(reqBody, &req); err != nil {
		fmt.Println(err.Error())
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		pkg.ResponseError(w, http.StatusBadRequest, err)
		return
	}

	err := h.walletModule.DepositBalanceByWalletID(r.Context(), req)
	if err != nil {
		pkg.ResponseError(w, http.StatusInternalServerError, err)
		return
	}

	pkg.ResponseSuccess(w, http.StatusCreated, nil)
}

func (h walletHandler) GetDepositByWalletID(w http.ResponseWriter, r *http.Request) {
	walletID := mux.Vars(r)["walletID"]

	wallet, err := h.walletModule.GetBalanceByWalletID(r.Context(), walletID)
	if err != nil {
		pkg.ResponseError(w, http.StatusInternalServerError, err)
		return
	}

	pkg.ResponseSuccess(w, http.StatusOK, wallet)
}
