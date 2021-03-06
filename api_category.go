package main

import (
	"net/http"
	"strings"

	b64 "encoding/base64"

	"github.com/romana/rlog"
	"go.mongodb.org/mongo-driver/bson"
)

func getProductSkusInCategory(w http.ResponseWriter, r *http.Request) {

	rlog.Debug("getProductsInCategory() handle function invoked ...")

	if !pre(w, r) {
		return
	}

	csx := getAccessToken(r)
	ctcol := csx + CategoryTreeExtension

	var ctrq CATEGORYREQUEST

	if !mapInput(w, r, &ctrq) {
		return
	}

	path := cleanCategoryPath(ctrq.Path)

	if !pathExists(w, r, path, ExternalDB+csx, ctcol) {
		respondWith(w, r, nil, "Category path does not exit ...", nil, http.StatusBadRequest, false)
		return
	}

	SKUs := getSKUsInTheCategoryPath(w, r, path, ExternalDB+csx, ctcol, true)

	respondWith(w, r, nil, "Products in category path ...", bson.M{path: SKUs}, http.StatusOK, true)

}

func getProductsInCategory(w http.ResponseWriter, r *http.Request) {

	rlog.Debug("getProductsInCategory() handle function invoked ...")

	if !pre(w, r) {
		return
	}

	csx := getAccessToken(r)
	ctcol := csx + CategoryTreeExtension

	pth := strings.Split(r.URL.Path, "/")
	cid := pth[len(pth)-1]

	pathx, _ := b64.StdEncoding.DecodeString(cid)

	path := cleanCategoryPath(string(pathx))

	if !pathExists(w, r, path, ExternalDB+csx, ctcol) {
		respondWith(w, r, nil, "Category path does not exit ...", nil, http.StatusBadRequest, false)
		return
	}

	products := getProductsInTheCategoryPath(w, r, path, ExternalDB+csx, ctcol, true, csx)

	respondWith(w, r, nil, "Products in category path ...", products, http.StatusOK, true)

}

func getRootCategory(w http.ResponseWriter, r *http.Request) {

	rlog.Debug("getRootCategory() handle function invoked ...")

	if !pre(w, r) {
		return
	}

	csx := getAccessToken(r)
	ctcol := csx + CategoryTreeExtension

	cats := getRootCategories(w, r, ExternalDB+csx, ctcol)

	respondWith(w, r, nil, "Root categories ...", cats, http.StatusOK, true)

}

func getImmediateSubCategories(w http.ResponseWriter, r *http.Request) {

	rlog.Debug("getImmediateSubCategories() handle function invoked ...")

	if !pre(w, r) {
		return
	}

	csx := getAccessToken(r)
	ctcol := csx + CategoryTreeExtension

	pth := strings.Split(r.URL.Path, "/")
	cid := pth[len(pth)-1]

	catNode := getCategoryNodeByID(w, r, cid, ExternalDB+csx, ctcol)
	var childNodes []*CATEGORYTREENODE

	if catNode.Children == nil {

		respondWith(w, r, nil, "Category with ID: "+cid+" does not have a sub category ...", nil, http.StatusNotFound, false)
		return

	}

	for _, child := range catNode.Children {
		childNodes = append(childNodes, getCategoryNode(w, r, child, ExternalDB+csx, ctcol))
	}

	respondWith(w, r, nil, "Immediate Sub categories ...", childNodes, http.StatusOK, true)

}

func getParentCategory(w http.ResponseWriter, r *http.Request) {

	rlog.Debug("getParentCategory() handle function invoked ...")

	if !pre(w, r) {
		return
	}

	csx := getAccessToken(r)
	ctcol := csx + CategoryTreeExtension

	pth := strings.Split(r.URL.Path, "/")
	cid := pth[len(pth)-1]

	catNode := getCategoryNodeByID(w, r, cid, ExternalDB+csx, ctcol)
	var parentNode *CATEGORYTREENODE

	if catNode.Parent == "" {

		respondWith(w, r, nil, "Category "+cid+" does not have a parent ...", nil, http.StatusNotFound, false)
		return

	}

	parentNode = getCategoryNode(w, r, catNode.Parent, ExternalDB+csx, ctcol)

	respondWith(w, r, nil, "Category parent ...", parentNode, http.StatusOK, true)

}
